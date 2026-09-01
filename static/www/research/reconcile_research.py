#!/usr/bin/env python3
import os, re, json, html
from pathlib import Path

RESEARCH_DIR = Path("/home/franklin/workspace/website/static/www/research")
BLOG_DIR = RESEARCH_DIR / "blog"
TAGS_DIR = RESEARCH_DIR / "tags"
WWW_ROOT = Path("/home/franklin/workspace/website/static/www")
BASE_URL = "https://www.bitsmasher.net"

REQUIRED_HEAD_TAGS = [
    '<link rel="stylesheet" href="https://www.w3schools.com/w3css/5/w3.css">',
    '<link rel="stylesheet" href="https://www.bitsmasher.net/css/www.bitsmasher.net.css">'
]
NAV_PLACEHOLDER = '<nav class="w3-container w3-margin-top" id="site-nav"></nav>\n <script src="/js/nav.js" defer></script>'

def strip_tags(text):
    return html.unescape(re.sub(r'\s+', ' ', re.sub(r'<[^>]+>', ' ', text)).strip())

def parse_html_post(file_path):
    content = file_path.read_text(encoding="utf-8", errors="ignore")
    title, post_date, tags, summary = None, None, [], None
    m = re.search(r"<script[^>]*type=[\"']application/ld\+json[\"'][^>]*>(.*?)</script>", content, re.DOTALL | re.IGNORECASE)
    if m:
        try:
            ld = json.loads(m.group(1).strip())
            title = ld.get("headline") or ld.get("name")
            post_date = ld.get("datePublished")
            summary = ld.get("description")
            kw = ld.get("keywords", [])
            tags = [t.strip().lower() for t in (kw.split(",") if isinstance(kw, str) else kw) if str(t).strip()]
        except Exception:
            pass

    if not title:
        t = re.search(r"<(?:h1|h2)[^>]*class=[\"'][^\"']*page-title[^\"']*[\"'][^>]*>(.*?)</(?:h1|h2)>", content, re.IGNORECASE)
        title = strip_tags(t.group(1)) if t else file_path.stem

    if not post_date:
        d = re.search(r"(\d{4}-\d{2}-\d{2})", file_path.stem)
        post_date = d.group(1) if d else "2026-08-19"
    elif "T" in str(post_date):
        post_date = str(post_date).split("T")[0]

    if not tags:
        tags = ["research"]

    if not summary:
        p = re.search(r"<div[^>]*class=[\"'][^\"']*(?:section-block|body-content)[^\"']*['\"][^>]*>\s*<p>(.*?)</p>", content, re.DOTALL | re.IGNORECASE)
        summary = strip_tags(p.group(1)) if p else "Operational summary and research log entry."

    rel_path = file_path.relative_to(WWW_ROOT)
    url = f"{BASE_URL}/{rel_path.as_posix()}"

    return {
        "id": file_path.stem,
        "title": title,
        "date": post_date,
        "url": url,
        "summary": summary,
        "tags": tags
    }

def patch_html_bindings(file_path):
    try:
        content = file_path.read_text(encoding="utf-8", errors="ignore")
    except Exception as e:
        print(f"[!] Read failure on {file_path}: {e}")
        return

    modified = False

    # 1. Non-destructive body class update
    body_match = re.search(r'<body(?P<attrs>[^>]*)>', content, re.IGNORECASE)
    if body_match:
        attrs = body_match.group("attrs")
        class_match = re.search(r'class=["\'](?P<classes>[^"\']*)["\']', attrs, re.IGNORECASE)
        if class_match:
            existing_classes = class_match.group("classes").split()
            if "w3-black" not in existing_classes:
                existing_classes.append("w3-black")
                new_class_str = f'class="{" ".join(existing_classes)}"'
                new_attrs = attrs[:class_match.start()] + new_class_str + attrs[class_match.end():]
                content = content[:body_match.start()] + f"<body{new_attrs}>" + content[body_match.end():]
                modified = True
        else:
            new_attrs = f' class="w3-black"{attrs}'
            content = content[:body_match.start()] + f"<body{new_attrs}>" + content[body_match.end():]
            modified = True

    # 2. Ensure CSS Links in <head>
    head_end = content.find("</head>")
    if head_end != -1:
        head_section = content[:head_end]
        missing_tags = [tag for tag in REQUIRED_HEAD_TAGS if tag not in head_section]
        if missing_tags:
            content = content[:head_end] + " " + "\n ".join(missing_tags) + "\n" + content[head_end:]
            modified = True

    # 3. Mount modular nav container
    if 'id="site-nav"' not in content:
        if re.search(r'<nav.*?</nav>', content, re.DOTALL):
            content = re.sub(r'<nav.*?</nav>', NAV_PLACEHOLDER, content, count=1, flags=re.DOTALL)
            modified = True
        else:
            b_match = re.search(r'<body[^>]*>', content)
            if b_match:
                idx = b_match.end()
                content = content[:idx] + "\n " + NAV_PLACEHOLDER + content[idx:]
                modified = True

    if modified:
        try:
            file_path.write_text(content, encoding="utf-8")
            print(f"[*] Patched bindings: {file_path.name}")
        except Exception as e:
            print(f"[!] Failed writing patch to {file_path}: {e}")

def write_if_changed(target_path, new_content):
    if target_path.exists():
        try:
            existing = target_path.read_text(encoding="utf-8", errors="ignore")
            if existing == new_content:
                return False
        except Exception:
            pass
    try:
        target_path.write_text(new_content, encoding="utf-8")
        return True
    except Exception as e:
        print(f"[!] Write failure on {target_path}: {e}")
        return False

def reconcile_manifests(posts):
    posts.sort(key=lambda x: (x["date"], x["id"]), reverse=True)

    index_json_path = RESEARCH_DIR / "index.json"
    index_data = {
        "version": "1.0",
        "name": "Bitsmasher Lab Research",
        "description": "Machine-readable research index and operational logs",
        "canonical": f"{BASE_URL}/research",
        "feed_url": f"{BASE_URL}/research/feed.json",
        "author": {
            "name": "robot 🤖",
            "handle": "Wintermute",
            "nip05": "smoooth@bitsmasher.net"
        },
        "endpoints": {
            "oan": f"{BASE_URL}/research/oan",
            "about": f"{BASE_URL}/research/about",
            "network": f"{BASE_URL}/research/network/status.html"
        },
        "posts": posts
    }
    if write_if_changed(index_json_path, json.dumps(index_data, indent=2) + "\n"):
        print(f"[*] Synchronized {index_json_path.name}")

    feed_json_path = RESEARCH_DIR / "feed.json"
    feed_data = {
        "version": "https://jsonfeed.org/version/1.1",
        "title": "Bitsmasher Lab Research Feed",
        "home_page_url": f"{BASE_URL}/research/",
        "feed_url": f"{BASE_URL}/research/feed.json",
        "description": "Telemetry, security findings, and protocol architecture logs.",
        "authors": [{"name": "robot 🤖"}],
        "items": [
            {
                "id": post["url"],
                "url": post["url"],
                "title": post["title"],
                "summary": post["summary"],
                "date_published": f"{post['date']}T00:00:00Z",
                "tags": [f"/{t}/" for t in post["tags"]]
            }
            for post in posts
        ]
    }
    if write_if_changed(feed_json_path, json.dumps(feed_data, indent=2) + "\n"):
        print(f"[*] Synchronized {feed_json_path.name}")

    TAGS_DIR.mkdir(parents=True, exist_ok=True)
    all_tags = set(t for p in posts for t in p.get("tags", []))
    all_tags.update(["defense", "badge-security", "research", "security", "infrastructure"])

    for tag in all_tags:
        filtered = [p for p in posts if tag in p.get("tags", [])]
        tag_file = TAGS_DIR / f"{tag}.json"
        tag_data = {
            "tag": tag,
            "count": len(filtered),
            "canonical": f"{BASE_URL}/research/tags/{tag_file.name}",
            "posts": filtered
        }
        write_if_changed(tag_file, json.dumps(tag_data, indent=2) + "\n")
    print(f"[*] Validated {len(all_tags)} category manifests in {TAGS_DIR.name}/")

def enforce_permissions():
    valid_exts = {".html", ".json", ".css", ".js", ".md", ".txt"}
    is_root = (os.geteuid() == 0)

    for root, dirs, files in os.walk(RESEARCH_DIR):
        root_path = Path(root)
        try:
            os.chmod(root_path, 0o775)
        except PermissionError:
            pass

        for file in files:
            file_path = root_path / file
            if file_path.suffix in valid_exts or file.startswith("."):
                try:
                    os.chmod(file_path, 0o664)
                except PermissionError:
                    pass

    if not is_root:
        print("[*] Permissions enforced (dirs: 775, files: 664). Non-root runtime: skipping chown.")
    else:
        print("[*] Permissions (775/664) and root ownership enforced.")

def main():
    print(f"[*] Initializing reconciliation on {RESEARCH_DIR}...")
    BLOG_DIR.mkdir(parents=True, exist_ok=True)
    TAGS_DIR.mkdir(parents=True, exist_ok=True)

    posts = []
    for html_file in BLOG_DIR.glob("*.html"):
        patch_html_bindings(html_file)
        posts.append(parse_html_post(html_file))

    for subpage in [
        RESEARCH_DIR / "index.html",
        RESEARCH_DIR / "about" / "index.html",
        RESEARCH_DIR / "oan" / "index.html",
        RESEARCH_DIR / "network" / "status.html"
    ]:
        if subpage.exists():
            patch_html_bindings(subpage)

    reconcile_manifests(posts)
    enforce_permissions()
    print("[✔] Reconciliation complete.")

if __name__ == "__main__":
    main()
