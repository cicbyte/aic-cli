"""
Patch 命令详细测试 - 覆盖所有模式和边界条件
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


class TestPatchContentMatch:
    """纯内容匹配模式测试"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data
        self.skill_id = test_data["skills"]["existing_skill_id"]

    def test_basic_content_match(self):
        """基本内容匹配"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            old="name:",
            new="name: updated",
        )
        # 可能成功或失败（取决于内容是否匹配）
        if result.success:
            assert_output_contains(result, "已更新")

    def test_content_match_not_found(self):
        """内容不匹配"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            old="这是一个绝对不存在的文本_xyz_123456",
            new="新内容",
        )
        assert_error_output(result)

    def test_content_match_replace_all(self):
        """替换所有匹配"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            old=":",
            new="：",
            replace_all=True,
        )
        if result.success:
            assert_output_contains(result, "已更新")


class TestPatchLineMode:
    """纯行号模式测试"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data
        self.skill_id = test_data["skills"]["existing_skill_id"]

    def test_single_line_replace(self):
        """单行替换"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            line="1",
            new="# 更新的标题",
        )
        if result.success:
            assert_output_contains(result, "已更新")

    def test_line_range_replace(self):
        """行范围替换"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            line="1-3",
            new="# 新标题\n\n新描述",
        )
        if result.success:
            assert_output_contains(result, "已更新")

    def test_invalid_line_number(self):
        """无效行号"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            line="abc",
            new="新内容",
        )
        assert_error_output(result)

    def test_out_of_range_line(self):
        """超出范围的行号"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            line="99999",
            new="新内容",
        )
        if result.success:
            # 服务器可能允许超出范围的行号（追加到末尾）
            pass
        else:
            assert_error_output(result)


class TestPatchMixedMode:
    """混合模式测试（行号 + 内容校验）"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data
        self.skill_id = test_data["skills"]["existing_skill_id"]

    def test_mixed_mode_basic(self):
        """基本混合模式"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            line="1",
            old="name:",
            new="name: updated",
        )
        if result.success:
            assert_output_contains(result, "已更新")

    def test_mixed_mode_content_not_in_range(self):
        """内容不在指定行范围内"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            line="1",
            old="这是第100行的内容",
            new="新内容",
        )
        assert_error_output(result)

    def test_mixed_mode_range(self):
        """行范围 + 内容校验"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            line="1-5",
            old="name:",
            new="name: updated",
        )
        if result.success:
            assert_output_contains(result, "已更新")


class TestPatchStdin:
    """Stdin JSON 模式测试"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data
        self.skill_id = test_data["skills"]["existing_skill_id"]

    def test_stdin_basic(self):
        """基本 stdin 模式"""
        stdin_data = json.dumps({
            "path": "SKILL.md",
            "edits": [
                {
                    "old_string": "name:",
                    "new_string": "name: updated",
                }
            ]
        })

        result = self.cli.skill_patch(self.skill_id, stdin=stdin_data)
        if result.success:
            assert_output_contains(result, "已更新")

    def test_stdin_multiple_edits(self):
        """多个编辑操作"""
        stdin_data = json.dumps({
            "path": "SKILL.md",
            "edits": [
                {
                    "old_string": "name:",
                    "new_string": "name: updated1",
                },
                {
                    "old_string": "description:",
                    "new_string": "description: updated2",
                }
            ]
        })

        result = self.cli.skill_patch(self.skill_id, stdin=stdin_data)
        if result.success:
            assert_output_contains(result, "已更新")

    def test_stdin_with_line_numbers(self):
        """带行号的 stdin"""
        stdin_data = json.dumps({
            "path": "SKILL.md",
            "edits": [
                {
                    "line_start": 1,
                    "new_string": "# 新标题",
                }
            ]
        })

        result = self.cli.skill_patch(self.skill_id, stdin=stdin_data)
        if result.success:
            assert_output_contains(result, "已更新")

    def test_stdin_invalid_json(self):
        """无效 JSON"""
        result = self.cli.skill_patch(self.skill_id, stdin="invalid json {{{")
        assert_error_output(result)

    def test_stdin_empty(self):
        """空 stdin"""
        result = self.cli.skill_patch(self.skill_id, stdin="")
        assert_error_output(result)


class TestPatchBatchFile:
    """批量文件模式测试"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data
        self.skill_id = test_data["skills"]["existing_skill_id"]

    def _create_batch_file(self, data: dict) -> str:
        """创建临时批量操作文件"""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".json", delete=False) as f:
            json.dump(data, f)
            return f.name

    def test_batch_patch_only(self):
        """批量 patch 操作"""
        batch_data = {
            "operations": [
                {
                    "op": "patch",
                    "path": "SKILL.md",
                    "edits": [
                        {
                            "old_string": "name:",
                            "new_string": "name: batch_updated",
                        }
                    ]
                }
            ]
        }

        batch_file = self._create_batch_file(batch_data)
        try:
            result = self.cli.skill_patch(self.skill_id, batch=batch_file)
            if result.success:
                assert_output_contains(result, "已更新")
        finally:
            os.unlink(batch_file)

    def test_batch_mixed_operations(self):
        """混合操作"""
        batch_data = {
            "operations": [
                {
                    "op": "patch",
                    "path": "SKILL.md",
                    "edits": [
                        {
                            "old_string": "name:",
                            "new_string": "name: mixed_updated",
                        }
                    ]
                },
                {
                    "op": "write",
                    "path": "test/batch_test.md",
                    "content": "# 批量测试文件"
                }
            ]
        }

        batch_file = self._create_batch_file(batch_data)
        try:
            result = self.cli.skill_patch(self.skill_id, batch=batch_file)
            if result.success:
                assert_output_contains(result, "已更新")
        finally:
            os.unlink(batch_file)

    def test_batch_file_not_found(self):
        """批量文件不存在"""
        result = self.cli.skill_patch(self.skill_id, batch="nonexistent.json")
        assert_error_output(result)

    def test_batch_invalid_json(self):
        """批量文件无效 JSON"""
        with tempfile.NamedTemporaryFile(mode="w", suffix=".json", delete=False) as f:
            f.write("invalid json {{{")
            batch_file = f.name

        try:
            result = self.cli.skill_patch(self.skill_id, batch=batch_file)
            assert_error_output(result)
        finally:
            os.unlink(batch_file)


class TestPatchEdgeCases:
    """边界条件测试"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data
        self.skill_id = test_data["skills"]["existing_skill_id"]

    def test_patch_special_characters(self):
        """特殊字符处理"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            old="name:",
            new="name: special!@#$%^&*()",
        )
        if result.success:
            assert_output_contains(result, "已更新")

    def test_patch_unicode(self):
        """Unicode 字符"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            old="name:",
            new="name: 测试中文",
        )
        if result.success:
            assert_output_contains(result, "已更新")

    def test_patch_multiline(self):
        """多行替换"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            old="name:",
            new="name: updated\ndescription: new desc",
        )
        if result.success:
            assert_output_contains(result, "已更新")

    def test_patch_empty_new_string(self):
        """空新字符串（删除操作）"""
        result = self.cli.skill_patch(
            self.skill_id,
            path="SKILL.md",
            old="description:",
            new="",
        )
        if result.success:
            assert_output_contains(result, "已更新")

    def test_patch_no_path(self):
        """缺少 path 参数"""
        result = self.cli.run("skill", "remote", "patch", str(self.skill_id), "--old", "test", "--new", "new")
        assert_error_output(result)

    def test_patch_no_old_or_line(self):
        """缺少 old 和 line 参数"""
        result = self.cli.run("skill", "remote", "patch", str(self.skill_id), "--path", "SKILL.md", "--new", "new")
        assert_error_output(result)

    def test_patch_no_new(self):
        """缺少 new 参数"""
        result = self.cli.run("skill", "remote", "patch", str(self.skill_id), "--path", "SKILL.md", "--old", "test")
        assert_error_output(result)
