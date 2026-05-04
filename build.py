#!/usr/bin/env python3
"""
AIC CLI 构建入口脚本

这是项目的快速构建入口，实际构建逻辑在 scripts/build.py 中
"""

import sys
from pathlib import Path

# 添加 scripts 目录到 Python 路径
scripts_dir = Path(__file__).parent / "scripts"
sys.path.insert(0, str(scripts_dir))

# 导入并执行构建
from build import main

if __name__ == "__main__":
    main()
