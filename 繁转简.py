import argparse
from zhconv import convert

def is_index_line(line):
    """检查是否为序号行（纯数字）"""
    return line.strip().isdigit()

def is_timestamp_line(line):
    """检查是否为时间轴行（包含'-->'）"""
    return '-->' in line

def process_line(line):
    """处理单行文本，条件过滤后转换繁体到简体"""
    if is_index_line(line) or is_timestamp_line(line) or line.strip() == '':
        return line
    return convert(line, 'zh-cn')

def main():
    parser = argparse.ArgumentParser(description='将SRT字幕文件中的繁体字转换为简体字')
    parser.add_argument('input', help='输入的SRT文件路径')
    parser.add_argument('output', help='输出的SRT文件路径')
    args = parser.parse_args()

    with open(args.input, 'r', encoding='utf-8') as infile, \
         open(args.output, 'w', encoding='utf-8') as outfile:
        for line in infile:
            processed_line = process_line(line)
            outfile.write(processed_line)

if __name__ == '__main__':
    main()