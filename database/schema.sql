DROP TABLE IF EXISTS vaults;
DROP TYPE IF EXISTS vault_status;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TYPE vault_status AS ENUM ('ACTIVE', 'DORMANT');

CREATE TABLE vaults (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alias VARCHAR(255) NOT NULL,
    beneficiary_pubkey TEXT NOT NULL,
    encrypted_payload TEXT NOT NULL,          
    check_in_interval_seconds INT NOT NULL,  
    last_check_in_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status vault_status DEFAULT 'ACTIVE',
    multisig_required INT NOT NULL DEFAULT 2,
    multisig_pubkeys TEXT[] NOT NULL DEFAULT '{}',
    multisig_address TEXT NOT NULL DEFAULT '',
    multisig_redeem_script TEXT NOT NULL DEFAULT '',
    multisig_descriptor TEXT NOT NULL DEFAULT '',
    multisig_network TEXT NOT NULL DEFAULT 'regtest',
    CHECK (multisig_required > 0),
    CHECK (array_length(multisig_pubkeys, 1) IS NULL OR multisig_required <= array_length(multisig_pubkeys, 1))
);
