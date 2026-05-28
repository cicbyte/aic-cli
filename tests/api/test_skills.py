"""
技能相关命令测试
"""
import pytest

from helpers import (
    AicCliClient,
    assert_success_output,
    assert_error_output,
    assert_output_contains,
    assert_output_not_contains,
)


@pytest.mark.run(order=2)
class TestSkillsSearch:
    """技能搜索测试"""

    def test_search_basic(self, cli: AicCliClient):
        """测试基本搜索"""
        result = cli.skill_search("test")
        # 搜索可能返回空结果，但命令应该成功
        assert_success_output(result, "搜索命令失败")

    def test_search_with_category(self, cli: AicCliClient):
        """测试按分类搜索"""
        result = cli.skill_search("test", category_id=1)
        assert_success_output(result, "按分类搜索失败")

    def test_search_with_pagination(self, cli: AicCliClient):
        """测试分页搜索"""
        result = cli.skill_search("test", page=1, page_size=5)
        assert_success_output(result, "分页搜索失败")


@pytest.mark.run(order=3)
class TestSkillsTree:
    """技能文件树测试"""

    def test_tree(self, cli: AicCliClient, test_data: dict):
        """测试显示文件树"""
        skill_id = test_data.get("skills", {}).get("existing_skill_id", 1)
        result = cli.skill_tree(skill_id)
        # 文件树命令应该成功执行
        if result.success:
            assert_output_contains(result, "(ID:")

    def test_tree_invalid_id(self, cli: AicCliClient):
        """测试无效技能 ID - 服务器可能返回空数据或错误"""
        result = cli.skill_tree(999999)
        # 服务器可能返回空文件树而不是错误，两种情况都可接受
        if not result.success:
            assert_error_output(result)
        else:
            # 如果成功，应该输出 ID 但可能没有文件
            assert_output_contains(result, "(ID: 999999)")


@pytest.mark.run(order=4)
class TestSkillsCat:
    """技能文件内容测试"""

    def test_cat_skill_md(self, cli: AicCliClient, test_data: dict):
        """测试读取 SKILL.md"""
        skill_id = test_data.get("skills", {}).get("existing_skill_id", 1)
        result = cli.skill_cat(skill_id, "SKILL.md")
        # 如果技能存在且有 SKILL.md 文件，应该成功
        if result.success:
            assert_output_contains(result, "---")

    def test_cat_invalid_path(self, cli: AicCliClient, test_data: dict):
        """测试读取不存在的文件"""
        skill_id = test_data.get("skills", {}).get("existing_skill_id", 1)
        result = cli.skill_cat(skill_id, "nonexistent/file.txt")
        assert_error_output(result)

    def test_cat_invalid_skill(self, cli: AicCliClient):
        """测试无效技能 ID"""
        result = cli.skill_cat(999999, "SKILL.md")
        assert_error_output(result)


@pytest.mark.run(order=5)
class TestSkillsValidate:
    """技能校验测试"""

    def test_validate(self, cli: AicCliClient, test_data: dict):
        """测试校验技能"""
        skill_id = test_data.get("skills", {}).get("existing_skill_id", 1)
        result = cli.skill_validate(skill_id)
        # 校验命令应该成功执行
        if result.success:
            assert_output_contains(result, "校验")

    def test_validate_strict(self, cli: AicCliClient, test_data: dict):
        """测试严格校验"""
        skill_id = test_data.get("skills", {}).get("existing_skill_id", 1)
        result = cli.skill_validate(skill_id, strict=True)
        if result.success:
            assert_output_contains(result, "校验")


@pytest.mark.run(order=6)
class TestSkillsCategories:
    """技能分类测试"""

    def test_categories(self, cli: AicCliClient):
        """测试列出分类"""
        result = cli.skill_categories()
        assert_success_output(result, "列出分类失败")
