# study-go

This repository contains a small sample application built with [Fiber](https://github.com/gofiber/fiber). The `scripts/setup.sh` helper installs Go and downloads the module dependencies. For Codex environments, `.codex/setup.sh` automatically runs this helper during provisioning so that dependencies are available before the network is disabled.

## Setup

Run the setup script from the repository root:

```bash
./scripts/setup.sh
```

Then build and run the app:

```bash
go run ./...
```

Run tests with:

```bash
go test ./...
```

## Codex environment

This repository provides `.codex/setup.sh`. When using Codex, configure this as
your environment's setup script so dependencies are installed automatically.
