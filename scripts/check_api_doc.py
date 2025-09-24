#!/usr/bin/env python3
"""
API文档跳转系统检查和重构脚本
"""
import re
import sys

def extract_headings(file_path):
    """提取所有标题及其位置"""
    headings = []
    with open(file_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()

    for i, line in enumerate(lines, 1):
        # 匹配 ## 和 ### 级别的标题
        if line.startswith('## '):
            level = 2
            title = line[3:].strip()
            headings.append({'level': level, 'title': title, 'line': i, 'raw': line})
        elif line.startswith('### '):
            level = 3
            title = line[4:].strip()
            headings.append({'level': level, 'title': title, 'line': i, 'raw': line})

    return headings

def generate_anchor_id(title):
    """生成锚点ID - 移除HTML标签和特殊字符"""
    # 移除HTML标签
    title = re.sub(r'<[^>]+>', '', title)
    # 移除emoji和特殊字符，只保留中文、英文、数字和空格
    title = re.sub(r'[^\u4e00-\u9fff\u3400-\u4dbf\a-zA-Z0-9\s]', '', title)
    # 替换空格为连字符
    anchor = title.replace(' ', '-').lower()
    return anchor

def check_current_links(file_path):
    """检查当前的链接和锚点"""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    # 查找目录中的链接
    toc_links = re.findall(r'\[([^\]]+)\]\((#[^\)]+)\)', content)

    # 查找HTML锚点
    html_anchors = re.findall(r'<a name="([^"]+)"></a>', content)

    # 查找Markdown链接
    md_links = re.findall(r'\[🔝 返回目录\]\((#[^\)]+)\)', content)

    print("=== 当前链接分析 ===")
    print(f"目录链接数量: {len(toc_links)}")
    print(f"HTML锚点数量: {len(html_anchors)}")
    print(f"返回链接数量: {len(md_links)}")

    print("\n目录链接:")
    for link_text, link_href in toc_links[:10]:  # 只显示前10个
        print(f"  - {link_text} -> {link_href}")

    print("\nHTML锚点:")
    for anchor in html_anchors[:10]:  # 只显示前10个
        print(f"  - {anchor}")

    return toc_links, html_anchors, md_links

def generate_new_document(file_path):
    """生成新的文档结构"""
    with open(file_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()

    new_lines = []
    heading_map = {}  # 标题到锚点的映射

    for line in lines:
        new_line = line

        # 处理 ## 级别的标题
        if line.startswith('## '):
            title = line[3:].strip()
            # 移除已有的HTML标签
            clean_title = re.sub(r'<[^>]+>', '', title)
            anchor = generate_anchor_id(clean_title)

            # 生成新的标题行（纯Markdown格式）
            new_line = f"## {title}\n"
            # 记录映射关系
            heading_map[clean_title] = anchor

        # 处理 ### 级别的标题
        elif line.startswith('### '):
            title = line[4:].strip()
            # 移除已有的HTML标签
            clean_title = re.sub(r'<[^>]+>', '', title)
            anchor = generate_anchor_id(clean_title)

            # 生成新的标题行（纯Markdown格式）
            new_line = f"### {title}\n"
            # 记录映射关系
            heading_map[clean_title] = anchor

        new_lines.append(new_line)

    # 第二步：处理链接
    final_lines = []
    for line in new_lines:
        # 处理目录链接
        if line.strip().startswith('- [') and '](' in line:
            # 提取链接文本和锚点
            match = re.search(r'\[([^\]]+)\]\((#[^\)]+)\)', line)
            if match:
                link_text, old_href = match.groups()
                # 生成新的锚点
                clean_text = re.sub(r'<[^>]+>', '', link_text)
                clean_text = re.sub(r'[^\u4e00-\u9fff\u3400-\u4dbf\a-zA-Z0-9\s]', '', clean_text)
                new_href = f"#{clean_text.replace(' ', '-').lower()}"

                # 特殊处理admin相关的链接
                if 'admin' in old_href.lower() and '文章管理' in link_text:
                    new_href = '#admin-文章管理'
                elif 'admin' in old_href.lower() and '分类管理' in link_text:
                    new_href = '#admin-分类管理'

                new_line = line.replace(f"[{link_text}]{old_href}", f"[{link_text}]{new_href}")
                final_lines.append(new_line)
            else:
                final_lines.append(line)

        # 处理返回目录链接
        elif '[🔝 返回目录]' in line:
            final_lines.append('[🔝 返回目录](#📑-目录)\n')
        else:
            final_lines.append(line)

    return final_lines, heading_map

def main():
    file_path = '/Volumes/R0sORICO/work_dir/r0website_server/API文档.md'

    print("=== API文档跳转系统检查工具 ===\n")

    # 1. 分析当前文档结构
    headings = extract_headings(file_path)
    print(f"发现 {len(headings)} 个标题:")
    for h in headings[:15]:  # 显示前15个
        print(f"  第{h['line']}行: {'#' * h['level']} {h['title']}")

    print("\n" + "="*50 + "\n")

    # 2. 检查当前链接
    toc_links, html_anchors, md_links = check_current_links(file_path)

    print("\n" + "="*50 + "\n")

    # 3. 生成新的文档结构
    print("正在生成新的文档结构...")
    new_content, heading_map = generate_new_document(file_path)

    print("\n标题到锚点映射:")
    for title, anchor in list(heading_map.items())[:10]:
        print(f"  '{title}' -> '#{anchor}'")

    # 4. 保存新文档
    output_path = '/Volumes/R0sORICO/work_dir/r0website_server/API文档_新.md'
    with open(output_path, 'w', encoding='utf-8') as f:
        f.writelines(new_content)

    print(f"\n✅ 新文档已生成: {output_path}")
    print(f"文档总行数: {len(new_content)}")

    # 5. 生成简单的跳转测试
    print("\n=== 跳转测试建议 ===")
    print("请手动测试以下跳转链接:")
    print("1. 点击目录中的'用户管理'链接")
    print("2. 点击'图床管理'下的子链接")
    print("3. 点击任意'🔝 返回目录'链接")
    print("4. 验证所有链接都能正确跳转")

if __name__ == '__main__':
    main()