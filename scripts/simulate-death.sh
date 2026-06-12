#!/bin/bash

# Default to a placeholder UUID if an argument isn't provided explicitly
VAULT_ID=${1:-"demo-vault-uuid-placeholder"}

echo "🕒 Triggering simulation protocol for Vault ID: $VAULT_ID"
echo "Sending network clock displacement request to Fiber runtime..."

RESPONSE=$(curl -s -X POST "http://localhost:8080/api/vaults/$VAULT_ID/warp")

if echo "$RESPONSE" | grep -q "SUCCESS"; then
    echo "✅ Time warp executed cleanly!"
    echo "Server response: $RESPONSE"
    echo "The background scheduler service will catch this and flag the vault DORMANT on its next 5-second tick."
else
    echo "❌ Execution failure. Make sure your Go backend is running on port 8080."
    echo "Server output: $RESPONSE"
fi