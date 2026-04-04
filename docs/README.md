# Stash

1. Consolidation & De-duplication

```sh
* [ ] Merge pass Guides: Consolidate docs/password-store-readme.txt, docs/stash.md, and the usage section of docs/README.md into a single
    docs/guides/local-stewardship.md.
* [ ] Clean up Root README: Update docs/README.md to serve as a high-level table of contents for the documentation folder, rather than containing
    fragmented code snippets.
* [ ] Formalize Scratchpads: Move raw URLs and command snippets from docs/nostr.txt and docs/latex.txt into structured documentation (e.g.,
    docs/architecture/nostr.md and docs/guides/building-from-source.md).
```

2. Architectural Clarity

```sh
* [ ] Create a Developer Guide: Add a docs/DEVELOPMENT.md explaining the Go project structure (cmd/, internal/), how to add new "ProtectLocal"
    operators, and how to run the test suite.
* [ ] Document the Materialization Boundary: The LaTeX paper identifies this as a "critical transition point." Create a dedicated markdown file
    (docs/architecture/materialization.md) that explains this concept for engineers, including the JIT injection patterns mentioned.
* [ ] Update Diagrams: The docs/stash.drawio file should be exported to a high-resolution PNG/SVG and embedded in the markdown docs for easier
    viewing without a specific editor.
```

3. Workflow & Tooling

```sh
* [ ] Integrate Demo Docs: Link the demo/ scripts clearly within the documentation. Create a docs/guides/walkthrough.md that guides a user through
    demo_setup.sh and the 01-03 sequence.
* [ ] Infrastructure Docs: Add documentation for the terraform/ layer. While the paper mentions GKE/FOKS, there is no "Operator's Guide" on how to
    actually deploy the federated backend.
* [ ] Keyringer Integration: The slides mention incorporating keyringer. Decide if this is a core part of the framework and document the integration
    path or remove the placeholder if it's no longer planned.
```

4. Maintenance & Standards

```sh
* [ ] Standardize Naming: Ensure the project is consistently referred to as "Stash House" (the framework) rather than just "Stash" (which is often
    confused with the tool pass or generic stashing).
* [ ] Automate LaTeX Builds: Ensure the Makefile.am in docs/paper and docs/slides is integrated into the root Makefile so make docs builds all PDFs.
* [ ] Verify Bibliography: Update docs/paper/bib.bib to ensure all URLs are current, particularly the FOKS whitepaper link.
```

5. Suggested Directory Structure Refactor

```sh
docs/
├── architecture/           # High-level design (Materialization, Nostr, FOKS)
├── guides/                 # How-to guides (Local Stewardship, Deployment, Demos)
├── paper/                  # Formal LaTeX White Paper
├── slides/                 # Presentation materials (Consolidate ec/ and bsides-cos/)
├── DEVELOPMENT.md          # Contributor guide
└── README.md               # Documentation Index
```

## Usage

To import a multiline secret into an env var from `pass` database:

```sh
export GOOGLE_APPLICATION_CREDENTIALS=(pass show gcp-gcs-pso-f9322c73b9bf.json| string split0)
```
