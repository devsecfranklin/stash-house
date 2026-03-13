# Stash House: GEMINI.md

This file provides architectural context and development guidelines for the Stash House project.

## Project Overview

**Stash House** is a multi-layer framework designed to reconstitute disciplined secret use under conditions of local disorder. It moves beyond passive detection to a model of "Authority Recovery," bridging local workstation hygiene with decentralized transport layers (Nostr) and federated identity brokers (FOKS).

## Core Architecture

- **Materialization Boundary:** The critical transition point where encrypted, stored authority becomes live runtime authority.
- **Authority Recovery:** A hybrid architecture integrating enterprise standards (Kerberos, LDAP) with sovereign, distributed backends.
- **Components:**
    - **FOKS:** Federated Identity Broker.
    - **Nostr:** Decentralized transport layer.
    - **GKE/Terraform:** Cloud infrastructure for deployment.

## Tech Stack

- **Languages:** Go (1.25+), Python (3.9+), Bash, LaTeX.
- **Infrastructure:** Terraform (AWS, GCP), Docker, Kubernetes (GKE).
- **Build System:** GNU Autotools (`autoconf`, `automake`).
- **Identity & Secrets:** Kerberos, LDAP, GPG, Nostr, FOKS.

## Development Workflow

### Initial Setup

The project uses a bootstrap script to configure the environment, including GNU Autotools and Go dependencies.

```bash
./bootstrap.sh
```

This script will:
1. Install necessary Debian packages.
2. Initialize GNU Autotools (`aclocal`, `autoreconf`, `automake`, `./configure`).
3. Initialize and update Go modules.

### Build Commands

- **Build Go binaries:** `go build ./cmd/...`
- **Build LaTeX White Paper:** `make paper`
- **Build LaTeX Slides:** `make slides`
- **Clean build artifacts:** `make clean`

### Linting & Formatting

- **Linting:** `make lint` (runs `shfmt` and `lacheck`).
- **Go Formatting:** `go fmt ./...`

## Project Structure

- `cmd/`: Go command-line tools.
- `docs/`: White papers, slides, and architectural diagrams (LaTeX and Markdown).
- `terraform/`: Infrastructure as Code for AWS and GCP.
- `container/`: Dockerfiles and container configurations.
- `demo/`: Setup scripts and walkthroughs for demonstrating the framework.
- `aclocal/`: Custom M4 macros for Autotools (specifically for LaTeX).

## Development Mandates

- **Surgical Updates:** Maintain the GNU Autotools structure when adding new files or directories.
- **Secret Safety:** Never commit secrets. Use the framework's own "Authority Recovery" principles for managing sensitive data.
- **LaTeX Documentation:** Scientific papers and presentations are core project artifacts. Ensure LaTeX builds are verified after documentation changes.
