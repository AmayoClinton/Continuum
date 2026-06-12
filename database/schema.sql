DROP TABLE IF EXISTS vaults;
DROP TYPE IF EXISTS vault_status;

CREATE TYPE vault_status AS ENUM ('ACTIVE', 'DORMANT');

CREATE TABLE vaults (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alias VARCHAR(255) NOT NULL,
    beneficiary_pubkey TEXT NOT NULL,
    encrypted_payload TEXT NOT NULL,          
    check_in_interval_seconds INT NOT NULL,  
    last_check_in_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status vault_status DEFAULT 'ACTIVE'
);