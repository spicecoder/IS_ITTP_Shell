#!/bin/bash
# Enhanced WiFi Detection Test

echo "=== Full en0 Interface Info ==="
ifconfig en0
echo ""

echo "=== Does en0 have an IP? ==="
ifconfig en0 | grep "inet " | grep -v "inet6"
echo ""

echo "=== Network Service for en0 ==="
networksetup -getinfo "Wi-Fi"
echo ""

echo "=== Current WiFi Network (alternative method) ==="
system_profiler SPAirPortDataType | grep "Current Network"
echo ""

echo "=== scutil check ==="
scutil --nwi | grep -A 10 "en0"