#!/usr/bin/env python3
"""
测试纯Markdown格式的跳转链接
"""
import re
from pathlib import Path

DOCS_DIR = Path(__file__).resolve().parent.parent / "docs"

def test_pure_md_links():
    """测试纯Markdown版本的跳转链接"""
    file_path = DOCS_DIR / "API文档.md"

    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # 提取目录链接
    toc_links = re.findall(r'\[([^\]]+)\]\((#[^\)]+)\)', content)

    print("=== 纯Markdown版本跳转链接测试 ===")
    print(f"文件路径: {file_path}")
    print(f"目录链接数量: {len(toc_links)}")

    # 检查每个链接是否有对应的标题
    missing_targets = []
    working_links = []

    for link_text, link_href in toc_links:
        anchor_id = link_href[1:]  # 移除#号

        # 查找对应的标题
        # 查找 ## 标题
        title_pattern = f'^##\\s+{re.escape(anchor_id)}\\s*$'
        # 查找 ### 标题
        subtitle_pattern = f'^###\\s+{re.escape(anchor_id)}\\s*$'

        if re.search(title_pattern, content, re.MULTILINE) or re.search(subtitle_pattern, content, re.MULTILINE):
            working_links.append((link_text, link_href))
        else:
            missing_targets.append((link_text, link_href))
            print(f"❌ [{link_text}] -> {link_href} (未找到对应标题)")

    print(f"\n✅ 正常工作的链接: {len(working_links)}")
    for link_text, link_href in working_links:
        print(f"  - [{link_text}] -> {link_href}")

    if missing_targets:
        print(f"\n❌ 缺失目标的链接: {len(missing_targets)}")
        for link_text, link_href in missing_targets:
            print(f"  - [{link_text}] -> {link_href}")

    # 检查所有标题
    print(f"\n📋 发现的所有标题:")
    titles = re.findall(r'^##\s+(.+)$', content, re.MULTILINE)
    subtitles = re.findall(r'^###\s+(.+)$', content, re.MULTILINE)

    print("## 级别标题:")
    for i, title in enumerate(titles, 1):
        print(f"  {i}. {title}")

    print("### 级别标题:")
    for i, title in enumerate(subtitles, 1):
        print(f"  {i}. {title}")

    return len(missing_targets) == 0

if __name__ == '__main__':
    success = test_pure_md_links()
    print(f"\n结果: {'✅ 所有链接正常' if success else '❌ 存在问题的链接'}")
