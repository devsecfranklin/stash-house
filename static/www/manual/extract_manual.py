import glob, os, re, sys

sections_dir = "/home/franklin/workspace/website/static/www/manual/sections"
output_path = "/home/franklin/workspace/website/static/www/manual/manual.md"

chapters = sorted(glob.glob(os.path.join(sections_dir, "*.tex")))
if not chapters:
    print("No .tex files found in sections/", file=sys.stderr)
    sys.exit(1)

def convert_list(match):
    content = match.group(0)
    items = re.findall(r"\\item\s+(.*)", content)
    result = []
    for i, item in enumerate(items, 1):
        item = re.sub(r"\\\\[\n\r]*", " ", item)
        item = re.sub(r"\\[a-zA-Z]+\{([^}]*)\}", r"\1", item)
        item = re.sub(r"\\[a-zA-Z]+", "", item).strip()
        result.append("- " + item)
    return "\n" + "\n".join(result) + "\n"

def clean_latex(text):
    lines = []
    for line in text.split("\n"):
        stripped = line.lstrip()
        if stripped.startswith("%") and not line.startswith("%%"):
            continue
        lines.append(line)
    text = "\n".join(lines)
    text = re.sub(r"\\section\{([^}]*)\}", r"\n## \1\n", text)
    text = re.sub(r"\\subsection\{([^}]*)\}", r"\n### \1\n", text)
    text = re.sub(r"\\subsubsection\{([^}]*)\}", r"\n#### \1\n", text)
    text = re.sub(r"\\begin\{(itemize|enumerate)\}.*?\\end\{\1\}", convert_list, text, flags=re.DOTALL)
    text = re.sub(r"\\textbf\{([^}]*)\}", r"**\1**", text)
    text = re.sub(r"\\textit\{([^}]*)\}", r"*\1*", text)
    text = re.sub(r"\\url\{([^}]*)\}", r"[`<http://\1>`]", text)
    text = re.sub(r"\\\\[\n\r]*", "\n", text)
    text = re.sub(r"\\[a-zA-Z]+\{([^}]*)\}", r"\1", text)
    text = re.sub(r"\\[a-zA-Z]+", "", text)
    text = re.sub(r"\n{4,}", "\n\n\n", text)
    return text.strip()

print(f"Processing {len(chapters)} chapters...")
output_parts = []

for tex_file in chapters:
    filename = os.path.basename(tex_file)
    chapter_title = filename.replace(".tex", "").replace("-", " ").title()
    with open(tex_file, "r", encoding="utf-8", errors="replace") as f:
        content = f.read()
    clean = clean_latex(content)
    if clean:
        output_parts.append("\n\n" + "="*60 + "\n## " + chapter_title + "\n" + "="*60 + "\n" + clean)

markdown = "# Bitsmasher Lab Operations Manual\n\nPlaintext Markdown export compiled from 22 operational chapters.\n"
output_parts.insert(0, markdown)
final = "\n".join(output_parts)

with open(output_path, "w", encoding="utf-8") as f:
    f.write(final)

chars = len(final)
print(f"Written {chars} chars to {output_path}")
