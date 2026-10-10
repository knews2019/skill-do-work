#!/usr/bin/env bash
# Pre-flight baseline for REQ-679 (run before the build): the two timeline probes it will move pass today with the browser lane on and are not skipped.
set -euo pipefail
exec bash "$(dirname "${BASH_SOURCE[0]}")/REQ-679-probe.sh"
