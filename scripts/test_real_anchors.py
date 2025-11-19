#!/usr/bin/env python3
"""
测试真实的Markdown锚点链接
"""
import re
from pathlib import Path

DOCS_DIR = Path(__file__).resolve().parent.parent / "docs"

def test_real_anchors():
    """测试真实的Markdown锚点链接"""
    file_path = DOCS_DIR / "API文档.md"

    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    print("=== 真实Markdown锚点链接测试 ===")
    print(f"文件路径: {file_path}")

    # 提取目录链接
    toc_links = re.findall(r'\[([^\]]+)\]\((#[^\)]+)\)', content)

    # 提取所有标题
    all_titles = []
    lines = content.split('\n')
    for line_num, line in enumerate(lines, 1):
        if line.startswith('## '):
            title = line[3:].strip()
            all_titles.append((2, title, line_num))
        elif line.startswith('### '):
            title = line[4:].strip()
            all_titles.append((3, title, line_num))
        elif line.startswith('#### '):
            title = line[5:].strip()
            all_titles.append((4, title, line_num))

    print(f"找到 {len(all_titles)} 个标题")
    print(f"找到 {len(toc_links)} 个目录链接")

    # 模拟Markdown锚点ID生成规则
    def generate_anchor_id(title):
        """生成Markdown锚点ID"""
        # 转换为小写
        anchor = title.lower()
        # 替换空格和特殊字符为连字符
        anchor = re.sub(r'[^\w\u4e00-\u9fff\u3040-\u309f\u30a0-\u30ff\- ]', '', anchor)
        # 替换空格为连字符
        anchor = re.sub(r'\s+', '-', anchor)
        # 移除重复的连字符
        anchor = re.sub(r'-+', '-', anchor)
        # 移除首尾连字符
        anchor = anchor.strip('-')
        return anchor

    print("\n=== 生成的锚点ID ===")
    for level, title, line_num in all_titles:
        anchor_id = generate_anchor_id(title)
        print(f"第{line_num}行: {title} -> #{anchor_id}")

    # 检查每个链接
    working_links = []
    broken_links = []

    for link_text, link_href in toc_links:
        target_anchor = link_href[1:]  # 移除#号

        # 检查是否有标题生成这个锚点ID
        found = False
        for level, title, line_num in all_titles:
            expected_anchor = generate_anchor_id(title)
            if target_anchor == expected_anchor:
                found = True
                working_links.append((link_text, link_href, title))
                break

        if not found:
            broken_links.append((link_text, link_href))

    print(f"\n=== 测试结果 ===")
    print(f"✅ 正常工作的链接: {len(working_links)}")
    for link_text, link_href, target_title in working_links:
        print(f"  - [{link_text}] -> {link_href} (目标: {target_title})")

    if broken_links:
        print(f"\n❌ 损坏的链接: {len(broken_links)}")
        for link_text, link_href in broken_links:
            print(f"  - [{link_text}] -> {link_href}")

    return len(broken_links) == 0

if __name__ == '__main__':
    success = test_real_anchors()
    print(f"\n结果: {'✅ 所有链接正常' if success else '❌ 存在损坏的链接'}")
