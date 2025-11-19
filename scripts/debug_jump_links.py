#!/usr/bin/env python3
"""
详细调试跳转链接问题
"""
import re
from pathlib import Path

DOCS_DIR = Path(__file__).resolve().parent.parent / "docs"

def debug_jump_links():
    """详细分析跳转链接问题"""
    file_path = DOCS_DIR / "API文档.md"

    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # 提取目录链接
    toc_links = re.findall(r'\[([^\]]+)\]\((#[^\)]+)\)', content)

    # 提取HTML锚点
    html_anchors = re.findall(r'<a name="([^"]+)"></a>', content)

    print("=== 详细跳转链接调试 ===")
    print(f"文件路径: {file_path}")
    print(f"文件大小: {len(content)} 字符")
    print(f"总行数: {len(content.splitlines())}")

    print(f"\n📋 目录链接发现: {len(toc_links)} 个")
    for i, (text, href) in enumerate(toc_links, 1):
        print(f"  {i:2d}. [{text}] -> {href}")

    print(f"\n⚓ HTML锚点发现: {len(html_anchors)} 个")
    for i, anchor in enumerate(html_anchors, 1):
        print(f"  {i:2d}. <a name=\"{anchor}\"></a>")

    print(f"\n🔗 链接匹配分析:")
    matched = 0
    unmatched = []

    for text, href in toc_links:
        anchor_id = href[1:]  # 移除#号
        if anchor_id in html_anchors:
            matched += 1
            print(f"  ✅ [{text}] -> {href} (匹配)")
        else:
            unmatched.append((text, href))
            print(f"  ❌ [{text}] -> {href} (未匹配)")

    print(f"\n📊 匹配结果: {matched}/{len(toc_links)} 个链接匹配")

    if unmatched:
        print(f"\n🔍 未匹配的链接:")
        for text, href in unmatched:
            print(f"    - [{text}] -> {href}")
            # 查找相似的锚点
            anchor_id = href[1:]
            similar = [a for a in html_anchors if anchor_id in a or a in anchor_id]
            if similar:
                print(f"      💡 可能的匹配: {similar}")

    # 检查特殊字符问题
    print(f"\n🔤 特殊字符分析:")
    emoji_anchors = [a for a in html_anchors if any(ord(c) > 127 for c in a)]
    if emoji_anchors:
        print(f"  发现emoji/特殊字符锚点: {emoji_anchors}")

    # 检查返回链接
    return_links = re.findall(r'\[🔝 返回目录\]\((#[^\)]+)\)', content)
    print(f"\n🔄 返回目录链接: {len(return_links)} 个")
    for href in return_links:
        anchor_id = href[1:]
        if anchor_id in html_anchors:
            print(f"  ✅ 返回目录 -> {href} (匹配)")
        else:
            print(f"  ❌ 返回目录 -> {href} (未匹配)")

    print(f"\n💡 建议:")
    print(f"  1. 确保使用的Markdown查看器支持HTML锚点")
    print(f"  2. 检查是否有JavaScript阻止了默认跳转行为")
    print(f"  3. 尝试在不同的Markdown查看器中测试")
    print(f"  4. 检查文件编码是否为UTF-8")

if __name__ == '__main__':
    debug_jump_links()
