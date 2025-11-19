#!/usr/bin/env python3
"""
跳转链接测试脚本
"""
import re
import sys
from pathlib import Path

DOCS_DIR = Path(__file__).resolve().parent.parent / "docs"

def test_jump_links(file_path):
    """测试跳转链接是否正常"""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # 提取目录链接
    toc_links = re.findall(r'\[([^\]]+)\]\((#[^\)]+)\)', content)

    # 提取HTML锚点
    html_anchors = re.findall(r'<a name="([^"]+)"></a>', content)

    print("=== 跳转链接测试 ===")
    print(f"目录链接数量: {len(toc_links)}")
    print(f"HTML锚点数量: {len(html_anchors)}")

    # 检查链接是否匹配
    broken_links = []
    working_links = []

    for link_text, link_href in toc_links:
        # 移除#号
        anchor_id = link_href[1:]

        if anchor_id in html_anchors:
            working_links.append((link_text, link_href))
        else:
            broken_links.append((link_text, link_href))

    print(f"\n✅ 正常工作的链接: {len(working_links)}")
    for link_text, link_href in working_links[:10]:  # 显示前10个
        print(f"  - {link_text} -> {link_href}")

    if broken_links:
        print(f"\n❌ 损坏的链接: {len(broken_links)}")
        for link_text, link_href in broken_links:
            print(f"  - {link_text} -> {link_href}")
    else:
        print(f"\n🎉 所有链接都正常工作!")

    # 检查返回链接
    return_links = re.findall(r'\[🔝 返回目录\]\((#[^\)]+)\)', content)
    print(f"\n🔄 返回目录链接数量: {len(return_links)}")

    if return_links and return_links[0] == '#📑-目录':
        print("✅ 返回目录链接正确")
    else:
        print(f"❌ 返回目录链接可能有问题: {return_links}")

    return len(broken_links) == 0

if __name__ == '__main__':
    file_path = DOCS_DIR / "API文档.md"
    success = test_jump_links(str(file_path))
    sys.exit(0 if success else 1)
