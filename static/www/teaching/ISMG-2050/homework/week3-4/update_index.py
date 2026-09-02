import sys

path = "/home/franklin/workspace/website/static/www/teaching/ISMG-2050/index.html"
with open(path, "r") as f:
    content = f.read()

new_section = '''				<div class="info-block" style="background:#1a1a1a; border-left:4px solid #4ade80;">
					<h3 style="margin-top:0; color:#4ade80;">Module 2 Practice &amp; Homework Deliverables</h3>
					<p>Flash Fill, math operators, conditional formatting, data import, and page setup. Six practice exercises plus two comprehensive projects.</p>

					<h4 style="margin-bottom:0.3rem; color:#fff;">Practice 4: Flash Fill and Cell Referencing</h4>
					<ul style="margin-top:0.3rem;">
						<li>\U0001f4c4 <a href="/teaching/ISMG-2050/homework/week3-4/practice4.html"><strong>Directions &amp; Guide</strong></a></li>
						<li>\U0001f4ca <a href="/teaching/ISMG-2050/homework/week3-4/starter_data_wk3.xlsx" download><strong>Download Starter Workbook (.xlsx)</strong></a></li>
					</ul>

					<h4 style="margin-bottom:0.3rem; color:#fff;">Practice 5: Math Operators and Conditional Formatting</h4>
					<ul style="margin-top:0.3rem;">
						<li>\U0001f4c4 <a href="/teaching/ISMG-2050/homework/week3-4/practice5.html"><strong>Directions &amp; Guide</strong></a></li>
					</ul>

					<h4 style="margin-bottom:0.3rem; color:#fff;">Practice 6: Customer Tracking with Data Import</h4>
					<ul style="margin-top:0.3rem;">
						<li>\U0001f4c4 <a href="/teaching/ISMG-2050/homework/week3-4/practice6.html"><strong>Directions &amp; Guide</strong></a></li>
						<li>\U0001f4ca <a href="/teaching/ISMG-2050/homework/week3-4/customer_data.csv" download><strong>Download Customer Data (.csv)</strong></a></li>
					</ul>

					<h4 style="margin-bottom:0.3rem; color:#fff;">Practice 7: Page Setup and Multi-Sheet Print Optimization</h4>
					<ul style="margin-top:0.3rem;">
						<li>\U0001f4c4 <a href="/teaching/ISMG-2050/homework/week3-4/practice7.html"><strong>Directions &amp; Guide</strong></a></li>
					</ul>

					<h4 style="margin-bottom:0.3rem; color:#fff;">Comprehensive Project: College Cost Calculator</h4>
					<ul style="margin-top:0.3rem;">
						<li>\U0001f4c4 <a href="/teaching/ISMG-2050/homework/week3-4/college-cost.html"><strong>Project Directions &amp; Guide</strong></a></li>
						<li>\U0001f4cc Build from blank workbook (submit as <code>[initials]CollegeCost.xlsx</code>)</li>
					</ul>

					<h4 style="margin-bottom:0.3rem; color:#fff;">Comprehensive Project: Cost Analysis Report</h4>
					<ul style="margin-top:0.3rem;">
						<li>\U0001f4c4 <a href="/teaching/ISMG-2050/homework/week3-4/cost-analysis.html"><strong>Project Directions &amp; Guide</strong></a></li>
						<li>\U0001f4ca Starter: <a href="/teaching/ISMG-2050/homework/week3-4/starter_data_wk3.xlsx" download><strong>starter_data_wk3.xlsx</strong></a></li>
						<li>\U0001f4cc Submit as <code>[initials]CostAnalysis.xlsx</code></li>
					</ul>

					<p style="margin-top:1rem; margin-bottom:0;">
						<a href="/teaching/ISMG-2050/homework/week3-4/"><strong>Open Full Week 3\u20134 Homework Portal \u2192</strong></a>
					</p>
				</div>

'''

insert_point = content.find('\t\t\t\t<div class="link-section">')
if insert_point == -1:
    print("ERROR: Could not find insertion point")
    sys.exit(1)
else:
    new_content = content[:insert_point] + new_section + content[insert_point:]
    with open(path, "w") as f:
        f.write(new_content)
    print("OK: Module 2 section inserted into index.html")
