# Installation

## Prerequisites

- **Go**: Version 1.21 or higher
- **Git**: For cloning the repository

## Install via go install

The easiest way to install goxcel:

```bash
go install github.com/ryo-arima/goxcel/cmd/goxcel@latest
```

This will install the `goxcel` binary to `$GOPATH/bin` (usually `~/go/bin`).

### Verify Installation

```bash
goxcel --version
```

## Build from Source

### Clone Repository

```bash
git clone https://github.com/ryo-arima/goxcel.git
cd goxcel
```

### Build Binary

```bash
make build
```

The binary will be created at `.bin/goxcel`.

### Install Locally

```bash
make install
```

Or copy manually:

```bash
cp .bin/goxcel /usr/local/bin/
# or
cp .bin/goxcel $GOPATH/bin/
```

## Docker (Optional)

Build Docker image:

```bash
docker build -t goxcel .
```

Run with Docker:

```bash
docker run -v $(pwd):/workspace goxcel \
  generate \
  --template /workspace/template.gxl \
  --data /workspace/data.json \
  --output /workspace/output.xlsx
```

## Verify Setup

Test with sample files:

```bash
# Navigate to repository
cd goxcel

# Generate sample
.bin/goxcel generate \
  --template .etc/sample.gxl \
  --data .etc/sample.json \
  --output sample.xlsx

# Check output
ls -lh sample.xlsx
```

## Environment Setup

### Add to PATH

If `goxcel` is not found, add Go bin to PATH:

**Linux/macOS (bash/zsh):**
```bash
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
# or for zsh
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.bashrc  # or source ~/.zshrc
```

**Windows (PowerShell):**
```powershell
$env:PATH += ";$env:USERPROFILE\go\bin"
```

### Configure Logging (Optional)

Set environment variables for logging:

```bash
export GOXCEL_LOG_LEVEL=DEBUG
export GOXCEL_LOG_STRUCTURED=true
```

## Troubleshooting

### Command Not Found

If `goxcel` command is not found:

1. Check Go is installed: `go version`
2. Verify GOPATH: `go env GOPATH`
3. Check binary location: `ls $GOPATH/bin/goxcel`
4. Ensure PATH includes Go bin directory

### Build Errors

If build fails:

1. Update Go: `go version` should be 1.21+
2. Clean and rebuild:
   ```bash
   go clean -cache
   make clean
   make build
   ```

### Permission Denied

On Linux/macOS, make binary executable:

```bash
chmod +x .bin/goxcel
# or for installed binary
chmod +x $GOPATH/bin/goxcel
```

## Next Steps

- [Quick Start Guide](./quick-start.md) - Create your first template
- [Basic Concepts](./concepts.md) - Understand GXL fundamentals
- [Specification](../specification/overview.md) - Detailed format reference
