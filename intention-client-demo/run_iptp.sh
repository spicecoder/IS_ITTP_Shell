#!/bin/bash
set -euo pipefail
echo "Installing client dependencies..."
npm install
echo "Running client demo (sends two intentions)..."
npm run send
