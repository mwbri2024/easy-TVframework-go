#!/bin/sh

# Exit on non defined variables and on non zero exit codes
set -eu

# Execute the compiled application with any provided arguments
exec /app/itv "$@"