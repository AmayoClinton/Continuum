#!/bin/bash
set -e

echo "⏳ Waiting for bitcoind container interface to stabilize..."
sleep 3

# Define a short helper function to execute RPC interactions quickly
btc-cli() {
  docker exec continuum-bitcoind bitcoin-cli -regtest -rpcuser=continuum -rpcpassword=finality "$@"
}

echo "🪙 Generating developer regtest wallet..."
btc-cli createwallet continuum_miner || true

echo "⛏️ Mining 101 initial blocks to mature block rewards..."
MINING_ADDR=$(btc-cli getnewaddress)
btc-cli generatetoaddress 101 "$MINING_ADDR"

echo "✅ Chain state synchronized. Current Block Count: $(btc-cli getblockcount)"