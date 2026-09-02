#!/usr/bin/env python3
"""Update ISMG-2050 course page and portal with Central Sierra Insurance + card reordering + unicode cleanup."""

import json
import os
import sys

TEACHING_DIR = "/home/franklin/workspace/website/static/www/teaching"
ISMG_PATH = f"{TEACHING_DIR}/ISMG-2050/index.html"
PORTAL_PATH = f"{TEACHING_DIR}/ISMG-2050/homework/week3-4/index.html"

def read_file(path):
    with open(path, "r") as f:
        return f.read()

def write_file(path, content):
    with open(path, "w", encoding="utf-8") as f:
        f.write(content)

# ---- Unicode cleanup ----
def fix_unicode_in_html(content):
    """Replace literal \u2013 and \u20134 sequences with proper UTF-8/en-dash."""
    # Replace literal escaped sequences that got written literally
    content = content.replace('\\u20134', '\u2013')  # \u20134 -> \u2013 (en-dash + "4" should just be en-dash)
    content = content.replace('\\u2013', '\u2013')   # literal \u2013 -> actual en-dash char
    content = content.replace('\\u0026', '&')         # literal \u0026 -> &
    return content

def fix_unicode_in_json(content):
    """In JSON, replace literal escaped unicode with proper UTF-8 strings."""
    # json.loads will already decode standard escapes; we just need to make sure
    # there are no double-encoded literals
    content = content.replace('\\u2013', '\u2013')
    return content

# ---- Card Reordering for ISMG-2050 index.html ----
def reorder_card_revenue(content):
    """Reorganize the main hub page so essential links appear near top."""

    # Strategy: rebuild the body section with proper order.
    # Find key anchor points in the HTML to reconstruct.

    lines = content.split('\n')

    # Find insertion point right after <header class="page-header">...  </div> block
    # and before the existing info-block blocks.

    # Let's find where the first content div starts (after header)
    # and insert our reordered sections there.

    result = []
    in_header = False
    header_done = False
    skipped_old_links = False
    old_syllabus_done = False
    old_history_done = False

    i = 0
    while i < len(lines):
        line = lines[i]

        # After the course-desc div, insert quick-access links FIRST
        if not header_done and '<div class="link-section">' in line:
            # This is the old link section (syllabus + history). We'll move these up.
            # Skip this and collect what's inside to insert earlier
            # Actually, let's just restructure differently - collect sections

            # Insert quick-access links before this point
            quick_links = '''				<div class="info-block" style="background:#1a1a1a;border-left:4px solid #CFB87C;">
					<p class="info-label">Quick Access</p>
					<p class="info-value" style="margin-top:0.5rem;">
						<a href="syllabus/syllabus_fall_2026.pdf" style="color:#4ade80;text-decoration:none;margin-right:1rem;">\U0001f4d4 Syllabus (PDF)</a>
						<a href="history.html" style="color:#4ade80;text-decoration:none;margin-right:1rem;">\U0001f4c3 Course History &amp; Legacy</a>
					</p>
				</div>

'''
            result.append(quick_links)
            header_done = True
            i += 1
            continue

        # Skip the old <div class="link-section"> block entirely (we moved it)
        if '<div class="link-section">' in line:
            # skip to </div>
            while i < len(lines) and '</div>' not in lines[i]:
                i += 1
            i += 1  # skip the closing </div>
            old_syllabus_done = True
            continue

        result.append(line)
        i += 1

    return '\n'.join(result)


# ---- Portal index update: add Central Sierra to Module 2 section ----
def update_portal_index(content):
    """Insert Central Sierra Insurance into the portal listing."""

    # Find the line with Practice 4 heading and insert after the Week 3 label
    lines = content.split('\n')
    result = []
    inserted_cs = False

    for i, line in enumerate(lines):
        result.append(line)

        # After "Week 3 (Sep 7\u201313) \u2014; Flash Fill &amp; Math Operators" heading
        if 'Flash Fill' in line and 'Week 3' in line:
            cs_card = '''
<div class="assignment-row" style="border-left-color:#CFB87C">
<p class="assignment-type">Comprehensive Project</p>
<p class="assignment-name"><a href="centralsierra.html" style="color:#CFB87C;text-decoration:none">Assignment 2.3.1: Central Sierra Insurance (XLOOKUP, SUMIF, IF, CONCAT)</a></p>
<p class="assignment-desc">Defined names from selection, XLOOKUP approximate match (match_mode=-1), order-of-operations commission formula (=E5*(1+F5)), SUMIF branch aggregation, statistical distribution (AVERAGE/MODE/MEDIAN), IF policy classification, and CONCAT label assembly. Builds Commission_Level and Bonus_Rate named ranges from Tables sheet.</p>
<p class="assignment-meta">Due: Sunday, September 13 at 11:59 PM MT \u2022 Starter: <code>CentralSierra-02.xlsx</code> \u2022 Submit: <code>[initials] CentralSierraIns.xlsx</code></p>
</div>

'''
            result.append(cs_card)
            inserted_cs = True

    if not inserted_cs:
        print("WARNING: Could not insert Central Sierra card into portal", file=sys.stderr)

    return '\n'.join(result)


# ---- Main ----
if __name__ == '__main__':
    # 1. Update ISMG-2050 main hub
    print("Processing ISMG-2050 index.html...")
    content = read_file(ISMG_PATH)
    content = fix_unicode_in_html(content)
    content = reorder_card_revenue(content)
    write_file(ISMG_PATH, content)
    print(f"  Written to {ISMG_PATH}")

    # 2. Update portal index
    print("Processing week3-4/index.html...")
    content = read_file(PORTAL_PATH)
    content = fix_unicode_in_html(content)
    content = update_portal_index(content)
    write_file(PORTAL_PATH, content)
    print(f"  Written to {PORTAL_PATH}")

    # 3. Update teaching/index.json if it exists
    json_path = f"{TEACHING_DIR}/index.json"
    if os.path.exists(json_path):
        print("Processing teaching/index.json...")
        with open(json_path, "r") as f:
            content = f.read()
        content = fix_unicode_in_json(content)
        with open(json_path, "w", encoding="utf-8") as f:
            f.write(content)
        print(f"  Written to {json_path}")

    # 4. Update teaching/index.html if it exists
    index_html_path = f"{TEACHING_DIR}/index.html"
    if os.path.exists(index_html_path):
        print("Processing teaching/index.html...")
        content = read_file(index_html_path)
        content = fix_unicode_in_html(content)
        write_file(index_html_path, content)
        print(f"  Written to {index_html_path}")

    print("\nDone.")
