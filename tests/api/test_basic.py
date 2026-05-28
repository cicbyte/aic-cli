"""
基础命令测试 - 版本、服务器状态等
"""
import pytest

from helpers import AicCliClient, assert_success_output, assert_output_contains


@pytest.mark.run(order=1)
class TestBasicCommands:
    """基础命令测试"""

    def test_version(self, cli: AicCliClient):
        """测试 version 命令"""
        result = cli.version()
        assert_success_output(result, "version 命令失败")
        assert_output_contains(result, "aic-cli")

    def test_server_status(self, cli: AicCliClient):
        """测试 server status 命令"""
        result = cli.server_status()
        # 服务器可能未运行，这里只测试命令能执行
        # 如果服务器运行中，应该返回成功
        if result.success:
            assert_output_contains(result, "连接")
