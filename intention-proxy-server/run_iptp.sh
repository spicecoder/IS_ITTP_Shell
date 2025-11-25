#!/bin/bash
set -euo pipefail
echo "Installing server dependencies..."
npm install
echo "Starting Intention Proxy server..."
SHARED_SECRET=${SHARED_SECRET:-topsecret} npm start
