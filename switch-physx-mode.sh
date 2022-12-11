#!/bin/bash

set -e

if [[ $# -ne 1 ]]; then
    echo -e "1) Checked mode\n2) Release mode"
    exit 0
fi

mode=$1
physxBinDir="../physx/physx/bin/win.x86_64.vc142.mt"
physxCBinDir="../physx-c/x64"
if [[ $mode -eq 1 ]]; then

    physxCheckedBinDir="$physxBinDir/checked"
    cp "$physxCheckedBinDir/PhysX_64.dll" "$physxCheckedBinDir/PhysXCommon_64.dll" "$physxCheckedBinDir/PhysXFoundation_64.dll" .

    physxCCheckedBinDir="$physxCBinDir/Checked"
    cp "$physxCCheckedBinDir/physx-c.dll" .

    echo "Switched PhysX to Checked mode"

elif [[ $mode -eq 2 ]]; then

    physxReleaseBinDir="$physxBinDir/release"
    cp "$physxReleaseBinDir/PhysX_64.dll" "$physxReleaseBinDir/PhysXCommon_64.dll" "$physxReleaseBinDir/PhysXFoundation_64.dll" .

    physxCReleaseBinDir="$physxCBinDir/Release"
    cp "$physxCReleaseBinDir/physx-c.dll" .

    echo "Switched PhysX to Release mode"

else

    echo "Unknown mode. Please select 1 or 2"
    exit 1

fi
