"""
Edit 命令测试 - 交互式编辑功能
"""
import pytest

from helpers import (
    AicCliClient,
    assert_success_output,
    assert_error_output,
    assert_output_contains,
)


class TestEditBasic:
    """Edit 基本功能测试"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data
        self.skill_id = test_data["skills"]["existing_skill_id"]

    def test_edit_help(self):
        """测试 edit 命令帮助"""
        result = self.cli.run("skill", "edit", "--help")
        assert_success_output(result)
        assert_output_contains(result, "编辑远程技能文件")
        assert_output_contains(result, "--dry-run")

    def test_edit_no_args(self):
        """测试 edit 缺少参数"""
        result = self.cli.run("skill", "edit")
        assert_error_output(result)

    def test_edit_invalid_skill_id(self):
        """测试无效技能 ID"""
        # edit 是交互式的，这里测试参数验证
        result = self.cli.run("skill", "edit", "abc")
        assert_error_output(result)

    def test_edit_dry_run_help(self):
        """测试 dry-run 选项在帮助中"""
        result = self.cli.run("skill", "edit", "--help")
        assert_success_output(result)
        assert_output_contains(result, "--dry-run")


class TestEditInteractive:
    """Edit 交互式测试（需要模拟用户输入）"""

    @pytest.fixture(autouse=True)
    def setup(self, cli, test_data):
        self.cli = cli
        self.test_data = test_data
        self.skill_id = test_data["skills"]["existing_skill_id"]

    @pytest.mark.skip(reason="交互式命令需要模拟用户输入，暂跳过")
    def test_edit_select_file(self):
        """测试文件选择交互"""
        # 这个测试需要模拟用户输入，暂跳过
        pass

    @pytest.mark.skip(reason="交互式命令需要打开编辑器，暂跳过")
    def test_edit_open_editor(self):
        """测试打开编辑器"""
        # 这个测试需要打开编辑器，暂跳过
        pass
