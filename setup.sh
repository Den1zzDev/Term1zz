#!/bin/sh
set -eu

# Term1zz setup wrapper
# Delegates to POSIX installer in scripts/install.sh

SCRIPT_DIR="$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)"

if [ -f "${SCRIPT_DIR}/scripts/install.sh" ]; then
    exec /bin/sh "${SCRIPT_DIR}/scripts/install.sh" "$@"
fi

# If piped or run without local clone
curl -fsSL https://raw.githubusercontent.com/Den1zzDev/Term1zz/main/scripts/install.sh | /bin/sh -s -- "$@"
