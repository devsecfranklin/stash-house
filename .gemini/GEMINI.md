# Workspace Mandates

## Context
This is the root workspace containing multiple sub-projects (lab-franklin, website, writing, etc.). This configuration applies globally across all sub-projects unless overridden by a more specific `GEMINI.md`.

## Project Stack
- **Languages:** Go, Python, Bash, LaTeX, C++, Perl.
- **Tools:** GNU Autotools, Terraform, Docker, Kubernetes, GCP.

## Development Workflow
- **Submodules:** This repository uses git submodules. Always use `--recursive` when performing git operations that affect submodules.
- **Bootstrapping:** Use the root `bootstrap.sh` or project-specific bootstrap scripts to initialize environments.
- **Build System:** GNU Autotools (`./configure`, `make`) is common across several projects.

## Engineering Standards
- **Secret Safety:** Strictly avoid committing secrets. Use local environment variables or project-specific secret management.
- **Surgical Updates:** When modifying code, adhere strictly to the local style and conventions of the specific sub-project.
- **Documentation:** Maintain READMEs and LaTeX documentation as core artifacts.
