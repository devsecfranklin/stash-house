#
# Makefile for building and pushing the Franklin Resume container images using Podman.
#

# --- Configuration ---
# Centralize your image names and tags here for easy updating.
REGISTRY        := ghcr.io/devsecfranklin
IMAGE_NAME      := franklin-resume
DEFAULT_TAG     := latest

# podman-compose uses the current directory name as the project name by default.
# We capture it here to reliably tag the correct image.
PROJECT_NAME    := $(shell basename $(CURDIR))

# Fully qualified image names are constructed from the variables above.
FULL_IMAGE_NAME       := $(REGISTRY)/$(IMAGE_NAME)
FULL_IMAGE_NAME_AMD64 := $(REGISTRY)/$(IMAGE_NAME)-amd64

# --- Phony Targets ---
# Declare targets that are not files.
.PHONY: all build build-amd64 build-all check-env push push-amd64 push-all run clean prune

# --- Main Targets ---

# Default target: running 'make' will just build the primary image.
all: build

# Build the primary ARM/multi-arch image.
# The 'podman tag' step is no longer needed because the image name is in compose.yml
build:
	@echo "--> Building service: franklin-resume (with podman-compose)..."
	podman-compose build franklin-resume

# Build the dedicated AMD64 image.
# The 'podman tag' step is no longer needed.
build-amd64:
	@echo "--> Building service: amd64 (with podman-compose)..."
	podman-compose build amd64

# Target to build both images at once.
build-all: build build-amd64

# Check for the required GitHub Container Registry PAT.
# NOTE: For Podman, you should first log in with `podman login ghcr.io -u <username> --password-stdin <<< $$CR_PAT`
check-env:
	@test -n "$${CR_PAT}" || (echo "ERROR: CR_PAT environment variable is not set." && exit 1)

# Push the primary image to the container registry.
push: check-env build
	@echo "--> Pushing $(FULL_IMAGE_NAME):$(DEFAULT_TAG)"
	podman push $(FULL_IMAGE_NAME):$(DEFAULT_TAG)

# Push the AMD64 image to the container registry.
push-amd64: check-env build-amd64
	@echo "--> Pushing $(FULL_IMAGE_NAME_AMD64):$(DEFAULT_TAG)"
	podman push $(FULL_IMAGE_NAME_AMD64):$(DEFAULT_TAG)

# Push both images.
push-all: push push-amd64

# Run the locally built primary image interactively.
# --rm automatically cleans up the container on exit.
run: build
	@echo "--> Running image $(FULL_IMAGE_NAME):$(DEFAULT_TAG) interactively..."
	podman run --rm -it $(FULL_IMAGE_NAME):$(DEFAULT_TAG) bash

# Clean up the specific images created by this Makefile.
# Uses '|| true' to prevent errors if the images don't exist.
clean:
	@echo "--> Removing local Podman images..."
	@podman image rm -f $(FULL_IMAGE_NAME):$(DEFAULT_TAG) || true
	@podman image rm -f $(FULL_IMAGE_NAME_AMD64):$(DEFAULT_TAG) || true
	@podman image rm -f $(PROJECT_NAME)_franklin-resume:latest || true
	@podman image rm -f $(PROJECT_NAME)_amd64:latest || true

# A more aggressive cleanup. Use with caution.
prune:
	@echo "WARNING: This will remove all stopped containers, networks, and dangling images."
	podman system prune -f