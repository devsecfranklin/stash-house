# MCP (Model Context Protocol) - Technical Deep Dive (Updated)

This document provides a technical overview of the MCP setup on the Stargate host and summarizes recent architectural updates.

## 1. MCP Architecture and Goals

The MCP bridges high-level LLM reasoning with Minecraft server administrative tools, optimized for low-latency I/O and performance.

## 2. Recent Updates (April 14, 2026)

*   **Directory Standardization**: Renamed `src/mcp-server/` to `src/mcp_server/` to align with Python package naming conventions.
*   **Infrastructure Refactoring**:
    *   Updated `mcp-servers.json` to reflect absolute workspace paths (`/mnt/clusterfs2/workspace/gaming/game-skynet-minecraft`).
    *   Migrated MCP service execution to the dedicated project virtual environment (`.venv/bin/python3`).
*   **Daemon Orchestration**:
    *   Updated `stargate-daemon.service` with correct working directory and executable paths.
    *   Patched `deploy.sh` to correctly locate the service file within the new workspace structure.
*   **Test Integration**:
    *   Introduced `test/test_mcp_stargate.py` for functional validation of the MCP utility tools.
*   **Documentation**:
    *   Synced `README.md` and `src/servers/stargate/README.md` to reflect the new directory structure and the `make stargate` build target.

## 3. Repository Structure & Emerald Integration

*   `src/mcp_server/`: Contains the Python scripts and configuration for MCP servers.
    *   `mcp-server.py`: Base utility server.
    *   `mcp-servers.json`: Configuration registry for MCP servers.
    *   `schematic_generator_orchestrator.py`: LLM orchestration logic for Emerald Mirror mutations.
*   **Emerald Mirror Tooling:**
    *   `spatial-snapshot-tool`: Ingests voxel data from Chonk for mutation.
    *   `vertex-ai-link`: Routes data to Vertex AI and retrieves greebled deltas.
    *   `BUC-throttler`: Ensures RCON deployments stay within 20 TPS limits.
*   `test/test_mcp_stargate.py`: Functional tests for MCP services.
*   `Makefile.am`: Added `stargate` target for environment synchronization.

---
*Created for theDevilsVoice | Last Updated: May 8, 2026*
