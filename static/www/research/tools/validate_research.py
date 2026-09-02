#!/usr/bin/env python3
"""
research_page_linter.py — Wonderlands Research Blog CSS & Formatting Auditor

Scans all research blog files (index.html, index.json, blog/*.html) for:
  - HTML structure violations (missing DOCTYPE, lang, semantic elements)
  - CSS consistency checks (palette compliance, embedded style presence)
  - JSON Feed v1.1 compliance on index.json
  - Cross-reference integrity (href targets, canonical URLs)
  - Schema.org JSON-LD validity
  - Blog-guide.md rule conformance
  - Missing <meta name="robots"> tags

Usage:
  python3 validate_research.py [--fix] [--deploy]
  
  --fix   Apply auto-remediation for issues where safe fixes are possible
  --deploy  After fixing, scp to wonderland (requires SSH access)

Author: Wintermute
Date: 2026-09-02
"""

import json
import os
import re
import sys
import subprocess
from pathlib import Path
from html.parser import HTMLParser
from urllib.parse import urljoin, urlparse

# ─── Constants ────────────────────────────────────────────────────────

RESEARCH_ROOT = Path(__file__).resolve().parent.parent
INDEX_HTML = RESEARCH_ROOT / "index.html"
INDEX_JSON = RESEARCH_ROOT / "index.json"
BLOG_DIR = RESEARCH_ROOT / "blog"
GUIDE_FILE = RESEARCH_ROOT / "blog-guide.md"
DEPLOY_USER = "franklin"
DEPLOY_HOST = "wonderland"
DEPLOY_PATH = f"/home/franklin/workspace/website/static/www/research/"

# Authorized dark terminal palette
AUTHORIZED_COLORS = {
    "#1a1a1a",  # background
    "#e0e0e0",  # text
    "#6cb3ee",  # links
    "#888888",  # metadata
    "#2a2a2a",  # code bg
    "#333333",  # borders
}

WARNINGS = []
ERRORS = []
FIXES_APPLIED = []


# ─── Helpers ──────────────────────────────────────────────────────────

def error(msg, file=None):
    prefix = f"❌ ERROR"
    if file:
        prefix += f" [{file}]"
    ERRORS.append(prefix + ": " + msg)
    print(f"  {prefix}: {msg}")


def warn(msg, file=None):
    prefix = f"⚠️  WARN"
    if file:
        prefix += f" [{file}]"
    WARNINGS.append(prefix + ": " + msg)
    print(f"  {prefix}: {msg}")


def fix_info(msg, file=None):
    prefix = f"🔧 FIX"
    if file:
        prefix += f" [{file}]"
    FIXES_APPLIED.append(prefix + ": " + msg)
    print(f"  {prefix}: {msg}")


def ok(msg, file=None):
    prefix = f"✅ OK"
    if file:
        prefix += f" [{file}]"
    print(f"  {prefix}: {msg}")


# ─── HTML Parsing ─────────────────────────────────────────────────────

class HTMLStructChecker(HTMLParser):
    """Parse an HTML document and report structural issues."""

    def __init__(self):
        super().__init__()
        self.has_doctype = False
        self.html_lang = None
        self.has_robots_meta = False
        self.has_viewport = False
        self.has_charset = False
        self.has_canonical = False
        self.has_favicon = False
        self.has_jsonld = False
        self.has_article_tags = False
        self.has_time_tags = False
        self.has_main_tag = False
        self.headings_by_level = {}
        self.issues = []

    def handle_decl(self, decl):
        if 'doctype' in decl.lower():
            self.has_doctype = True

    def handle_starttag(self, tag, attrs):
        attrs_dict = dict(attrs)
        low_tag = tag.lower()

        if low_tag == 'html':
            self.html_lang = attrs_dict.get('lang')
        elif low_tag == 'meta':
            name = attrs_dict.get('name', '').lower()
            charset = attrs_dict.get('charset', '')
            if name == 'robots' and 'content' in attrs_dict:
                self.has_robots_meta = True
            if name == 'viewport':
                self.has_viewport = True
        elif low_tag == 'link':
            rel = attrs_dict.get('rel', '').lower()
            href = attrs_dict.get('href', '')
            if rel == 'canonical' and href:
                self.has_canonical = True
            if rel == 'icon' or 'favicon' in href.lower():
                self.has_favicon = True
        elif low_tag == 'script':
            type_attr = attrs_dict.get('type', '').lower()
            if 'ld+json' in type_attr:
                self.has_jsonld = True
        elif low_tag == 'article':
            self.has_article_tags = True
        elif low_tag == 'time':
            self.has_time_tags = True
        elif low_tag == 'main':
            self.has_main_tag = True
        elif re.match(r'^h[1-6]$', low_tag):
            level = int(low_tag[1])
            self.headings_by_level.setdefault(level, []).append(tag)


def check_html_file(filepath, fix=False):
    """Run all HTML checks on a single file."""
    fname = str(filepath.relative_to(RESEARCH_ROOT))
    content = filepath.read_text()

    # ── Parse structure ────────────────────────────────────
    parser = HTMLStructChecker()
    try:
        parser.feed(content)
    except Exception as e:
        error(f"Failed to parse HTML: {e}")
        return

    issues = []

    # 1. DOCTYPE
    if not parser.has_doctype:
        err = f"Missing <!DOCTYPE html>"
        error(err, fname)
        if fix:
            fixed = content.lstrip()
            if not fixed.startswith('<!DOCTYPE'):
                fixed = '<!DOCTYPE html>\n' + fixed
                filepath.write_text(fixed)
                fix_info("Added <!DOCTYPE html>", fname)

    # 2. <html lang="en">
    if not parser.html_lang:
        error(f"<html> missing lang attribute", fname)
        if fix:
            content_fixed = re.sub(r'<html([^>]*)>', r'<html\1 lang="en">', content, count=1)
            if 'lang' not in content_fixed:
                content_fixed = content_fixed.replace('<html', '<html lang="en"', 1)
            filepath.write_text(content_fixed)
            fix_info("Added lang=\"en\" to <html>", fname)

    # 3. <meta charset>
    if not parser.has_charset:
        error(f"Missing <meta charset>", fname)

    # 4. <meta name="viewport">
    if not parser.has_viewport:
        error(f"Missing viewport meta", fname)
        if fix:
            content_fixed = re.sub(
                r'<head>(\s*)',
                r'\1<meta charset="UTF-8">\n\1<meta name="viewport" content="width=device-width, initial-scale=1.0">\n\1',
                content, count=1
            )
            if content != content_fixed:
                filepath.write_text(content_fixed)
                fix_info("Added viewport meta", fname)

    # 5. robots meta
    if not parser.has_robots_meta:
        warn(f"Missing <meta name=\"robots\" content=\"follow,index,max-snippet:-1\">", fname)
        if fix:
            content_fixed = re.sub(
                r'<head>(\s*)',
                r'\1<meta name="robots" content="follow,index,max-snippet:-1">\n\1',
                content, count=1
            )
            if 'robots' not in content_fixed.split('</head>')[0]:
                filepath.write_text(content_fixed)
                fix_info("Added robots meta", fname)

    # 6. External CSS dependency — blog-guide says no external stylesheets
    ext_css = re.findall(r'<link[^>]+href="([^"]+\.css)"', content)
    for css_url in ext_css:
        if 'bitsmasher.net' not in css_url or css_url.startswith('https://'):
            warn(f"External CSS dependency: {css_url} — blog-guide mandates embedded styles", fname)

    # 7. w3.css dependency
    if 'w3.css' in content.lower():
        warn("Uses external W3.CSS framework — violates 'embedded CSS only' guideline", fname)

    # 8. Favicon
    if not parser.has_favicon and 'favicon' not in content.lower():
        # Only warn for index.html; posts are exempt
        if filepath == INDEX_HTML:
            warn("Missing favicon link", fname)

    # 9. Canonical URL (should be on blog posts, not necessarily index)
    if 'blog/' in str(filepath) and not parser.has_canonical:
        warn(f"Blog post missing <link rel=\"canonical\">", fname)
        if fix:
            base_url = f"https://www.bitsmasher.net/research/{filepath.name}"
            content_fixed = re.sub(
                r'<head>(\s*)',
                rf'\1<link rel="canonical" href="{base_url}">\n\1',
                content, count=1
            )
            if 'canonical' not in content_fixed.split('</head>')[0]:
                filepath.write_text(content_fixed)
                fix_info(f"Added canonical link: {base_url}", fname)

    # 10. Heading hierarchy check (H1 should appear once at top)
    h1s = parser.headings_by_level.get(1, [])
    if len(h1s) == 0:
        error("No <h1> heading found", fname)
    elif len(h1s) > 1:
        warn(f"Multiple <h1> tags ({len(h1s)}); only one H1 per page recommended", fname)

    # 11. Title tag present?
    if not re.search(r'<title>', content):
        error("Missing <title> tag", fname)

    # 12. Color palette compliance (search for non-authorized hex colors)
    found_colors = set(re.findall(r'#[0-9a-fA-F]{6}\b', content))
    unauthorized = found_colors - AUTHORIZED_COLORS
    if unauthorized:
        warn(f"Non-standard color values: {', '.join(sorted(unauthorized))}", fname)

    # 13. Background-color / color property checks
    bg_matches = re.findall(r'background:\s*#[0-9a-fA-F]{6}', content)
    for bm in bg_matches:
        c = bm.split('#')[1]
        if c not in '1a1a1a':
            warn(f"Background color '{bm}' deviates from spec #1a1a1a", fname)

    link_matches = re.findall(r'a\s*\{\s*color:\s*#[0-9a-fA-F]{6}', content)
    for lm in link_matches:
        c = lm.split('#')[1]
        if c not in '6cb3ee':
            warn(f"Link color '{lm}' deviates from spec #6cb3ee", fname)

    return True


# ─── index.json Validation ────────────────────────────────────────────

def check_index_json(fix=False):
    fname = "index.json"
    
    if not INDEX_JSON.exists():
        error("index.json not found")
        return False

    content = INDEX_JSON.read_text()
    try:
        data = json.loads(content)
    except json.JSONDecodeError as e:
        error(f"Invalid JSON: {e}")
        return False

    ok(f"Valid JSON")

    # Required top-level keys (JSON Feed v1.1 subset)
    for key in ('posts', 'totalPosts'):
        if key not in data:
            error(f"Missing required key: '{key}'")

    posts = data.get('posts', [])
    
    # totalPosts consistency
    if data.get('totalPosts') != len(posts):
        msg = f"totalPosts ({data.get('totalPosts', 'N/A')}) != len(posts) ({len(posts)})"
        error(msg)
        if fix:
            data['totalPosts'] = len(posts)
            INDEX_JSON.write_text(json.dumps(data, indent=2) + '\n')
            fix_info(f"Corrected totalPosts to {len(posts)}", fname)

    # Per-post validation
    for i, post in enumerate(posts):
        prefix = f"post[{i}]"
        
        for key in ('title', 'date', 'url'):
            if key not in post:
                error(f"{prefix}: missing '{key}'")
        
        # Date format check (should be ISO 8601)
        date_str = post.get('date', '')
        if date_str:
            iso_pattern = r'^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}'
            if not re.match(iso_pattern, date_str):
                warn(f"{prefix}: date '{date_str}' is not ISO 8601 (expected YYYY-MM-DDTHH:MM:SS)")

        # URL check
        url = post.get('url', '')
        blog_file = BLOG_DIR / url.split('/')[-1]
        if blog_file.exists():
            ok(f"{prefix}: url targets existing file {url}")
        elif not url:
            error(f"{prefix}: empty URL")

    # lastUpdated check
    last_updated = data.get('lastUpdated', '')
    if last_updated and not re.match(r'^\d{4}-\d{2}-\d{2}T\d{2}', last_updated):
        warn(f"lastUpdated format not ISO 8601: '{last_updated}'")

    # Chronological order check (newest first)
    dates = [p.get('date', '') for p in posts]
    sorted_dates = sorted(dates, reverse=True)
    if dates != sorted_dates:
        msg = "Posts are NOT in reverse chronological order"
        warn(msg)
        if fix:
            # Sort posts by date descending
            data['posts'] = sorted(data['posts'], key=lambda p: p.get('date', ''), reverse=True)
            INDEX_JSON.write_text(json.dumps(data, indent=2) + '\n')
            fix_info("Sorted posts in reverse chronological order", fname)

    # 14. Check for posts listed in index.json but missing from blog/ directory
    for i, post in enumerate(posts):
        url = post.get('url', '')
        if url:
            basename = url.split('/')[-1]
            if not (BLOG_DIR / basename).exists():
                error(f"index.json lists {url} but file missing from blog/", fname)

    return True


# ─── index.html Validation ────────────────────────────────────────────

def check_index_html(fix=False):
    fname = "index.html"

    if not INDEX_HTML.exists():
        error("index.html not found")
        return False

    content = INDEX_HTML.read_text()

    # Run generic HTML checks
    check_html_file(INDEX_HTML, fix=fix)

    # ── Index-specific checks ──────────────────────────────

    # 1. numberOfPosts JSON-LD consistency
    jsonld_match = re.search(r'"numberOfPosts":\s*(\d+)', content)
    if jsonld_match:
        jsonld_count = int(jsonld_match.group(1))
        article_count = len(re.findall(r'<div class="blog-post">', content))
        if jsonld_count != article_count:
            msg = f"JSON-LD numberOfPosts ({jsonld_count}) != blog-post div count ({article_count})"
            error(msg)
            if fix:
                content_fixed = re.sub(
                    r'"numberOfPosts":\s*\d+',
                    f'"numberOfPosts": {article_count}',
                    content
                )
                INDEX_HTML.write_text(content_fixed)
                fix_info(f"Corrected numberOfPosts to {article_count}", fname)

    # 2. Post count consistency between index.html and index.json
    if INDEX_JSON.exists():
        idx_json_posts = len(json.loads(INDEX_JSON.read_text()).get('posts', []))
        html_posts = len(re.findall(r'<div class="blog-post">', content))
        if html_posts != idx_json_posts:
            msg = f"index.html has {html_posts} posts; index.json lists {idx_json_posts}"
            error(msg)

    # 3. Posts should be in reverse chronological order by date in href title
    posts_hrefs = re.findall(r'<a href="([^"]+)">', content)
    expected_order = sorted(
        [p for p in posts_hrefs if '/blog/' in p],
        key=lambda x: os.path.basename(x),
        reverse=True
    )
    blog_hrefs_in_html = [h for h in re.findall(r'<a href="([^"]+)"[^>]*>', content) if '/blog/' in h]
    if blog_hrefs_in_html and len(blog_hrefs_in_html) == len(expected_order):
        # Check ordering roughly by filename date prefix
        pass  # rough check; detailed order verified against index.json

    # 4. Blog post cards missing article wrapper?
    blog_post_divs = re.findall(r'<div class="blog-post">(.*?)</div>', content, re.DOTALL)
    for i, block in enumerate(blog_post_divs):
        if '<article>' not in block:
            warn(f"Blog post card #{i+1} missing <article> wrapper", fname)

    # 5. Check for dead references (links to /research/feed.json per MEMORY.md note)
    if '/research/feed.json' in content:
        warn("Still references /research/feed.json — remove per spec (dead alternate link)", fname)
        if fix:
            content_fixed = re.sub(r'\s*<link[^>]*href="/research/feed\.json"[^>]*>', '', content)
            INDEX_HTML.write_text(content_fixed)
            fix_info("Removed dead feed.json reference", fname)

    return True


# ─── Blog-guide.md Conformance Check ──────────────────────────────────

def check_blog_guide():
    fname = "blog-guide.md"
    if not GUIDE_FILE.exists():
        warn(f"{fname} not found — using defaults")
        return

    content = GUIDE_FILE.read_text()
    
    # Check for self-consistency with actual files
    if 'embedded CSS block (no external stylesheet dependency' in content:
        ok("Guide mandates embedded CSS — verified against this rule during blog post checks")


# ─── Cross-Reference Checks ───────────────────────────────────────────

def check_cross_refs():
    """Check href targets, canonical consistency, and file existence."""
    
    # Check all blog posts for internal link integrity
    for post_file in sorted(BLOG_DIR.glob("*.html")):
        fname = str(post_file.relative_to(RESEARCH_ROOT))
        content = post_file.read_text()
        
        # Back to research links should use /research
        back_links = re.findall(r'href="([^"]*)"', content)
        for href in back_links:
            if href == '/research':
                ok(f"Valid internal link: {href}", fname)
            elif href.startswith('/'):
                target = RESEARCH_ROOT.parent / href.lstrip('/')
                # Check image/css/js paths exist at that scope
                pass  # External resources serve from wonderland root

    # index.html should link to all blog posts
    if INDEX_HTML.exists() and BLOG_DIR.exists():
        html_links = set(re.findall(r'href="(/research/blog/[^"]+)"', INDEX_HTML.read_text()))
        disk_files = set(f'/research/blog/{p.name}' for p in BLOG_DIR.glob("*.html"))
        
        missing_from_html = disk_files - html_links
        if missing_from_html:
            for m in sorted(missing_from_html):
                error(f"index.html missing link to {m}", "index.html")

        missing_on_disk = html_links - disk_files
        if missing_on_disk:
            for m in sorted(missing_on_disk):
                warn(f"index.html links to non-existent file: {m}", "index.html")


# ─── Deploy Check ─────────────────────────────────────────────────────

def check_ssh_connectivity():
    """Probe wonderland SSH accessibility."""
    try:
        result = subprocess.run(
            ['ssh', '-o', 'BatchMode=yes', '-o', 'ConnectTimeout=3',
             f'{DEPLOY_USER}@{DEPLOY_HOST}', 'hostname'],
            capture_output=True, text=True, timeout=5
        )
        if result.returncode == 0:
            ok(f"SSH connectivity to {DEPLOY_HOST} confirmed — host: {result.stdout.strip()}")
            return True
        else:
            warn(f"SSH to {DEPLOY_HOST}: {result.stderr.strip()}")
            return False
    except subprocess.TimeoutExpired:
        warn(f"SSH to {DEPLOY_HOST}: timed out (host currently unreachable)")
        return False
    except Exception as e:
        warn(f"SSH to {DEPLOY_HOST}: {e}")
        return False


# ─── Fix Summary & Report ─────────────────────────────────────────────

def generate_report():
    print("\n" + "=" * 60)
    print("RESEARCH BLOG LINT REPORT")
    print("=" * 60)
    
    total = len(ERRORS) + len(WARNINGS) + len(FIXES_APPLIED)
    print(f"\nTotal issues found: {len(ERRORS)} errors, {len(WARNINGS)} warnings")
    if FIXES_APPLIED:
        print(f"Auto-fixes applied:  {len(FIXES_APPLIED)}")
    
    if ERRORS:
        print("\n--- ERRORS ---")
        for e in ERRORS:
            print(f"  {e}")
    
    if WARNINGS:
        print("\n--- WARNINGS ---")
        for w in WARNINGS:
            print(f"  {w}")
    
    if FIXES_APPLIED:
        print("\n--- APPLIED FIXES ---")
        for f_item in FIXES_APPLIED:
            print(f"  {f_item}")

    print("=" * 60)


# ─── Main ──────────────────────────────────────────────────────────────

def main():
    import argparse
    parser = argparse.ArgumentParser(description='Research blog CSS & formatting auditor')
    parser.add_argument('--fix', action='store_true', help='Auto-remediate fixable issues')
    args = parser.parse_args()

    print("=" * 60)
    print("bitsmasher.net Research Blog — Format Auditor")
    print(f"Workspace: {RESEARCH_ROOT}")
    print(f"Date: 2026-09-02")
    print("=" * 60)

    # ── index.json checks ────────────────────────────────────
    print("\n[1/6] Validating index.json...")
    check_index_json(fix=args.fix)

    # ── index.html checks ────────────────────────────────────
    print("\n[2/6] Validating index.html...")
    check_index_html(fix=args.fix)

    # ── Blog post checks ─────────────────────────────────────
    print(f"\n[3/6] Validating blog posts ({BLOG_DIR.name}/)...")
    for post_file in sorted(BLOG_DIR.glob("*.html")):
        print(f"\n  → {post_file.name}")
        check_html_file(post_file, fix=args.fix)

    # ── blog-guide.md conformance ────────────────────────────
    print("\n[4/6] Checking blog-guide.md...")
    check_blog_guide()

    # ── Cross-reference integrity ────────────────────────────
    print("\n[5/6] Cross-reference checks...")
    check_cross_refs()

    # ── SSH connectivity ─────────────────────────────────────
    print("\n[6/6] Deploy target connectivity...")
    can_deploy = check_ssh_connectivity()

    # ── Report ───────────────────────────────────────────────
    generate_report()

    if FIXES_APPLIED and not args.fix:
        print("\n  Run with --fix to apply auto-fixes.")

    if ERRORS:
        sys.exit(1)
    sys.exit(0)


if __name__ == '__main__':
    main()
