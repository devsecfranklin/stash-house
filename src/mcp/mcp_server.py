from pathlib import Path
from fastmcp import FastMCP
from fastmcp.resources import FileResource

mcp = FastMCP("research.bitsmasher.net")


@mcp.tool
def greet(name: str) -> str:
    return f"Hello, {name}!"

@mcp.resource("file://README.md")
def get_file_content() -> str:
    """Returns the content of the specified file."""
    with open("README.md", "r") as file:
        return file.read()

if __name__ == "__main__":
    mcp.run(transport="http", host="0.0.0.0", port=8989)
