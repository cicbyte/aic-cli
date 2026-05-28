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
uv run pytest -v

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
├── test_workspace/          # 测试工作目录（已忽略）
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
    ├── test_edit.py         # Edit 命令测试
    ├── test_errors.py       # 错误处理测试
    ├── test_lifecycle.py    # 生命周期测试
    ├── test_patch.py        # Patch 命令测试
    ├── test_remote.py       # Remote 全流程测试
    └── test_skills.py       # 技能相关测试
```

## 测试数据

测试数据位于 `data/test_data.yaml`，包含：
- 技能创建/更新/发布数据
- 文件操作测试数据
- 错误测试数据

## 注意事项

1. 测试需要 aic 服务器运行中
2. 测试会自动从 aic-cli 配置文件读取 API 地址和 Token
3. 测试不会修改服务器数据（只读操作）
4. 交互式测试（edit 命令）已跳过，需要模拟用户输入
