"""
通用断言辅助函数
"""
from .cli_client import CliResult


def assert_success_output(result: CliResult, msg: str = None):
    """断言命令执行成功"""
    assert result.success, (
        f"命令执行失败 (exit code {result.returncode})"
        f"{': ' + msg if msg else ''}"
        f"\nstdout: {result.stdout[:500]}"
        f"\nstderr: {result.stderr[:500]}"
    )


def assert_error_output(result: CliResult, expected_msg: str = None):
    """断言命令执行失败"""
    assert not result.success, (
        f"命令应该失败但成功了"
        f"\nstdout: {result.stdout[:500]}"
    )
    if expected_msg:
        assert expected_msg in result.output, (
            f"输出中未找到预期错误信息: {expected_msg}"
            f"\n实际输出: {result.output[:500]}"
        )


def assert_output_contains(result: CliResult, text: str):
    """断言输出包含指定文本"""
    assert text in result.output, (
        f"输出中未找到: {text}"
        f"\n实际输出: {result.output[:500]}"
    )


def assert_output_not_contains(result: CliResult, text: str):
    """断言输出不包含指定文本"""
    assert text not in result.output, (
        f"输出中不应包含: {text}"
        f"\n实际输出: {result.output[:500]}"
    )
