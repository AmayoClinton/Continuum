#!/bin/bash
set -e

echo "⚡ Initializing Lightning Network Node Dev Wallets..."

# Short helpers to interact with bitcoind and lnd containers smoothly
btc-cli() {
  docker exec continuum-bitcoind bitcoin-cli -regtest -rpcuser=continuum -rpcpassword=finality "$@"
}

ln-cli() {
  docker exec continuum-lnd lncli --network=regtest "$@"
}

# 1. Wait briefly for LND's gRPC/REST interface layer to wake up
echo "⏳ Waiting for LND container state engine to boot..."
sleep 5

# 2. Generate a new on-chain Bech32 address inside your LND node wallet
echo "📥 Generating on-chain deposit address for Continuum LND Node..."
LND_ADDRESS=$(ln-cli newaddress p2wkh | grep '"address"' | awk -F'"' '{print $4}')

if [ -z "$LND_ADDRESS" ]; then
    echo "❌ Error: Could not retrieve a valid deposit address from LND node container."
    exit 1
fi
echo "🎯 LND Target Deposit Address: $LND_ADDRESS"

# 3. Disburse testnet funds from your miner pool straight to the LND node
echo "💸 Sending 10 BTC from matured miner wallet pool to LND node..."
TX_ID=$(btc-cli sendtoaddress "$LND_ADDRESS" 10.0)
echo "📦 Transaction Broadcast ID: $TX_ID"

# 4. Mine blocks to confirm the transaction so the balance registers as spendable
echo "⛏️ Mining 6 blocks to force on-chain confirmation..."
MINER_ADDR=$(btc-cli getnewaddress)
btc-cli generatetoaddress 6 "$MINER_ADDR"

# 5. Output confirmation state summary parameters to verify wallet stability
echo "📊 Verifying node balances..."
sleep 2
WALLET_BALANCE=$(ln-cli walletbalance | grep '"confirmed_balance"' | awk -F'"' '{print $4}')

echo "=============================================================================="
echo "🎉 SUCCESS: LND Node Loaded and Funded!"
echo "💰 Current Spendable LND Confirmed Balance: $WALLET_BALANCE Satoshis"
echo "=============================================================================="