#!/usr/bin/env python3
"""
生成 Release Notes

从上次 tag 到 HEAD 的 commit 消息中提取 release notes，写入 msg.txt
"""

import subprocess
import sys
from pathlib import Path


def run_git(*args):
    """运行 git 命令并返回输出"""
    result = subprocess.run(
        ["git"] + list(args),
        capture_output=True,
        text=True,
        check=False
    )
    if result.returncode != 0:
        return None
    return result.stdout.strip()


def get_last_tag():
    """获取最近的 tag"""
    tag = run_git("describe", "--tags", "--abbrev=0")
    if not tag:
        # 如果没有 tag，返回空字符串表示使用全部提交
        return ""
    return tag


def gen_release_notes():
    """生成 release notes"""
    # 获取最近 tag
    last_tag = get_last_tag()

    # 确定提交范围
    if last_tag:
        range_ref = f"{last_tag}..HEAD"
        print(f"从 tag {last_tag} 到 HEAD 生成 release notes")
    else:
        range_ref = "HEAD"
        print("未找到 tag，使用全部提交记录")

    # 提取 commit 标题，排除 merge commit 和 docs: 开头的 commit
    result = run_git(
        "log", range_ref,
        "--pretty=format:%s",
        "--no-merges",
        "--invert-grep",
        "--grep=^docs:"
    )

    if not result:
        print("没有找到相关 commit")
        return False

    # 写入 msg.txt
    output_file = Path("msg.txt")
    with open(output_file, "w", encoding="utf-8") as f:
        f.write("# Release Notes\n\n")
        f.write(result)

    print(f"\nRelease notes 已写入 {output_file}")
    print(f"共 {len(result.splitlines())} 个 commit")

    # 显示预览
    print("\n=== 预览 ===")
    print(result[:500] + "..." if len(result) > 500 else result)

    return True


if __name__ == "__main__":
    try:
        success = gen_release_notes()
        sys.exit(0 if success else 1)
    except Exception as e:
        print(f"错误: {e}", file=sys.stderr)
        sys.exit(1)
