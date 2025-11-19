#!/usr/bin/env python3
"""
测试新的API文档跳转链接
"""
import re
import sys
from pathlib import Path

DOCS_DIR = Path(__file__).resolve().parent.parent / "docs"

def test_new_jump_links():
    """测试新的API文档跳转链接"""
    file_path = DOCS_DIR / "API文档_新.md"

    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # 提取目录链接
    toc_links = re.findall(r'\[([^\]]+)\]\((#[^\)]+)\)', content)

    # 提取HTML锚点
    html_anchors = re.findall(r'<a name="([^"]+)"></a>', content)

    print("=== 新API文档跳转链接测试 ===")
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
    for link_text, link_href in working_links:
        print(f"  - {link_text} -> {link_href}")

    if broken_links:
        print(f"\n❌ 损坏的链接: {len(broken_links)}")
        for link_text, link_href in broken_links:
            print(f"  - {link_text} -> {link_href}")
        return False
    else:
        print(f"\n🎉 所有链接都正常工作!")
        return True

if __name__ == '__main__':
    success = test_new_jump_links()
    sys.exit(0 if success else 1)
