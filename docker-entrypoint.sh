#!/bin/sh
set -e

# this if will check if the first argument is a flag
if [ "$#" -eq 0 ] || [ "${1#-}" != "$1" ]; then
    set -- getFailedStepsLogs "$@"
fi

# check for the expected command
if [ "$1" = 'getFailedStepsLogs' ]; then
    exec getFailedStepsLogs "$@"
fi

# else default to run whatever the user wanted like "bash" or "sh"
exec "$@"

