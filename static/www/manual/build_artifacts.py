#!/usr/bin/env python3
"""Generate manual.md and index.json from sections/*.tex."""
import glob, os, re, json, subprocess

sections_dir = 'sections'
output_md = 'manual.md'
output_json = 'index.json'

chapters = sorted(glob.glob(os.path.join(sections_dir, '*.tex')))
if not chapters:
    raise SystemExit('No .tex files found in sections/')

def convert_list(match):
    content = match.group(0)
    items = re.findall(r'\\item\s+(.*)', content)
    result = []
    for i, item in enumerate(items, 1):
        item = re.sub(r'\\\\[\n\r]*', ' ', item)
        item = re.sub(r'\\[a-zA-Z]+\{([^}]*)\}', r'\1', item)
        item = re.sub(r'\\[a-zA-Z]+', '', item).strip()
        result.append('- ' + item)
    return '\n' + '\n'.join(result) + '\n'

def clean_latex(text):
    lines = []
    for line in text.split('\n'):
        stripped = line.lstrip()
        if stripped.startswith('%') and not line.startswith('%%'):
            continue
        lines.append(line)
    text = '\n'.join(lines)
    text = re.sub(r'\\section\{([^}]*)\}', r'\n## \1\n', text)
    text = re.sub(r'\\subsection\{([^}]*)\}', r'\n### \1\n', text)
    text = re.sub(r'\\subsubsection\{([^}]*)\}', r'\n#### \1\n', text)
    text = re.sub(r'\\begin\{(itemize|enumerate)\}.*?\\end\{\1\}', convert_list, text, flags=re.DOTALL)
    text = re.sub(r'\\textbf\{([^}]*)\}', r'**\1**', text)
    text = re.sub(r'\\textit\{([^}]*)\}', r'*\1*', text)
    text = re.sub(r'\\url\{([^}]*)\}', lambda m: '[' + m.group(1) + ']', text)
    text = re.sub(r'\\\\[\n\r]*', '\n', text)
    text = re.sub(r'\\[a-zA-Z]+\{([^}]*)\}', r'\1', text)
    text = re.sub(r'\\[a-zA-Z]+', '', text)
    text = re.sub(r'\n{4,}', '\n\n\n', text)
    return text.strip()

# Read PDF page count
result = subprocess.run(['pdfinfo', 'lab-manual.pdf'], capture_output=True, text=True)
page_count = 0
for line in result.stdout.split('\n'):
    if 'Pages:' in line:
        page_count = int(line.split(':')[1].strip())
        break

print(f'Processing {len(chapters)} chapters...')
output_parts = []
sections_info = []
pages_per_section = page_count / len(chapters) if chapters else 1

for i, tex_file in enumerate(chapters):
    filename = os.path.basename(tex_file)
    chapter_title = filename.replace('.tex', '').replace('-', ' ').title()
    with open(tex_file, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    clean = clean_latex(content)
    if clean:
        output_parts.append('\n\n' + '='*60 + '\n## ' + chapter_title + '\n' + '='*60 + '\n' + clean)
    sections_info.append({
        'title': chapter_title,
        'file': filename,
        'pageRange': [int(i * pages_per_section) + 1, min(int((i+1) * pages_per_section), page_count)]
    })

markdown = '# Bitsmasher Lab Operations Manual\n\nPlaintext Markdown export compiled from ' + str(len(chapters)) + ' operational chapters.\n'
output_parts.insert(0, markdown)
final_md = '\n'.join(output_parts)

with open(output_md, 'w', encoding='utf-8') as f:
    f.write(final_md)

sections_json = {
    'title': 'Bitsmasher Lab Operations Manual',
    'pageCount': page_count,
    'sections': sections_info,
    'generatedAt': '2026-08-27T21:35:00Z'
}
with open(output_json, 'w', encoding='utf-8') as f:
    json.dump(sections_json, f, indent=2)

print(f'Written {len(final_md)} chars to {output_md}')
print(f'Written JSON with {len(sections_info)} sections to {output_json}')
print(f'PDF: {page_count} pages')
