# Docker Multi-Platform Build Configuration

This document explains how MMock's Docker images are built with multi-platform support and attestations.

## Overview

MMock uses GoReleaser to build Docker images for multiple platforms:
- `linux/amd64` (x86_64)
- `linux/arm64` (ARM 64-bit)

Images are automatically built and pushed to Docker Hub on every push to the `master` branch.

## Build Architecture

### Components

1. **GoReleaser** - Orchestrates the entire build process
2. **Docker Buildx** - Builds multi-platform Docker images
3. **docker-container driver** - Ensures attestations (SBOM/provenance) are preserved
4. **Docker Hub** - Final image registry

### Build Flow

```
┌─────────────┐
│ Push to     │
│ master      │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────────┐
│ GitHub Actions Workflow                         │
│                                                  │
│  1. Setup Docker Buildx (docker-container)      │
│  2. Build Go binaries (all platforms)           │
│  3. Build Docker images:                        │
│     - linux/amd64 → jordimartin/mmock:vX.X.X-amd64 │
│     - linux/arm64 → jordimartin/mmock:vX.X.X-arm64 │
│  4. Create multi-platform manifest              │
│  5. Attach attestations (SBOM + Provenance)     │
│  6. Push to Docker Hub                          │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
          ┌────────────────┐
          │  Docker Hub    │
          │                │
          │  Images with:  │
          │  ✓ Multi-arch  │
          │  ✓ SBOM        │
          │  ✓ Provenance  │
          └────────────────┘
```

## Local Development

### Prerequisites

Ensure you have the `docker-container` driver configured:

```bash
# Check current builders
docker buildx ls

# If you don't have a docker-container builder, create one:
docker buildx create --name goreleaser --driver docker-container --bootstrap --use

# Enable multi-platform support
docker run --privileged --rm tonistiigi/binfmt --install all
```

### Building Locally

```bash
# Snapshot build (no push to registry)
goreleaser release --snapshot --clean --skip=publish

# Test the images
docker run --rm jordimartin/mmock:latest-amd64 --help
docker run --rm jordimartin/mmock:latest-arm64 --help

# Inspect the multi-platform manifest
docker manifest inspect jordimartin/mmock:latest
```

## GitHub Actions Configuration

### Prerequisites: Environment and Secrets

Both workflows require the **Build** environment with Docker Hub credentials:

**Required GitHub Secrets (in "Build" environment):**
- `DOCKERHUB_USERNAME` - Your Docker Hub username
- `DOCKERHUB_TOKEN` - Docker Hub access token (not password!)

**Setup:**
1. Go to: Repository → Settings → Environments → Build
2. Add the required secrets
3. Use Docker Hub **access tokens**, not passwords:
   - Docker Hub → Account Settings → Security → New Access Token
   - Create with "Read & Write" permissions
   - Add to GitHub as `DOCKERHUB_TOKEN`

### Release Workflow

Located at `.github/workflows/release.yml`, the workflow:

1. **Sets up Docker Buildx** with docker-container driver
   ```yaml
   - name: Set up Docker Buildx
     uses: docker/setup-buildx-action@v3
     with:
       driver: docker-container
       platforms: linux/amd64,linux/arm64
       install: true
   ```

2. **Verifies the builder** is correctly configured
3. **Runs GoReleaser** which handles the entire build and push process
4. **Verifies attestations** were created and attached

### Test Workflow

A manual test workflow (`.github/workflows/test-docker-build.yml`) is available for testing without creating releases:

```bash
# Trigger via GitHub UI:
# Actions → Test Docker Build (Manual) → Run workflow
```

This workflow:
- Uses the **Build** environment (same secrets as release workflow)
- Builds snapshot images without pushing to Docker Hub
- Verifies docker-container driver is active
- Tests both amd64 and arm64 images
- Confirms multi-platform setup works correctly

**Note:** While this workflow logs into Docker Hub (to match production conditions), it uses `--skip=publish` so no images are pushed.

## Dockerfile

Located at `Dockerfile`, uses the `${TARGETPLATFORM}` variable:

```dockerfile
ARG TARGETPLATFORM
COPY ${TARGETPLATFORM}/mmock /usr/local/bin/mmock
```

**Key Point:** The braces `${}` around `TARGETPLATFORM` are **required** for proper variable expansion. Without them, Docker may misinterpret the path.

## GoReleaser Configuration

Located at `.goreleaser.yml`:

```yaml
dockers_v2:
  - id: mmock
    dockerfile: Dockerfile
    platforms:
      - linux/amd64
      - linux/arm64
    images:
      - "jordimartin/mmock"
    tags:
      - "{{ .Tag }}"
      - "latest"
```

GoReleaser's `dockers_v2` feature automatically:
- Builds separate images per platform
- Creates a multi-platform manifest combining both
- Preserves attestations when using docker-container driver

## Attestations

### What Are Attestations?

Attestations are cryptographic metadata attached to Docker images:

1. **SBOM (Software Bill of Materials)** - Lists all components and dependencies
2. **Provenance** - Documents who built the image, when, and how

### Verifying Attestations

```bash
# Inspect an image from Docker Hub
docker buildx imagetools inspect jordimartin/mmock:latest

# Look for entries like:
#   vnd.docker.reference.type: attestation-manifest
```

### Why Attestations Matter

- **Security Scanning**: Tools can identify vulnerabilities without downloading images
- **Supply Chain Verification**: Users can verify image authenticity
- **Compliance**: Meets security requirements for enterprise environments
- **Transparency**: Shows exactly what's in your images

## Troubleshooting

### "invalid docker buildx driver" Warning

If you see this warning from GoReleaser:

```
invalid docker buildx driver - docker buildx is using the docker
driver, which isn't tested and may cause issues.
```

**Solution:** Switch to docker-container driver:
```bash
docker buildx create --name goreleaser --driver docker-container --use
```

### Images Build But Attestations Missing

**Cause:** Using the default `docker` driver instead of `docker-container`

**Solution:**
- Local: Run `docker buildx use goreleaser` (assumes you created the builder)
- GitHub Actions: Ensure `docker/setup-buildx-action@v3` is in your workflow

### Multi-Platform Build Fails

**Cause:** Missing QEMU emulation for cross-platform builds

**Solution:**
```bash
docker run --privileged --rm tonistiigi/binfmt --install all
```

## Verification Checklist

Before pushing to master, verify:

- [ ] `goreleaser check` passes
- [ ] `docker buildx ls` shows docker-container driver active
- [ ] Local snapshot build succeeds
- [ ] Both amd64 and arm64 images run correctly
- [ ] `docker manifest inspect` shows both platforms
- [ ] GitHub Actions workflow includes `setup-buildx-action`

## References

- [GoReleaser Docker v2 Documentation](https://goreleaser.com/customization/dockers_v2/)
- [Docker Buildx Drivers](https://docs.docker.com/build/builders/drivers/)
- [Docker Attestations](https://docs.docker.com/go/attestations/)
- [Docker setup-buildx-action](https://github.com/docker/setup-buildx-action)

## Maintenance

### Updating Go Version

Update in two places:
1. `go.mod` - The module's Go version
2. `.github/workflows/release.yml` - Uses `go-version-file: go.mod`

### Updating Base Image

The Dockerfile uses `alpine:3.20`. To update:
1. Change in `Dockerfile`
2. Test locally with snapshot build
3. Verify image sizes haven't increased significantly

### Adding New Platforms

To add new platforms (e.g., `linux/386`):

1. Add to `.goreleaser.yml`:
   ```yaml
   builds:
     - goarch:
         - amd64
         - arm64
         - 386  # Add new architecture

   dockers_v2:
     - platforms:
         - linux/amd64
         - linux/arm64
         - linux/386  # Add new platform
   ```

2. Update GitHub Actions:
   ```yaml
   - name: Set up Docker Buildx
     with:
       platforms: linux/amd64,linux/arm64,linux/386
   ```

3. Test locally before pushing
