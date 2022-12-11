#!/bin/bash

set -e

if [[ $# -ne 1 ]]; then
    echo -e "1) Build Debug mode\n2) Build Release mode"
    exit 0
fi

mode=$1
if [[ $mode -eq 1 ]]; then

    ./switch-physx-mode.sh 1
    go build .
    echo "Debug build finished"

elif [[ $mode -eq 2 ]]; then

    ./switch-physx-mode.sh 2
    go build -tags "nmage_release,physx_release" .
    echo "Release build finished"

else

    echo "Unknown build option. Please select 1 for a Debug build or 2 for a Release build"
    exit 1

fi
