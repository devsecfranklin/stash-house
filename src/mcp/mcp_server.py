from mcp.server.fastmcp import FastMCP

mcp = FastMCP("stargate", auth_callback=auth_token)

@mcp.resource("resource://greeting")
def get_greeting() -> str:
    """Provides a simple greeting message."""
    return "Hello from FastMCP Resources!"


@mcp.resource("resource://context")
def get_context(filename: str = "context.md") -> str:
    """
    Returns the contents of a markdown file.

    Parameters
    ----------
    filename: str
        Name of the markdown file relative to the server’s base directory.
        Defaults to ``context.md``.

    Security considerations
    -----------------------
    * The path is resolved against ``BASE_DIR`` – any ``..`` components are
      stripped to prevent directory‑traversal attacks.
    * Only files ending with ``.md`` are allowed.
    """
    BASE_DIR = pathlib.Path(__file__).parent.resolve()

    # Prevent “../secret.txt” tricks
    safe_name = pathlib.Path(filename).name  # drop any directory parts
    if not safe_name.lower().endswith(".md"):
        raise ValueError("Only *.md files may be requested.")

    target = BASE_DIR / safe_name
    if not target.is_file():
        raise FileNotFoundError(f"{safe_name!r} not found on the server.")

    return target.read_text(encoding="utf-8")


if __name__ == "__main__":
    mcp.run(transport="tcp", host="0.0.0.0", port=8989)
