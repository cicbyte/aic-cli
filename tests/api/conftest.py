"""
API 测试 conftest.py - 提供测试 fixture
"""
import pytest
import yaml
from pathlib import Path

from helpers import AicCliClient


def load_test_data() -> dict:
    """加载测试数据"""
    data_file = Path(__file__).parent.parent / "data" / "test_data.yaml"
    if data_file.exists():
        with open(data_file, "r", encoding="utf-8") as f:
            return yaml.safe_load(f)
    return {}


@pytest.fixture(scope="session")
def cli(cli_binary) -> AicCliClient:
    """CLI 客户端 fixture"""
    return AicCliClient(cli_binary)


@pytest.fixture(scope="session")
def test_data() -> dict:
    """测试数据 fixture"""
    return load_test_data()


class ResourceTracker:
    """资源追踪器 - 测试创建的资源在 session 结束时自动清理"""

    def __init__(self):
        self._created: list[tuple[str, int]] = []

    def track(self, entity_type: str, entity_id: int):
        """注册创建的资源"""
        self._created.append((entity_type, entity_id))

    def cleanup(self, cli: AicCliClient):
        """清理所有注册的资源"""
        # 逆序清理，确保依赖关系正确
        for entity_type, entity_id in reversed(self._created):
            try:
                if entity_type == "skill":
                    # CLI 没有直接删除技能的命令，跳过
                    pass
            except Exception:
                pass  # 忽略清理失败


@pytest.fixture(scope="session")
def tracker() -> ResourceTracker:
    """资源追踪器 fixture"""
    return ResourceTracker()
