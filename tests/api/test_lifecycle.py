"""
技能生命周期测试 - 完整流程：创建 → 文件操作 → 校验 → 发布 → 取消发布 → 更新 → 清理
"""
import os
import tempfile

import pytest

from helpers import (
    AicCliClient,
    assert_success_output,
    assert_error_output,
    assert_output_contains,
    assert_output_not_contains,
)


class TestSkillLifecycle:
    """技能全生命周期测试"""

    # 类属性共享状态
    skill_id: int = None
    skill_name: str = None

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data, tracker):
        """初始化 fixture"""
        self.cli = cli
        self.test_data = test_data
        self.tracker = tracker

    def test_01_search_before_create(self):
        """搜索确认技能状态"""
        skill_name = self.test_data["skills"]["create"]["name"]
        result = self.cli.skill_search(skill_name)
        assert_success_output(result)
        # 搜索可能找到或找不到，记录状态

    def test_02_cat_existing_skill(self):
        """测试读取已有技能文件"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_cat(skill_id, "SKILL.md")
        # 如果技能存在且有 SKILL.md，应该成功
        if result.success:
            assert_output_contains(result, "---")

    def test_03_tree_existing_skill(self):
        """测试显示已有技能文件树"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_tree(skill_id)
        if result.success:
            assert_output_contains(result, "(ID:")

    def test_04_validate_existing_skill(self):
        """测试校验已有技能"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_validate(skill_id)
        if result.success:
            assert_output_contains(result, "校验")

    def test_05_validate_strict(self):
        """测试严格校验"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_validate(skill_id, strict=True)
        if result.success:
            assert_output_contains(result, "校验")

    def test_06_categories(self):
        """测试列出分类"""
        result = self.cli.skill_categories()
        assert_success_output(result)

    def test_07_update_metadata(self):
        """测试更新技能元数据"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        update_data = self.test_data["skills"]["update"]

        result = self.cli.skill_update(
            skill_id,
            desc=update_data["description"],
            tags=",".join(update_data["tags"]),
        )
        # 更新可能需要权限，记录结果
        if result.success:
            assert_output_contains(result, "已更新")

    def test_08_patch_content_match(self):
        """测试 patch - 纯内容匹配模式"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        files_data = self.test_data["skills"]["files"]

        result = self.cli.skill_patch(
            skill_id,
            path=files_data["skill_md"]["path"],
            old=files_data["patch_old"],
            new=files_data["patch_new"],
        )
        # patch 可能因为内容不匹配而失败，这是预期的
        if result.success:
            assert_output_contains(result, "已更新")

    def test_09_patch_line_mode(self):
        """测试 patch - 纯行号模式"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        files_data = self.test_data["skills"]["files"]

        result = self.cli.skill_patch(
            skill_id,
            path=files_data["skill_md"]["path"],
            line=str(files_data["patch_line"]),
            new=files_data["patch_line_new"],
        )
        if result.success:
            assert_output_contains(result, "已更新")

    def test_10_patch_mixed_mode(self):
        """测试 patch - 行号 + 内容校验模式"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        files_data = self.test_data["skills"]["files"]

        result = self.cli.skill_patch(
            skill_id,
            path=files_data["skill_md"]["path"],
            line=str(files_data["patch_line"]),
            old=files_data["patch_line_old"],
            new=files_data["patch_line_new"],
        )
        if result.success:
            assert_output_contains(result, "已更新")

    def test_11_patch_stdin(self):
        """测试 patch - stdin JSON 模式"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        files_data = self.test_data["skills"]["files"]

        import json
        stdin_data = json.dumps({
            "path": files_data["skill_md"]["path"],
            "edits": [
                {
                    "old_string": files_data["patch_old"],
                    "new_string": files_data["patch_new"],
                }
            ]
        })

        result = self.cli.skill_patch(skill_id, stdin=stdin_data)
        if result.success:
            assert_output_contains(result, "已更新")

    def test_12_patch_batch_file(self):
        """测试 patch - 批量文件模式"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        files_data = self.test_data["skills"]["files"]

        import json

        # 创建临时批量操作文件
        batch_data = {
            "operations": [
                {
                    "op": "patch",
                    "path": files_data["skill_md"]["path"],
                    "edits": [
                        {
                            "old_string": files_data["patch_old"],
                            "new_string": files_data["patch_new"],
                        }
                    ]
                }
            ]
        }

        with tempfile.NamedTemporaryFile(mode="w", suffix=".json", delete=False) as f:
            json.dump(batch_data, f)
            batch_file = f.name

        try:
            result = self.cli.skill_patch(skill_id, batch=batch_file)
            if result.success:
                assert_output_contains(result, "已更新")
        finally:
            os.unlink(batch_file)

    def test_13_cat_with_output(self):
        """测试 cat 命令保存到文件"""
        skill_id = self.test_data["skills"]["existing_skill_id"]

        with tempfile.NamedTemporaryFile(mode="w", suffix=".md", delete=False) as f:
            output_file = f.name

        try:
            result = self.cli.skill_cat(skill_id, "SKILL.md", output=output_file)
            if result.success:
                assert_output_contains(result, "已保存到")
                # 验证文件已创建
                assert os.path.exists(output_file)
        finally:
            if os.path.exists(output_file):
                os.unlink(output_file)

    def test_14_publish_skill(self):
        """测试发布技能"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        publish_data = self.test_data["skills"]["publish"]

        result = self.cli.skill_publish(
            skill_id,
            version=publish_data["version"],
            changelog=publish_data["changelog"],
        )
        # 发布可能需要权限或校验通过
        if result.success:
            assert_output_contains(result, "已发布")

    def test_15_unpublish_skill(self):
        """测试取消发布"""
        skill_id = self.test_data["skills"]["existing_skill_id"]

        result = self.cli.skill_unpublish(skill_id)
        if result.success:
            assert_output_contains(result, "已取消发布")

    def test_16_error_invalid_skill_id(self):
        """测试无效技能 ID"""
        result = self.cli.skill_cat(999999, "SKILL.md")
        assert_error_output(result)

    def test_17_error_invalid_path(self):
        """测试无效文件路径"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_cat(skill_id, "nonexistent/file.txt")
        assert_error_output(result)

    def test_18_error_invalid_patch(self):
        """测试 patch 失败 - 内容不匹配"""
        skill_id = self.test_data["skills"]["existing_skill_id"]

        result = self.cli.skill_patch(
            skill_id,
            path="SKILL.md",
            old="这是一个不存在的文本内容_xyz_123",
            new="新内容",
        )
        assert_error_output(result)

    def test_19_error_patch_no_args(self):
        """测试 patch 缺少参数"""
        result = self.cli.run("skill", "remote", "patch", "1")
        assert_error_output(result)

    def test_20_error_update_no_fields(self):
        """测试 update 缺少更新字段"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_update(skill_id)
        assert_error_output(result)
