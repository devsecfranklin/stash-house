#!/usr/bin/env bash
# ==============================================================================
# LaTeX Lint Tool Deployment Script for Linux CLI
# ==============================================================================
# This script installs the 'latextools' Python package and creates a 
# system-accessible CLI wrapper ('latexlint') to parse and check LaTeX files.
# ==============================================================================

set -e

# Define color codes for UI feedback
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}[*] Starting LaTeX Lint tool installation...${NC}"

# 1. Verify python3 and pip are available
if ! command -v python3 &> /dev/null; then
    echo -e "${RED}[!] Error: python3 is required but not installed.${NC}" >&2
    exit 1
fi

if ! command -v pip3 &> /dev/null; then
    echo -e "${RED}[!] Error: pip3 is required but not installed.${NC}" >&2
    exit 1
fi

# 2. Install the latextools package safely via pip
echo -e "${GREEN}[*] Installing 'latextools' from PyPI...${NC}"
python3 -m pip install --upgrade latextools

# 3. Create a dedicated local binary wrapper directory if it doesn't exist
BIN_DIR="$HOME/.local/bin"
mkdir -p "$BIN_DIR"

# 4. Generate the 'latexlint' CLI command file
CLI_PATH="$BIN_DIR/latexlint"
echo -e "${GREEN}[*] Generating CLI wrapper at $CLI_PATH...${NC}"

cat << 'EOF' > "$CLI_PATH"
#!/usr/bin/env python3
import sys
import os
import latextools

def lint_latex_file(filepath):
    if not os.path.exists(filepath):
        print(f"Error: File '{filepath}' not found.", file=sys.stderr)
        sys.exit(1)
        
    print(f"Linting: {filepath}\n" + "-"*40)
    
    try:
        # Parse the LaTeX file using latextools
        file_buffer = latextools.FileBuffer.from_file(filepath)
        parser = latextools.LaTeXParser(file_buffer)
        document = parser.parse()
        
        # Check for structural components or empty anomalies
        content_nodes = document.get_content()
        
        if not content_nodes:
            print("Warning: Document structural content appears to be empty.")
        else:
            print(f"Success: Successfully parsed {len(content_nodes)} core structural nodes.")
            print("No critical syntax formatting anomalies detected by latextools.")
            
    except Exception as e:
        print(f"Lint Fail: Anomalies or structural syntax errors found!\nDetails: {e}", file=sys.stderr)
        sys.exit(2)

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: latexlint <path_to_latex_file.tex>")
        sys.exit(1)
        
    lint_latex_file(sys.argv[1])
EOF

# 5. Make the lint tool executable
chmod +x "$CLI_PATH"

echo -e "${GREEN}[+] Deployment successful!${NC}"
echo -e "Ensure ${GREEN}$BIN_DIR${NC} is in your system PATH environment variable."
echo -e "Usage: ${GREEN}latexlint your_file.tex${NC}"
