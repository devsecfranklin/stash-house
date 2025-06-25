#!/usr/bin/env bash
#
# SPDX-FileCopyrightText: 2023 DE:AD:10:C5 <franklin@dead10c5.org>
#
# SPDX-License-Identifier: GPL-3.0-or-later
#
# v0.2 - A tool to inspect and aggressively clean a Git repository.

# --- Script Configuration ---
# Exit immediately if a command exits with a non-zero status.
set -e
# Treat unset variables as an error.
set -u
# The return value of a pipeline is the status of the last command to exit
# with a non-zero status, or zero if no command failed.
set -o pipefail

# --- Color and Logging Functions ---
# Using tput for compatibility and to check if the terminal supports color.
if tput setaf 1 &> /dev/null; then
    RED=$(tput setaf 1)
    GREEN=$(tput setaf 2)
    YELLOW=$(tput setaf 3)
    CYAN=$(tput setaf 6)
    BOLD=$(tput bold)
    NC=$(tput sgr0) # No Color
else
    RED=""
    GREEN=""
    YELLOW=""
    CYAN=""
    BOLD=""
    NC=""
fi

# Centralized logging functions for consistent output.
log_info() { echo -e "${CYAN}==>${NC} ${BOLD}$1${NC}"; }
log_success() { echo -e "${GREEN}==>${NC} ${BOLD}$1${NC}"; }
log_warn() { echo -e "${YELLOW}WARN:${NC} $1"; }
log_error() { >&2 echo -e "${RED}ERROR:${NC} $1"; } # Errors to stderr

# --- Check Functions ---

# Check if the script is being run as root.
check_if_root() {
    if [[ "${EUID:-$(id -u)}" -eq 0 ]]; then
        log_error "This script should not be run as root."
        exit 1
    fi
}

# Check if the current directory is a Git repository.
check_if_git_repo() {
    if ! git rev-parse --is-inside-work-tree &> /dev/null; then
        log_error "This script must be run from within a Git repository."
        exit 1
    fi
}

# Determines the default remote branch (main or master).
get_default_branch() {
    # This is the most reliable way to find the default branch name.
    git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@'
}

# Checks for modern branch conventions (main vs. master).
check_branch_conventions() {
    log_info "Checking branch conventions..."
    local default_branch
    default_branch=$(get_default_branch)

    if [[ "$default_branch" == "master" ]]; then
        log_warn "Your default branch is 'master'. Consider renaming it to 'main' for modern conventions."
    else
        log_success "Default branch is '$default_branch'."
    fi

    # Check if a 'master' branch exists locally when 'main' is the default.
    if [[ "$default_branch" == "main" ]] && git show-ref --quiet --verify refs/heads/master; then
        log_warn "A local 'master' branch exists. It may be a stale remnant."
    fi
}

# --- Core Logic Functions ---

# Display the size of the .git directory.
get_repo_size() {
    du -sh .git | awk '{print $1}'
}

# Perform the cleanup operations.
clean_repo() {
    log_info "Starting repository cleanup..."

    log_info "Pruning stale remote-tracking branches..."
    git remote prune origin

    log_info "Repacking objects to reduce disk space..."
    git repack -ad # -a packs all, -d removes redundant objects

    log_info "Pruning loose objects..."
    git prune

    log_info "Expiring reflog entries older than 30 days..."
    git reflog expire --expire=30.days.ago --all

    log_info "Running garbage collection..."
    git gc --aggressive --prune=now
}

# Display help message.
usage() {
    cat <<EOF
Usage: $(basename "$0") [options]

A tool to inspect and aggressively clean a Git repository.

Options:
  -f, --force    Skip the confirmation prompt and run the cleanup immediately.
  -c, --check    Run checks only; do not perform any cleanup operations.
  -h, --help     Display this help message and exit.
EOF
}

# --- Main Execution ---
main() {
    # Default options
    local FORCE=false
    local CHECK_ONLY=false

    # Parse command-line arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -f|--force)
                FORCE=true
                shift
                ;;
            -c|--check)
                CHECK_ONLY=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done

    check_if_root
    check_if_git_repo
    check_branch_conventions

    if [[ "$CHECK_ONLY" == true ]]; then
        log_info "Check-only mode enabled. No cleanup will be performed."
        exit 0
    fi

    local start_size
    start_size=$(get_repo_size)
    log_info "Current repository size: ${start_size}"

    if [[ "$FORCE" == false ]]; then
        read -p "${YELLOW}This will perform aggressive cleanup operations. Are you sure? (y/n) ${NC}" -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Cleanup cancelled by user."
            exit 1
        fi
    fi

    clean_repo

    local end_size
    end_size=$(get_repo_size)
    log_success "Cleanup complete!"
    log_info "Final repository size: ${end_size} (was ${start_size})"
}

# Run the main function with all command-line arguments.
main "$@"
