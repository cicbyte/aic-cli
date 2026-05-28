# AIC CLI 测试

基于 pytest 的集成测试框架，参考 aic 项目的测试机制。

## 快速开始

### 1. 安装依赖

```bash
cd tests
uv sync
```

### 2. 构建 CLI

```bash
# 在项目根目录
go build -o dist/aic-cli.exe .
```

### 3. 运行测试

```bash
cd tests

# 运行所有测试
uv run pytest

# 运行指定测试文件
uv run pytest api/test_basic.py -v

# 运行指定测试类
uv run pytest api/test_skills.py::TestSkillsSearch -v

# 运行指定测试方法
uv run pytest api/test_skills.py::TestSkillsSearch::test_search_basic -v

# 失败即停
uv run pytest -x

# 按名称筛选
uv run pytest -k "test_search"
```

## 目录结构

```
tests/
├── conftest.py              # 根级 fixture：自动读取项目配置
├── pyproject.toml           # 依赖声明 + pytest 配置
├── README.md                # 本文件
├── data/
│   └── test_data.yaml       # 测试数据
├── helpers/
│   ├── __init__.py
│   ├── cli_client.py        # CLI 客户端封装
│   └── assertions.py        # 通用断言
└── api/
    ├── __init__.py
    ├── conftest.py          # API fixture
    ├── test_basic.py        # 基础命令测试
    └── test_skills.py       # 技能相关测试
```

## 测试组织

测试按模块组织，使用 `@pytest.mark.run(order=N)` 控制执行顺序：

| 测试类 | order | 职责 |
|:---|:---|:---|
| `TestBasicCommands` | 1 | 版本、服务器状态 |
| `TestSkillsSearch` | 2 | 技能搜索 |
| `TestSkillsTree` | 3 | 文件树显示 |
| `TestSkillsCat` | 4 | 文件内容读取 |
| `TestSkillsValidate` | 5 | 技能校验 |
| `TestSkillsCategories` | 6 | 分类列表 |

## 测试数据

测试数据位于 `data/test_data.yaml`，使用 `__test_` 前缀标识测试数据。

## 注意事项

1. 测试需要 aic 服务器运行中
2. 测试会自动从 aic-cli 配置文件读取 API 地址和 Token
3. 测试不会修改服务器数据（只读操作）
