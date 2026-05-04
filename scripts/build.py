#!/usr/bin/env python3
"""
AIC CLI 构建脚本

功能：
1. 构建 web 前端（npm run build）
2. 复制 web 构建产物到 resources/static
3. 编译 Go 应用（支持交叉编译和版本注入）
"""

import os
import shutil
import subprocess
import sys
import time
import logging
from pathlib import Path

# 配置日志
log = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")

# 全局配置
USE_GOX = False  # 是否使用 gox 进行交叉编译

# 版本信息包路径（用于 ldflags）
VERSION_IMPORT_PATH = "github.com/cicbyte/aic-cli/cmd/version"

# 项目根目录
ROOT = Path(__file__).parent.parent


def run_command(cmd, cwd=None, shell=False, direct_output=False):
    """运行命令并实时显示输出

    Args:
        cmd: 要运行的命令
        cwd: 工作目录
        shell: 是否使用shell执行
        direct_output: 是否直接输出到终端（用于显示进度条等）
    """
    if direct_output:
        # 直接输出到终端，用于显示进度条等
        process = subprocess.Popen(
            cmd,
            stdout=sys.stdout,
            stderr=sys.stderr,
            cwd=cwd,
            shell=shell,
            universal_newlines=True
        )
        return_code = process.wait()
        return return_code
    else:
        # 通过logging输出
        process = subprocess.Popen(
            cmd,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            encoding='utf-8',
            errors='replace',
            cwd=cwd,
            shell=shell,
            bufsize=1,
            universal_newlines=True
        )

        # 实时读取并显示输出
        while True:
            output = process.stdout.readline()
            if output == '' and process.poll() is not None:
                break
            if output:
                log.info(output.strip())

        # 获取返回码
        return_code = process.poll()

        # 检查是否有错误输出
        stderr = process.stderr.read()
        if stderr:
            log.error(stderr.strip())

        return return_code


def run_npm_build():
    """构建 web 前端"""
    web_dir = ROOT / "web"
    if not web_dir.exists():
        log.info("web 目录不存在，跳过 npm 构建")
        return

    log.info(f"web 目录: {web_dir}")

    # 在Windows系统上使用npm.cmd
    npm_cmd = "npm.cmd" if os.name == 'nt' else "npm"

    try:
        return_code = run_command(
            [npm_cmd, "run", "build"],
            cwd=web_dir,
            shell=True
        )

        if return_code == 0:
            log.info("npm run build 执行成功")
        else:
            log.error("npm run build 执行失败")
            raise Exception("npm run build 执行失败")
    except Exception as e:
        log.error(f"执行npm命令时出错: {str(e)}")
        raise


def remove_old_build():
    """清理旧的构建产物"""
    target_dir = ROOT / "resources" / "static"

    # 检查并创建目标目录
    try:
        target_dir.mkdir(parents=True, exist_ok=True)
    except Exception as e:
        log.error(f"创建目录失败: {str(e)}")
        raise

    # 删除目录内容（包括子目录）
    for item in target_dir.glob('*'):
        try:
            if item.is_file():
                item.unlink()
            elif item.is_dir():
                shutil.rmtree(item)
        except Exception as e:
            log.error(f"删除 {item} 失败: {str(e)}")
            raise


def copy_new_build():
    """复制新的构建产物"""
    source_dir = ROOT / "web" / "dist"
    target_dir = ROOT / "resources" / "static"

    # 确保源目录存在
    if not source_dir.exists():
        log.warning(f"源目录不存在: {source_dir}，跳过复制")
        return

    # 复制整个目录结构
    try:
        shutil.copytree(
            source_dir,
            target_dir,
            dirs_exist_ok=True  # 允许目标目录已存在
        )
        log.info(f"成功复制文件从 {source_dir} 到 {target_dir}")
    except Exception as e:
        log.error(f"复制文件失败: {str(e)}")
        raise


def check_upx():
    """检查UPX是否可用"""
    try:
        return_code = run_command(["upx", "--version"])
        return return_code == 0
    except FileNotFoundError:
        return False


def compress_with_upx(output_name):
    """使用UPX压缩可执行文件"""
    if not check_upx():
        log.info("UPX未安装，跳过压缩步骤")
        return

    try:
        # 获取压缩前文件大小
        original_size = os.path.getsize(output_name)
        log.info(f"开始UPX压缩，原始文件大小: {original_size/1024/1024:.2f}MB")

        # 执行UPX压缩，使用--verbose选项显示详细信息，并直接输出到终端以显示进度条
        return_code = run_command(
            ["upx", "--best", "--verbose", output_name],
            direct_output=True
        )

        if return_code == 0:
            # 获取压缩后文件大小
            compressed_size = os.path.getsize(output_name)
            # 计算压缩率
            compression_ratio = (1 - compressed_size / original_size) * 100
            log.info(f"压缩完成: {original_size/1024/1024:.2f}MB -> {compressed_size/1024/1024:.2f}MB (压缩率: {compression_ratio:.1f}%)")

            # 清理UPX临时文件
            base_name = os.path.splitext(output_name)[0]  # 移除.exe扩展名
            temp_files = [
                f"{base_name}.000",
                f"{base_name}.upx"
            ]
            for temp_file in temp_files:
                if os.path.exists(temp_file):
                    try:
                        os.remove(temp_file)
                        log.info(f"已清理临时文件: {temp_file}")
                    except Exception as e:
                        log.warning(f"清理临时文件 {temp_file} 失败: {str(e)}")
        else:
            log.error("UPX压缩失败")
    except Exception as e:
        log.error(f"UPX压缩时出错: {str(e)}")


def check_gox():
    """检查gox是否已安装"""
    try:
        # gox 没有 -version 参数，直接运行 gox 命令检查是否可用
        return_code = run_command(["gox", "-h"])
        return return_code == 0
    except FileNotFoundError:
        return False


def get_version():
    """获取版本号

    优先从 VERSION 文件读取，如果没有则使用 Git 标签
    """
    # 优先从 VERSION 文件读取
    version_file = ROOT / "VERSION"
    if version_file.exists():
        return version_file.read_text().strip()

    # 回退到 Git 标签
    try:
        result = subprocess.run(
            ["git", "describe", "--tags", "--abbrev=0"],
            capture_output=True,
            text=True,
            check=False
        )
        if result.returncode == 0:
            return result.stdout.strip()
    except Exception:
        pass

    return "0.0.0-dev"


def get_build_info():
    """获取构建信息"""
    try:
        # 获取 Git 提交哈希
        commit = subprocess.run(
            ["git", "rev-parse", "HEAD"],
            capture_output=True,
            text=True,
            check=False
        ).stdout.strip()

        # 获取 Git 分支名
        branch = subprocess.run(
            ["git", "rev-parse", "--abbrev-ref", "HEAD"],
            capture_output=True,
            text=True,
            check=False
        ).stdout.strip()

        return {
            "commit": commit,
            "branch": branch,
            "build_time": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
        }
    except Exception:
        return {
            "commit": "unknown",
            "branch": "unknown",
            "build_time": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
        }


def build_with_gox():
    """使用gox进行交叉编译"""
    if not check_gox():
        log.error("gox未安装，请先安装gox: go install github.com/mitchellh/gox@v1.0.1")
        raise Exception("gox未安装")

    try:
        # 设置输出目录
        output_dir = ROOT / "build"

        # 清理输出目录
        if output_dir.exists():
            shutil.rmtree(output_dir)
        output_dir.mkdir()

        # 获取版本信息和构建信息
        version = get_version()
        build_info = get_build_info()

        log.info(f"构建版本: {version}")
        log.info(f"Git commit: {build_info['commit']}")
        log.info(f"构建时间: {build_info['build_time']}")

        # 设置编译参数
        ldflags = (
            f"-s -w "
            f"-X {VERSION_IMPORT_PATH}.Version={version} "
            f"-X {VERSION_IMPORT_PATH}.GitCommit={build_info['commit']} "
            f"-X {VERSION_IMPORT_PATH}.BuildTime={build_info['build_time']}"
        )

        # 执行gox命令，在项目根目录下执行
        return_code = run_command([
            "gox",
            "-os", "windows linux darwin",
            "-arch", "amd64",
            "-ldflags", ldflags,
            "-output", str(output_dir / "aic-cli_{{.OS}}_{{.Arch}}")
        ])

        if return_code == 0:
            log.info(f"gox交叉编译成功，版本: {version}")

            # 为Windows版本添加.exe扩展名
            for file in output_dir.glob("aic-cli_windows*"):
                if not file.name.endswith(".exe"):
                    new_path = file.with_suffix(".exe")
                    file.rename(new_path)
                    log.info(f"重命名Windows可执行文件: {file.name} -> {new_path.name}")

            # 只对当前平台的可执行文件进行UPX压缩
            current_os = "windows" if os.name == 'nt' else "linux" if sys.platform == "linux" else "darwin"
            current_pattern = f"aic-cli_{current_os}_amd64*"
            for file in output_dir.glob(current_pattern):
                if file.is_file():
                    log.info(f"开始压缩当前平台文件: {file.name}")
                    compress_with_upx(str(file))
        else:
            log.error("gox交叉编译失败")
            raise Exception("gox交叉编译失败")
    except Exception as e:
        log.error(f"使用gox编译时出错: {str(e)}")
        raise


def build_go_app():
    """编译Go应用"""
    try:
        if USE_GOX:
            build_with_gox()
            return

        # 设置输出文件名
        output_name = "aic-cli.exe" if os.name == 'nt' else "aic-cli"

        # 获取版本信息
        version = get_version()
        build_info = get_build_info()

        log.info(f"构建版本: {version}")
        log.info(f"Git commit: {build_info['commit']}")
        log.info(f"构建时间: {build_info['build_time']}")

        # 设置编译参数
        ldflags = (
            f"-s -w "
            f"-X {VERSION_IMPORT_PATH}.Version={version} "
            f"-X {VERSION_IMPORT_PATH}.GitCommit={build_info['commit']} "
            f"-X {VERSION_IMPORT_PATH}.BuildTime={build_info['build_time']}"
        )

        # 执行go build命令
        return_code = run_command(
            ["go", "build", "-ldflags", ldflags, "-o", output_name, "."],
            cwd=ROOT
        )

        if return_code == 0:
            log.info("Go应用编译成功")
            # 尝试使用UPX压缩
            compress_with_upx(output_name)
        else:
            log.error("Go应用编译失败")
            raise Exception("Go应用编译失败")
    except Exception as e:
        log.error(f"编译Go应用时出错: {str(e)}")
        raise


def main():
    """主函数"""
    log.info("=" * 60)
    log.info("AIC CLI 构建流程")
    log.info("=" * 60)

    # 1. 构建 web 前端
    if (ROOT / "web").exists():
        log.info("\n[1/4] 构建 web 前端...")
        run_npm_build()
    else:
        log.info("\n[1/4] 跳过 web 构建（web 目录不存在）")

    # 2. 清理旧构建
    log.info("\n[2/4] 清理旧构建产物...")
    if (ROOT / "resources").exists():
        remove_old_build()
        copy_new_build()
    else:
        log.info("跳过资源复制（resources 目录不存在）")

    # 3. 编译 Go 应用
    log.info("\n[3/4] 编译 Go 应用...")
    build_go_app()

    # 4. 完成
    log.info("\n[4/4] 构建完成!")
    log.info("=" * 60)

    # 显示输出文件
    output_files = []
    if USE_GOX:
        # gox 模式
        build_dir = ROOT / "build"
        if build_dir.exists():
            output_files = list(build_dir.glob("aic-cli_*"))
    else:
        # 单平台模式
        output_name = "aic-cli.exe" if os.name == 'nt' else "aic-cli"
        output_file = ROOT / output_name
        if output_file.exists():
            output_files = [output_file]

    if output_files:
        log.info("\n输出文件:")
        for file in output_files:
            size = file.stat().st_size / 1024 / 1024
            log.info(f"  {file.name:40} {size:6.2f}MB")


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        log.info("\n中断构建")
        sys.exit(1)
    except Exception as e:
        log.error(f"\n构建失败: {e}")
        sys.exit(1)
