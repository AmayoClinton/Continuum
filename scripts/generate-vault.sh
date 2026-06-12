#!/bin/bash

echo "🔐 Generating automated sample cryptographic capsule..."

# This simulates a pre-packaged Base64 string payload generated from encryption.ts
MOCK_ENCRYPTED_PAYLOAD="eyJlcGhlbWVyYWxQdWJLZXkiOiIwM2U5YSIsIml2IjoiYWYzZSIsImF1dGhUYWciOiI5OWJiIiwiY2lwaGVydGV4dCI6IjEyM2EifQ=="
MOCK_PUBKEY="02e9a2631247d5124b893a71b25076eefc432d56a29851a7eef1109bcfa0329a1d"

echo "📡 Transmitting registration payload to local Fiber API endpoint..."

RESPONSE=$(curl -s -X POST http://localhost:8080/api/vaults \
  -H "Content-Type: application/json" \
  -d "{
    \"alias\": \"CLI-Automated-Vault\",
    \"beneficiary_pubkey\": \"$MOCK_PUBKEY\",
    \"encrypted_payload\": \"$MOCK_ENCRYPTED_PAYLOAD\",
    \"check_in_interval_seconds\": 15
  }")

echo "📦 API Response Server Output:"
echo "$RESPONSE"