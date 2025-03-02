#!/usr/bin/env bash

set -ex


test_version() {
    local version="$1"
    # If version is "latest", run without --version flag
    local ver_flag="--version=${version}"
    if [[ "$version" == "latest" ]]; then
      ver_flag=""
    fi

    # Status should fail with OLM not installed
    commandoutput=$(operator-sdk olm status 2>&1 || true)
    echo $commandoutput | grep -F "Failed to get OLM status"

    # Uninstall should fail with OLM not installed
    commandoutput=$(operator-sdk olm uninstall 2>&1 || true)
    echo $commandoutput | grep -F "Failed to uninstall OLM"

    # Install should succeed with nothing installed
    commandoutput=$(operator-sdk olm install $ver_flag 2>&1)
    echo $commandoutput | grep -F "Successfully installed OLM version \\\"${version}\\\""

    # Install should fail with OLM Installed
    commandoutput=$(operator-sdk olm install $ver_flag 2>&1 || true)
    echo $commandoutput | grep -F "Failed to install OLM version \\\"${version}\\\": detected existing OLM resources: OLM must be completely uninstalled before installation"

    # Status should succeed with OLM installed
    commandoutput=$(operator-sdk olm status 2>&1)
    echo $commandoutput | grep -F "Successfully got OLM status"

    # Uninstall should succeed with OLM installed
    commandoutput=$(operator-sdk olm uninstall 2>&1)
    echo $commandoutput | grep -F "Successfully uninstalled OLM"
}

test_version "0.28.0"
test_version "0.17.0" # Check installation of OLM for locally stored version of binaries
