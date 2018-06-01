#!/bin/sh
set -e

# this if will check if the first argument is a flag
# and add getFailedStepsLogs
if [ "$#" -eq 0 ] || [ "${1#-}" != "$1" ]; then
    # prefix with getFailedStepsLogs
    set -- getFailedStepsLogs "$@"
fi

exec "$@"

