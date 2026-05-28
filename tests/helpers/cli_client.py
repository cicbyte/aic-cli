"""
AIC CLI 客户端封装 - 用于集成测试
"""
import json
import subprocess
from dataclasses import dataclass
from pathlib import Path
from typing import Optional


@dataclass
class CliResult:
    """CLI 命令执行结果"""
    returncode: int
    stdout: str
    stderr: str

    @property
    def success(self) -> bool:
        return self.returncode == 0

    @property
    def output(self) -> str:
        """合并的输出"""
        stdout = self.stdout or ""
        stderr = self.stderr or ""
        return stdout + stderr

    def json(self) -> Optional[dict]:
        """尝试解析 JSON 输出"""
        try:
            return json.loads(self.stdout)
        except (json.JSONDecodeError, TypeError):
            return None


class AicCliClient:
    """AIC CLI 客户端封装"""

    def __init__(self, binary_path: str):
        self.binary = binary_path

    def run(self, *args: str, input_data: str = None) -> CliResult:
        """执行 CLI 命令"""
        cmd = [self.binary] + list(args)
        result = subprocess.run(
            cmd,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
            input=input_data,
            timeout=30,
        )
        return CliResult(
            returncode=result.returncode,
            stdout=result.stdout,
            stderr=result.stderr,
        )

    # ========== Skill 常用命令 ==========

    def skill_search(self, keyword: str, category_id: int = None, page: int = None, page_size: int = None) -> CliResult:
        """搜索技能"""
        args = ["skill", "search", keyword]
        if category_id:
            args.extend(["-c", str(category_id)])
        if page:
            args.extend(["-p", str(page)])
        if page_size:
            args.extend(["-n", str(page_size)])
        return self.run(*args)

    def skill_add(self, skill_id: int, agent: str = None, output_dir: str = None, mode: str = None) -> CliResult:
        """添加技能到本地"""
        args = ["skill", "add", str(skill_id)]
        if agent:
            args.extend(["--agent", agent])
        if output_dir:
            args.extend(["-o", output_dir])
        if mode:
            args.extend(["-m", mode])
        return self.run(*args)

    def skill_download(self, skill_id: int, output: str = None) -> CliResult:
        """下载技能 ZIP 包"""
        args = ["skill", "download", str(skill_id)]
        if output:
            args.extend(["-o", output])
        return self.run(*args)

    def skill_import(self, path: str, description: str = None, category_id: int = None) -> CliResult:
        """导入本地技能到服务器"""
        args = ["skill", "import", path]
        if description:
            args.extend(["-d", description])
        if category_id:
            args.extend(["-c", str(category_id)])
        return self.run(*args)

    def skill_categories(self) -> CliResult:
        """列出所有分类"""
        return self.run("skill", "categories")

    # ========== Skill Remote 命令 (AI Agent) ==========

    def skill_cat(self, skill_id: int, path: str, output: str = None) -> CliResult:
        """查看远程文件内容"""
        args = ["skill", "remote", "cat", str(skill_id), path]
        if output:
            args.extend(["-o", output])
        return self.run(*args)

    def skill_tree(self, skill_id: int) -> CliResult:
        """显示远程技能文件树"""
        return self.run("skill", "remote", "tree", str(skill_id))

    def skill_patch(self, skill_id: int, path: str = None, old: str = None, new: str = None,
                    line: str = None, replace_all: bool = False, stdin: str = None,
                    batch: str = None) -> CliResult:
        """增量编辑技能文件

        Args:
            skill_id: 技能 ID
            path: 文件路径
            old: 要替换的原始文本
            new: 替换后的文本
            line: 行号范围 (如: "5" 或 "20-25")
            replace_all: 替换所有匹配
            stdin: 通过管道传入的 JSON 数据（自动检测）
            batch: 批量操作 JSON 文件
        """
        args = ["skill", "remote", "patch", str(skill_id)]

        # 管道模式：通过 input_data 传入
        if stdin is not None:
            return self.run(*args, input_data=stdin)

        # batch 模式
        if batch:
            args.extend(["--batch", batch])
            return self.run(*args)

        # 普通模式需要 path
        if not path:
            raise ValueError("path 参数是必需的（除非使用管道或 batch 模式）")

        args.extend(["--path", path])
        if old:
            args.extend(["--old", old])
        if new:
            args.extend(["--new", new])
        if line:
            args.extend(["--line", line])
        if replace_all:
            args.append("--replace-all")
        return self.run(*args)

    def skill_validate(self, skill_id: int, strict: bool = False) -> CliResult:
        """校验技能文件"""
        args = ["skill", "remote", "validate", str(skill_id)]
        if strict:
            args.append("--strict")
        return self.run(*args)

    def skill_publish(self, skill_id: int, version: str = None, changelog: str = None) -> CliResult:
        """发布技能"""
        args = ["skill", "remote", "publish", str(skill_id)]
        if version:
            args.extend(["--version", version])
        if changelog:
            args.extend(["--changelog", changelog])
        return self.run(*args)

    def skill_unpublish(self, skill_id: int) -> CliResult:
        """取消发布"""
        return self.run("skill", "remote", "unpublish", str(skill_id))

    def skill_update(self, skill_id: int, name: str = None, desc: str = None,
                     tags: str = None, category: int = None) -> CliResult:
        """更新技能元数据"""
        args = ["skill", "remote", "update", str(skill_id)]
        if name:
            args.extend(["--name", name])
        if desc:
            args.extend(["--desc", desc])
        if tags:
            args.extend(["--tags", tags])
        if category:
            args.extend(["--category", str(category)])
        return self.run(*args)

    # ========== Server 命令 ==========

    def server_status(self) -> CliResult:
        """查看服务器连接状态"""
        return self.run("server", "status")

    # ========== Version 命令 ==========

    def version(self) -> CliResult:
        """显示版本信息"""
        return self.run("version")
