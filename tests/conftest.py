"""
根级 conftest.py - 自动读取项目配置
"""
import json
import os
import subprocess
from pathlib import Path

import pytest
import yaml


def _get_project_root() -> Path:
    """获取项目根目录"""
    return Path(__file__).parent.parent


def _load_cli_config() -> dict:
    """从 aic-cli 配置文件加载配置"""
    config_path = Path.home() / ".config" / "aic-cli" / "config.yaml"
    if not config_path.exists():
        # Windows 路径
        config_path = Path.home() / ".aic-cli" / "config.yaml"

    if config_path.exists():
        with open(config_path, "r", encoding="utf-8") as f:
            config = yaml.safe_load(f)
            if config and "aic" in config:
                return config["aic"]

    return {}


def _get_base_url() -> str:
    """获取 API 基础 URL"""
    config = _load_cli_config()
    base_url = config.get("base_url", "http://localhost:8000")
    # 清理 URL
    base_url = base_url.rstrip("/")
    if not base_url.startswith("http"):
        base_url = f"http://{base_url}"
    return base_url


def _get_token() -> str:
    """获取认证 Token"""
    config = _load_cli_config()
    return config.get("token", "")


# 模块级常量
CLI_PROJECT_ROOT = _get_project_root()
API_BASE_URL = _get_base_url()
AUTH_TOKEN = _get_token()


@pytest.fixture(scope="session")
def project_root() -> Path:
    """项目根目录"""
    return CLI_PROJECT_ROOT


@pytest.fixture(scope="session")
def api_base_url() -> str:
    """API 基础 URL"""
    return API_BASE_URL


@pytest.fixture(scope="session")
def auth_token() -> str:
    """认证 Token"""
    return AUTH_TOKEN


@pytest.fixture(scope="session")
def cli_binary(project_root) -> str:
    """CLI 二进制文件路径"""
    # 先尝试构建
    binary_name = "aic-cli"
    if os.name == "nt":
        binary_name = "aic-cli.exe"

    binary_path = project_root / "dist" / binary_name

    if not binary_path.exists():
        # 尝试构建
        print(f"构建 CLI: {project_root}")
        result = subprocess.run(
            ["go", "build", "-o", str(binary_path), "."],
            cwd=project_root,
            capture_output=True,
            text=True,
        )
        if result.returncode != 0:
            pytest.fail(f"构建 CLI 失败: {result.stderr}")

    return str(binary_path)
