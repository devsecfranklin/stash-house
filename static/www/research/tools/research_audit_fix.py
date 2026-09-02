#!/usr/bin/env python3
"""
Deterministic audit-and-fix for the wonderland research blog.
All work runs in-process — no LLM context, no multi-turn edits.

Phases:
  P1: Audit inline <style> blocks across all files
  P2: Strip inline CSS from blog posts that already link canonical stylesheet
  P3: Fix index.html structure (h1 position, robots meta, JSON-LD count)
  P4: Append research-specific CSS rules to canonical stylesheet if missing
  P5: File permissions fixup
"""

import os
import re
import json
import sys
from pathlib import Path

BASE = Path("/home/franklin/workspace/website/static/www/research")
BLOG = BASE / "blog"
CSS_PATH = Path("/home/franklin/workspace/website/static/www/css/www.bitsmasher.net.css")

audit_results = {
    "inline_style_files": [],
    "stripped_count": 0,
    "css_rules_missing": [],
    "index_html_issues": [],
    "errors": [],
}


def strip_inline_style(filepath):
    content = filepath.read_text()
    if '<style>' not in content:
        return False
    new_content = re.sub(r'\s*<style>\s*\n.*?\n\s*</style>\s*', '', content, flags=re.DOTALL)
    filepath.write_text(new_content)
    print(f"  STRIPPED: {filepath.relative_to(BASE)}")
    audit_results["stripped_count"] += 1
    return True


def fix_index_html():
    target = BASE / "index.html"
    content = target.read_text()

    issues_found = []

    # Fix H1 position — must be above first blog-post div inside <main>
    h1_pos = content.find('<h1>')
    first_blog_div = content.find('<div class="blog-post">')
    if h1_pos > first_blog_div and h1_pos >= 0 and first_blog_div >= 0:
        h1_match = re.search(r'<h1>[^<]*</h1>', content[h1_pos:])
        if h1_match:
            h1_text = h1_match.group()
            after_main = re.search(r'<main(?:\s+[^>]*)?>', content)
            if after_main:
                insert_pos = after_main.end()
                clean_content = content[:insert_pos] + '\n' + h1_text + '\n' + content[insert_pos:]
                # Remove old H1 occurrence
                old_h1_end = content.find('</h1>', h1_pos) + 6
                after_new = clean_content.rfind(h1_text)
                if after_new > -1 and after_new != insert_pos:
                    pre_old = clean_content[:after_new]
                    post_old = clean_content[old_h1_end:]
                    clean_content = pre_old + post_old
                target.write_text(clean_content)
                content = clean_content
                issues_found.append("Moved H1 to top of <main>")

    # Add robots meta if missing
    if 'name="robots"' not in content:
        content_fixed = re.sub(
            r'(<title>[^<]*</title>)',
            r'\1\n    <meta name="robots" content="follow,index,max-snippet:-1">',
            content, count=1
        )
        target.write_text(content_fixed)
        issues_found.append("Added robots meta")

    # Add feed.json alternate link if missing
    if '/research/feed.json' not in content:
        content = target.read_text()
        content_fixed = re.sub(
            r'(<title>[^<]*</title>)',
            r'\1\n    <link rel="alternate" type="application/feed+json" href="/research/feed.json">',
            content, count=1
        )
        target.write_text(content_fixed)
        issues_found.append("Added feed.json alternate link")

    # Fix JSON-LD numberOfPosts to match actual div count
    content = target.read_text()
    html_div_count = content.count('class="blog-post"')
    jsonld_match = re.search(r'"numberOfPosts":\s*(\d+)', content)
    if jsonld_match:
        jsonld_num = int(jsonld_match.group(1))
        if jsonld_num != html_div_count:
            content = re.sub(
                r'"numberOfPosts":\s*\d+',
                f'"numberOfPosts": {html_div_count}',
                content
            )
            target.write_text(content)
            issues_found.append(f"Corrected JSON-LD numberOfPosts to {html_div_count}")

    # Add canonical link if missing
    content = target.read_text()
    if 'rel="canonical"' not in content:
        content_fixed = re.sub(
            r'(<title>[^<]*</title>)',
            r'\1\n    <link rel="canonical" href="https://www.bitsmasher.net/research/">',
            content, count=1
        )
        target.write_text(content_fixed)
        issues_found.append("Added canonical link")

    # Add favicon if missing
    content = target.read_text()
    if 'favicon' not in content.lower():
        content_fixed = re.sub(
            r'(<title>[^<]*</title>)',
            r'\1\n    <link rel="icon" type="image/x-icon" href="https://www.bitsmasher.net/images/favicon.ico">',
            content, count=1
        )
        target.write_text(content_fixed)
        issues_found.append("Added favicon link")

    audit_results["index_html_issues"] = issues_found


def check_and_append_css():
    if not CSS_PATH.exists():
        print(f"  ERROR: Canonical stylesheet not found at {CSS_PATH}")
        audit_results["errors"].append("Canonical stylesheet missing")
        return

    current = CSS_PATH.read_text()

    needed_rules = ['.research-content', '.page-title', '.page-footer', '.blog-post']
    missing = [r for r in needed_rules if r not in current]

    if not missing:
        print("  OK: Required research CSS rules already present in canonical stylesheet")
        return

    # Check if any subpages have inline copies that we need to preserve the content of
    css_to_append = """\n/* ── Research Blog Specific Styles ─────────────────────────────── */

.research-content {
  text-align: left;
  max-width: 900px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2.5rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid rgba(26,115,232,0.2);
}

.page-title {
  color: #1a73e8;
  font-size: 2rem;
  margin: 0.3rem 0;
}

.page-subtitle {
  color: #9e9e9e;
  font-size: 1rem;
  margin-bottom: 0.5rem;
}

h2, h3 {
  color: #64b5f6;
}

h4 {
  color: #90caf9;
}

p {
  line-height: 1.7;
}

code {
  background: rgba(26,115,232,0.15);
  padding: 2px 6px;
  border-radius: 3px;
  color: #82b1ff;
}

.back-link {
  margin-bottom: 2rem;
  font-size: 0.9rem;
}

.back-link a {
  color: #4fc3f7;
  text-decoration: none;
}

.back-link a:hover {
  text-decoration: underline;
}

.page-footer {
  margin-top: 3rem;
  padding-top: 1.5rem;
  border-top: 1px solid rgba(26,115,232,0.2);
  font-size: 0.8rem;
  color: #9e9e9e;
  text-align: center;
}

.blog-post {
  margin-bottom: 1.5rem;
  padding: 1rem 0;
  border-bottom: 1px solid rgba(26,115,232,0.1);
}

.blog-post article {
  background: transparent;
  padding: 0;
}

.blog-post h2 {
  font-size: 1.2rem;
  margin: 0.3rem 0;
}

.blog-post h2 a {
  color: var(--link-color, #3399ff);
  text-decoration: underline;
}

.blog-post h2 a:hover {
  color: var(--link-hover, #00ffcc);
  text-decoration: none;
}

.blog-post time {
  color: var(--text-muted, #94a3b8);
  font-size: 0.85rem;
  margin-bottom: 0.5rem;
  display: block;
}

.blog-post p {
  margin: 0.3rem 0 0 0;
  color: var(--text-main, #e2e8f0);
  font-size: 0.95rem;
}
"""
    with open(CSS_PATH, 'a') as f:
        f.write(css_to_append)
    audit_results["css_rules_missing"] = missing
    print(f"  APPENDED {len(missing)} missing CSS rule groups to canonical stylesheet")


def fix_permissions():
    for root, dirs, files in os.walk(BASE):
        for d in dirs:
            p = Path(root) / d
            try:
                os.chmod(p, 0o755)
            except OSError:
                pass
        for f in files:
            p = Path(root) / f
            try:
                os.chmod(p, 0o644)
            except OSError:
                pass
    print(f"  Fixed permissions on {BASE}")


def main():
    print("=" * 60)
    print("RESEARCH BLOG AUDIT & HARDENING — DETERMINISTIC FIX")
    print("=" * 60)

    # Phase 1: Audit inline styles
    print("\n[Phase 1] Auditing inline <style> blocks...")
    for post in sorted(BLOG.glob("*.html")):
        content = post.read_text()
        if '<style>' in content:
            audit_results["inline_style_files"].append(str(post.relative_to(BASE)))
            print(f"  FOUND: {post.name}")

    # Phase 2: Strip inline CSS from files with canonical stylesheet
    print("\n[Phase 2] Stripping inline CSS (where canonical stylesheet exists)...")
    stripped = 0
    for post in sorted(BLOG.glob("*.html")):
        if '<style>' in post.read_text():
            content = post.read_text()
            if '/css/www.bitsmasher.net.css' in content:
                if strip_inline_style(post):
                    stripped += 1

    # Phase 3: Fix index.html
    print("\n[Phase 3] Fixing index.html structure...")
    fix_index_html()

    # Phase 4: Ensure canonical CSS has research rules
    print("\n[Phase 4] Checking canonical stylesheet...")
    check_and_append_css()

    # Phase 5: Permissions
    print("\n[Phase 5] Fixing permissions...")
    fix_permissions()

    # Summary
    print("\n" + "=" * 60)
    print("SUMMARY")
    print("=" * 60)
    print(f"Inline style files found: {len(audit_results['inline_style_files'])}")
    print(f"CSS blocks stripped: {audit_results['stripped_count']}")
    if audit_results['index_html_issues']:
        print("index.html fixes:")
        for issue in audit_results['index_html_issues']:
            print(f"  - {issue}")
    if audit_results['errors']:
        print("\nERRORS:")
        for e in audit_results['errors']:
            print(f"  x {e}")
    print("=" * 60)


if __name__ == '__main__':
    main()
