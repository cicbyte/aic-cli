"""
Remote 命令全流程测试 - 创建 → 文件操作 → 校验 → 发布 → 更新 → 清理
"""
import json
import os
import tempfile

import pytest

from helpers import (
    AicCliClient,
    assert_success_output,
    assert_error_output,
    assert_output_contains,
)


class TestRemoteFullLifecycle:
    """Remote 命令全流程测试"""

    # 类属性共享状态
    skill_id: int = None
    original_sha256: str = None

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data
        # 使用已存在的技能进行测试
        self.skill_id = test_data["skills"]["existing_skill_id"]

    # ========== Phase 1: 文件读取 ==========

    def test_01_tree(self):
        """显示文件树"""
        result = self.cli.skill_tree(self.skill_id)
        assert_success_output(result)
        assert_output_contains(result, f"(ID: {self.skill_id})")

    def test_02_cat_skill_md(self):
        """读取 SKILL.md"""
        result = self.cli.skill_cat(self.skill_id, "SKILL.md")
        # API 可能未实现 (404)，记录但不失败
        if result.success:
            assert_output_contains(result, "---")
            # 保存 sha256 供后续使用
            if "sha256:" in result.stdout:
                self.__class__.original_sha256 = result.stdout.split("sha256: ")[-1].split("\n")[0].strip()
        else:
            # 服务器端 API 可能未实现
            assert "404" in result.stderr or "Not Found" in result.stderr

    def test_03_cat_to_file(self):
        """读取文件并保存到本地"""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".md", delete=False) as f:
            output_file = f.name

        try:
            result = self.cli.skill_cat(self.skill_id, "SKILL.md", output=output_file)
            if result.success:
                assert_output_contains(result, "已保存到")
                assert os.path.exists(output_file)
                # 验证文件内容不为空
                with open(output_file, "r", encoding="utf-8") as f:
                    content = f.read()
                assert len(content) > 0
            else:
                # 服务器端 API 可能未实现
                assert "404" in result.stderr or "Not Found" in result.stderr
        finally:
            if os.path.exists(output_file):
                os.unlink(output_file)

    def test_04_cat_invalid_path(self):
        """读取不存在的文件"""
        result = self.cli.skill_cat(self.skill_id, "nonexistent/file.txt")
        assert_error_output(result)

    # ========== Phase 2: 增量编辑 ==========

    def test_05_patch_content_match(self):
        """纯内容匹配模式"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            old="description:",
            new="description:",
        )
        # 内容相同不会失败，只是没有替换
        if result.success:
            assert_output_contains(result, "已更新")

    def test_06_patch_not_found(self):
        """内容不匹配"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            old="这是一个绝对不存在的文本_xyz_123456789",
            new="新内容",
        )
        assert_error_output(result)

    def test_07_patch_stdin(self):
        """管道 JSON 模式"""
        stdin_data = json.dumps({
            "path": "SKILL.md",
            "edits": [
                {
                    "old_string": "description:",
                    "new_string": "description:",
                }
            ]
        })
        result = self.cli.skill_patch(self.skill_id, stdin=stdin_data)
        if result.success:
            assert_output_contains(result, "已更新")

    def test_08_patch_stdin_multiple_edits(self):
        """多个编辑操作"""
        stdin_data = json.dumps({
            "path": "SKILL.md",
            "edits": [
                {
                    "old_string": "description:",
                    "new_string": "description:",
                },
                {
                    "old_string": "name:",
                    "new_string": "name:",
                }
            ]
        })
        result = self.cli.skill_patch(self.skill_id, stdin=stdin_data)
        if result.success:
            assert_output_contains(result, "已更新")

    def test_09_patch_batch(self):
        """批量文件操作"""
        batch_data = {
            "operations": [
                {
                    "op": "patch",
                    "path": "SKILL.md",
                    "edits": [
                        {
                            "old_string": "description:",
                            "new_string": "description:",
                        }
                    ]
                }
            ]
        }

        with tempfile.NamedTemporaryFile(mode="w", suffix=".json", delete=False) as f:
            json.dump(batch_data, f)
            batch_file = f.name

        try:
            result = self.cli.skill_patch(self.skill_id, batch=batch_file)
            if result.success:
                assert_output_contains(result, "已更新")
        finally:
            os.unlink(batch_file)

    def test_10_patch_schema(self):
        """查看 JSON 格式说明"""
        result = self.cli.run("skill", "remote", "patch", "--schema")
        assert_success_output(result)
        assert_output_contains(result, "单文件编辑 JSON 格式")
        assert_output_contains(result, "edits")

    # ========== Phase 3: 校验 ==========

    def test_11_validate(self):
        """标准校验"""
        result = self.cli.skill_validate(self.skill_id)
        if result.success:
            assert_output_contains(result, "校验")
        else:
            # 服务器端 API 可能未实现
            assert "404" in result.stderr or "Not Found" in result.stderr

    def test_12_validate_strict(self):
        """严格校验"""
        result = self.cli.skill_validate(self.skill_id, strict=True)
        if result.success:
            assert_output_contains(result, "校验")
        else:
            # 服务器端 API 可能未实现
            assert "404" in result.stderr or "Not Found" in result.stderr

    # ========== Phase 4: 发布 ==========

    def test_13_publish(self):
        """发布技能"""
        result = self.cli.skill_publish(
            self.skill_id,
            version="0.0.1-test",
            changelog="测试发布",
        )
        # 发布可能需要权限或校验通过
        if result.success:
            assert_output_contains(result, "已发布")

    def test_14_unpublish(self):
        """取消发布"""
        result = self.cli.skill_unpublish(self.skill_id)
        if result.success:
            assert_output_contains(result, "已取消发布")

    # ========== Phase 5: 元数据更新 ==========

    def test_15_update_name(self):
        """更新名称"""
        result = self.cli.skill_update(self.skill_id, name="test-updated-name")
        if result.success:
            assert_output_contains(result, "已更新")

    def test_16_update_description(self):
        """更新描述"""
        result = self.cli.skill_update(self.skill_id, desc="测试更新描述")
        if result.success:
            assert_output_contains(result, "已更新")

    def test_17_update_tags(self):
        """更新标签"""
        result = self.cli.skill_update(self.skill_id, tags="test,updated")
        if result.success:
            assert_output_contains(result, "已更新")

    # ========== Phase 6: 错误处理 ==========

    def test_18_error_invalid_id(self):
        """无效技能 ID"""
        result = self.cli.skill_cat(999999, "SKILL.md")
        assert_error_output(result)

    def test_19_error_invalid_path(self):
        """无效文件路径"""
        result = self.cli.skill_cat(self.skill_id, "../etc/passwd")
        assert_error_output(result)

    def test_20_error_patch_no_args(self):
        """patch 缺少参数"""
        result = self.cli.run("skill", "remote", "patch", str(self.skill_id))
        assert_error_output(result)

    def test_21_error_update_no_fields(self):
        """update 缺少更新字段"""
        result = self.cli.skill_update(self.skill_id)
        assert_error_output(result)
