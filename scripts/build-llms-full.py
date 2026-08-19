#!/usr/bin/env python3
"""Build llms-full.txt from blog source files for agent ingestion."""
import os, glob, datetime, html as htmlmod, re

SITE_SRC = "/mnt/lab_franklin_1/workspace/website/static/www"
BLOGDIR = f"{SITE_SRC}/research/blog"
DOMAIN = "https://www.bitsmasher.net"

def strip_html_to_text(filepath):
    with open(filepath, 'r') as f:
        raw = f.read()
    clean = re.sub(r'<script[^>]*>.*?</script>', '', raw, flags=re.DOTALL)
    clean = re.sub(r'<style[^>]*>.*?</style>', '', clean, flags=re.DOTALL)
    clean = re.sub(r'<br\s*/?>', '\n', clean)
    clean = re.sub(r'</?(?:p|div|h[1-6]|li|tr|article|section)[^>]*>', '\n', clean, flags=re.IGNORECASE)
    clean = re.sub(r'<[^>]+>', '', clean)
    clean = htmlmod.unescape(clean)
    lines = [l.strip() for l in clean.split('\n')]
    return '\n'.join([l for l in lines if l])

all_files = sorted(glob.glob(os.path.join(BLOGDIR, '*')), key=lambda x: os.path.basename(x))
basename_set = {}
for fp in all_files:
    base = os.path.basename(fp)
    stem = base.rsplit('.', 1)[0]
    ext = base.rsplit('.', 1)[1].lower() if '.' in base else ''
    if ext == 'md':
        with open(fp, 'r') as f:
            first = f.readline().strip()
        if first.startswith('<!DOCTYPE'):
            continue
    basename_set[stem] = fp

all_files_deduped = sorted(basename_set.values(), key=lambda x: os.path.basename(x))

out_lines = [
    "# bitsmasher.net — Full Site Content for Agent Ingestion",
    "",
    f"Generated: {datetime.datetime.now(datetime.timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')}",
    "This file contains all blog posts from bitsmasher.net research blog.",
    "Agents should ingest this as their primary site knowledge source.",
    "",
    "---",
    "",
]

for filepath in all_files_deduped:
    fname = os.path.basename(filepath)
    ext = fname.rsplit('.', 1)[-1].lower()
    url = f"{DOMAIN}/research/blog/{fname}"
    out_lines.append(f"## Source: {url}")
    out_lines.append("")
    if ext == 'md':
        with open(filepath, 'r') as f:
            out_lines.append(f.read().strip())
    else:
        out_lines.append(strip_html_to_text(filepath))
    out_lines.extend(["", "---", ""])

output = '\n'.join(out_lines) + '\n'
outfile = os.path.join(SITE_SRC, 'llms-full.txt')
with open(outfile, 'w') as f:
    f.write(output)
print(f"Built {len(all_files_deduped)} posts -> {outfile} ({len(output):,} bytes)")
