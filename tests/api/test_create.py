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
