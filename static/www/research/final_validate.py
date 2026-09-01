import os, json
from pathlib import Path

BASE = Path("/home/franklin/workspace/website/static/www/research")
BLOG = BASE / "blog"
TAGS = BASE / "tags"
errors = []

# 1. Index vs blog
with open(BASE / "index.json") as f:
    idx = json.load(f)
idx_urls = set(p["url"] for p in idx["posts"])
disk_urls = set("https://www.bitsmasher.net/research/blog/" + f.stem + ".html" for f in BLOG.glob("*.html"))
if idx_urls != disk_urls:
    errors.append("index.json URL mismatch with disk: %s" % (idx_urls ^ disk_urls))

# 2. Feed vs index
with open(BASE / "feed.json") as f:
    feed = json.load(f)
feed_urls = set(i["url"] for i in feed["items"])
if idx_urls != feed_urls:
    errors.append("Feed URL mismatch with index: %s" % (idx_urls ^ feed_urls))

# 3. Tag format
for item in feed["items"]:
    for t in item.get("tags", []):
        if not (t.startswith("/") and t.endswith("/") and len(t) > 2):
            errors.append('Bad tag "%s" in %s' % (t, item["title"]))

# 4. Required tags
for tag in ["defense", "badge-security", "research"]:
    p = TAGS / (tag + ".json")
    if not p.exists():
        errors.append("Missing required tag: " + tag)

# 5. Nav/styling
html_files = []
for root, dirs, files in os.walk(BASE):
    for f in files:
        if f.endswith(".html"):
            html_files.append(Path(root) / f)

bad_nav = []
search_strings = [
    'id="site-nav"',
    "/js/nav.js",
    "/w3css/5/w3.css",
    "/css/www.bitsmasher.net.css",
]
for hp in html_files:
    try:
        c = hp.read_text(encoding="utf-8", errors="ignore")
    except Exception as e:
        bad_nav.append((str(hp), ["unreadable: " + str(e)]))
        continue
    found = {}
    for s in search_strings:
        found[s] = s in c
    missing = [s for s, v in found.items() if not v]
    if missing:
        bad_nav.append((str(hp), missing))

if bad_nav:
    errors.append("Nav/styling issues (%d files):" % len(bad_nav))
    for name, issues in bad_nav[:3]:
        errors.append("  " + name + ": " + str(issues))

# 6. Permissions
for root, dirs, files in os.walk(BASE):
    m = oct(Path(root).stat().st_mode)[-3:]
    if m != "775":
        errors.append("DIR perms %s: %s" % (root, m))
for hp in html_files:
    m = oct(hp.stat().st_mode)[-3:]
    if m != "664":
        errors.append("FILE perms %s: %s" % (hp, m))

# Print results
print("=" * 60)
print("FINAL VALIDATION RESULT")
print("=" * 60)
total_files = len(html_files)
total_posts = len(list(BLOG.glob("*.html")))
total_tags = len(list(TAGS.glob("*.json")))
total_idx = len(idx["posts"])
total_feed = len(feed["items"])

print("Files scanned: %d HTML" % total_files)
print("Blog posts on disk: %d" % total_posts)
print("Index.json entries: %d" % total_idx)
print("Feed.json items: %d" % total_feed)
print("Tag manifests: %d" % total_tags)

if errors:
    print("\nERRORS (%d):" % len(errors))
    for e in errors:
        print("  X " + e)
else:
    print("\nALL CHECKS PASSED")
print("=" * 60)
