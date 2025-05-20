#!/usr/bin/env bash
set -euo pipefail

# Run the repository setup helper during environment provisioning
cd "$(dirname "$0")/.."

./scripts/setup.sh
