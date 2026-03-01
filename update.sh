#!/bin/bash
set -e

echo "=> Upgrading dependencies..."
make upgrade

echo "=> Update complete! Run 'ft help' to test."
