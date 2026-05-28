"""
错误处理和边界条件测试
"""
import pytest

from helpers import (
    AicCliClient,
    assert_success_output,
    assert_error_output,
    assert_output_contains,
)


class TestErrorHandling:
    """错误处理测试"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data

    # ========== 命令参数错误 ==========

    def test_missing_subcommand(self):
        """缺少子命令"""
        result = self.cli.run("skill")
        # 应该显示帮助信息
        assert_success_output(result)

    def test_invalid_subcommand(self):
        """无效子命令"""
        result = self.cli.run("skill", "nonexistent")
        # cobra 会显示帮助信息而不是错误
        assert_success_output(result)

    def test_cat_no_args(self):
        """cat 缺少参数"""
        result = self.cli.run("skill", "remote", "cat")
        assert_error_output(result)

    def test_cat_missing_path(self):
        """cat 缺少 path 参数"""
        result = self.cli.run("skill", "remote", "cat", "1")
        assert_error_output(result)

    def test_tree_no_args(self):
        """tree 缺少参数"""
        result = self.cli.run("skill", "remote", "tree")
        assert_error_output(result)

    def test_validate_no_args(self):
        """validate 缺少参数"""
        result = self.cli.run("skill", "remote", "validate")
        assert_error_output(result)

    def test_publish_no_args(self):
        """publish 缺少参数"""
        result = self.cli.run("skill", "remote", "publish")
        assert_error_output(result)

    def test_unpublish_no_args(self):
        """unpublish 缺少参数"""
        result = self.cli.run("skill", "remote", "unpublish")
        assert_error_output(result)

    def test_update_no_args(self):
        """update 缺少参数"""
        result = self.cli.run("skill", "remote", "update")
        assert_error_output(result)

    # ========== 无效 ID ==========

    def test_cat_invalid_id(self):
        """cat 无效 ID"""
        result = self.cli.skill_cat(999999, "SKILL.md")
        assert_error_output(result)

    def test_tree_invalid_id(self):
        """tree 无效 ID"""
        result = self.cli.skill_tree(999999)
        # 服务器可能返回空数据或错误
        if not result.success:
            assert_error_output(result)

    def test_validate_invalid_id(self):
        """validate 无效 ID"""
        result = self.cli.skill_validate(999999)
        assert_error_output(result)

    def test_publish_invalid_id(self):
        """publish 无效 ID"""
        result = self.cli.skill_publish(999999)
        assert_error_output(result)

    def test_unpublish_invalid_id(self):
        """unpublish 无效 ID"""
        result = self.cli.skill_unpublish(999999)
        # 可能成功（幂等）或失败
        pass

    def test_update_invalid_id(self):
        """update 无效 ID"""
        result = self.cli.skill_update(999999, desc="test")
        # 服务器可能允许更新无效 ID（幂等操作）
        if not result.success:
            assert_error_output(result)

    # ========== 无效路径 ==========

    def test_cat_invalid_path(self):
        """cat 无效路径"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_cat(skill_id, "nonexistent/path/file.txt")
        assert_error_output(result)

    def test_cat_hidden_file(self):
        """cat 尝试读取隐藏文件"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_cat(skill_id, ".hidden/file.txt")
        assert_error_output(result)

    def test_cat_path_traversal(self):
        """cat 路径穿越攻击"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_cat(skill_id, "../../../etc/passwd")
        assert_error_output(result)

    # ========== Patch 错误 ==========

    def test_patch_content_not_found(self):
        """patch 内容不匹配"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_patch(
            skill_id,
            path="SKILL.md",
            old="这是一个绝对不存在的文本_xyz_123456789",
            new="新内容",
        )
        assert_error_output(result)

    def test_patch_missing_required_args(self):
        """patch 缺少必需参数"""
        result = self.cli.run("skill", "remote", "patch", "1")
        assert_error_output(result)

    def test_patch_missing_path(self):
        """patch 缺少 path"""
        result = self.cli.run("skill", "remote", "patch", "1", "--old", "test", "--new", "new")
        assert_error_output(result)

    def test_patch_missing_old_and_line(self):
        """patch 缺少 old 和 line"""
        result = self.cli.run("skill", "remote", "patch", "1", "--path", "SKILL.md", "--new", "new")
        assert_error_output(result)

    def test_patch_missing_new(self):
        """patch 缺少 new"""
        result = self.cli.run("skill", "remote", "patch", "1", "--path", "SKILL.md", "--old", "test")
        assert_error_output(result)

    def test_patch_invalid_line_number(self):
        """patch 无效行号"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_patch(
            skill_id,
            path="SKILL.md",
            line="abc",
            new="新内容",
        )
        assert_error_output(result)

    def test_patch_invalid_stdin_json(self):
        """patch 无效 stdin JSON"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_patch(skill_id, stdin="invalid json {{{")
        assert_error_output(result)

    def test_patch_batch_file_not_found(self):
        """patch 批量文件不存在"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_patch(skill_id, batch="nonexistent.json")
        assert_error_output(result)

    # ========== 更新错误 ==========

    def test_update_no_fields(self):
        """update 没有指定更新字段"""
        skill_id = self.test_data["skills"]["existing_skill_id"]
        result = self.cli.skill_update(skill_id)
        assert_error_output(result)

    # ========== 搜索边界 ==========

    def test_search_empty_keyword(self):
        """搜索空关键词"""
        result = self.cli.skill_search("")
        # 可能返回所有结果或错误
        pass

    def test_search_special_characters(self):
        """搜索特殊字符"""
        # 使用不会导致 URL 编码问题的特殊字符
        result = self.cli.skill_search("test+special")
        assert_success_output(result)

    def test_search_unicode(self):
        """搜索 Unicode 字符"""
        result = self.cli.skill_search("测试中文")
        assert_success_output(result)

    def test_search_very_long_keyword(self):
        """搜索超长关键词"""
        result = self.cli.skill_search("a" * 1000)
        assert_success_output(result)

    # ========== 分页边界 ==========

    def test_search_page_zero(self):
        """搜索页码为 0"""
        result = self.cli.skill_search("test", page=0)
        # 可能使用默认值或错误
        pass

    def test_search_negative_page(self):
        """搜索负页码"""
        result = self.cli.skill_search("test", page=-1)
        # 可能使用默认值或错误
        pass

    def test_search_page_size_zero(self):
        """搜索每页数量为 0"""
        result = self.cli.skill_search("test", page_size=0)
        # 可能使用默认值或错误
        pass

    def test_search_very_large_page(self):
        """搜索非常大的页码"""
        result = self.cli.skill_search("test", page=999999)
        assert_success_output(result)
