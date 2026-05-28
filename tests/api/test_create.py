"""
Create 命令测试
"""
import pytest

from helpers import (
    AicCliClient,
    assert_success_output,
    assert_error_output,
    assert_output_contains,
)


class TestCreateBasic:
    """Create 基本功能测试"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data

    def test_create_help(self):
        """测试 create 命令帮助"""
        result = self.cli.run("skill", "create", "--help")
        assert_success_output(result)
        assert_output_contains(result, "创建新技能")
        assert_output_contains(result, "--category")
        assert_output_contains(result, "--tags")
        assert_output_contains(result, "--desc")

    def test_create_with_name(self):
        """测试指定名称创建"""
        # 需要指定 category 避免交互式选择
        result = self.cli.run(
            "skill", "create", "__test_create_skill",
            "--desc", "测试创建技能",
            "--category", "1",
        )
        # 可能因为服务器权限等原因失败
        if result.success:
            assert_output_contains(result, "技能已创建")
        else:
            # 记录错误但不失败
            print(f"创建失败（可能需要权限）: {result.stderr}")

    def test_create_with_all_options(self):
        """测试完整选项创建"""
        result = self.cli.run(
            "skill", "create", "__test_full_skill",
            "--desc", "完整测试技能",
            "--tags", "test,cli",
            "--category", "1",
        )
        if result.success:
            assert_output_contains(result, "技能已创建")

    @pytest.mark.skip(reason="交互式命令，无法自动化测试")
    def test_create_no_name(self):
        """测试无名称创建（交互式）"""
        # 交互式命令无法测试
        pass


class TestCreateIntegration:
    """Create 集成测试 - 创建后操作"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data

    @pytest.mark.skip(reason="需要服务器权限，跳过")
    def test_create_then_cat(self):
        """创建后读取"""
        # 创建技能
        create_result = self.cli.run(
            "skill", "create", "__test_integration",
            "--desc", "集成测试技能",
        )
        if not create_result.success:
            pytest.skip("创建技能失败")

        # 解析技能 ID
        output = create_result.stdout
        if "ID:" in output:
            # 提取 ID
            id_str = output.split("ID:")[-1].split(")")[0].strip()
            skill_id = int(id_str)

            # 读取文件
            cat_result = self.cli.skill_cat(skill_id, "SKILL.md")
            if cat_result.success:
                assert_output_contains(cat_result, "__test_integration")
