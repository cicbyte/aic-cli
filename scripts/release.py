#!/usr/bin/env python3
"""
版本发布脚本

自动升级版本号、提交、打 tag 并推送，触发 CI 构建
"""

import argparse
import re
import subprocess
import sys
from pathlib import Path


def run_git(*args, check=True):
    """运行 git 命令"""
    print(f"$ git {' '.join(args)}")
    result = subprocess.run(
        ["git"] + list(args),
        capture_output=True,
        text=True,
        check=check
    )
    if result.stdout:
        print(result.stdout.strip())
    return result


def read_version():
    """读取 VERSION 文件"""
    version_file = Path("VERSION")
    if not version_file.exists():
        print("错误: VERSION 文件不存在", file=sys.stderr)
        sys.exit(1)

    return version_file.read_text().strip()


def write_version(version: str):
    """写入版本号到 VERSION 文件"""
    version_file = Path("VERSION")
    version_file.write_text(version + "\n")
    print(f"更新 VERSION: {version}")


def validate_version(version: str):
    """校验版本号格式"""
    pattern = r"^\d+\.\d+\.\d+$"
    if not re.match(pattern, version):
        print(f"错误: 版本号格式无效，应为 MAJOR.MINOR.PATCH (如: 1.0.0)", file=sys.stderr)
        sys.exit(1)
    return version


def bump_version(current: str):
    """自动升级 patch 版本号"""
    parts = current.split(".")
    major, minor, patch = int(parts[0]), int(parts[1]), int(parts[2])
    patch += 1
    return f"{major}.{minor}.{patch}"


def read_release_notes():
    """读取并删除 msg.txt"""
    msg_file = Path("msg.txt")
    if not msg_file.exists():
        print("错误: msg.txt 不存在，请先运行: python scripts/gen_release_notes.py", file=sys.stderr)
        sys.exit(1)

    release_notes = msg_file.read_text().strip()
    msg_file.unlink()  # 读取后删除
    print(f"读取 release notes，已删除 {msg_file}")
    return release_notes


def release(version: str = None):
    """执行发布流程"""
    # 读取当前版本号
    current_version = read_version()
    print(f"当前版本: {current_version}")

    # 确定新版本号
    if version:
        # 使用指定的版本号
        new_version = validate_version(version)
    else:
        # 自动 patch 升级
        new_version = bump_version(current_version)
        print(f"自动 patch 升级: {current_version} -> {new_version}")

    # 校验版本号不能与当前相同
    if new_version == current_version:
        print(f"错误: 新版本号与当前版本号相同 ({current_version})", file=sys.stderr)
        sys.exit(1)

    # 读取 release notes
    release_notes = read_release_notes()
    print(f"\nRelease Notes:\n{release_notes}\n")

    # 写入新版本号
    write_version(new_version)

    # Git 操作
    print("\n开始 Git 操作...")

    # 提交 VERSION 文件
    run_git("add", "VERSION")
    run_git("commit", "-m", release_notes)

    # 推送
    run_git("push")

    # 创建 tag
    tag = f"v{new_version}"
    run_git("tag", tag)
    run_git("push", "origin", tag)

    print(f"\n✅ 发布完成！")
    print(f"版本: {new_version}")
    print(f"Tag: {tag}")
    print(f"\nGitHub Actions 将自动构建并创建 Release")


def main():
    parser = argparse.ArgumentParser(
        description="版本发布脚本",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  python scripts/release.py              # 自动 patch 升级 (0.1.0 -> 0.1.1)
  python scripts/release.py 1.0.0        # 指定版本号
  python scripts/release.py 2.1.3        # 指定版本号
        """
    )

    parser.add_argument(
        "version",
        nargs="?",
        help="新版本号 (MAJOR.MINOR.PATCH)，不指定则自动 patch 升级"
    )

    args = parser.parse_args()
    release(args.version)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n中断发布", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"错误: {e}", file=sys.stderr)
        sys.exit(1)
