-- +goose Up
CREATE TABLE transactions (
  id         SERIAL NOT NULL PRIMARY KEY,
  hash       TEXT NOT NULL,
  block_hash TEXT NOT NULL,
  height     INTEGER NOT NULL,
  time       TIMESTAMP WITH TIME ZONE NOT NULL,
  signer     VARCHAR,
  signer_key VARCHAR,
  receiver   VARCHAR,
  signature  VARCHAR,
  amount     VARCHAR,
  gas_burnt  VARCHAR,
  success    BOOLEAN,
  actions    JSONB,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_transactions_block_hash
  ON transactions(block_hash);

CREATE UNIQUE INDEX idx_transactions_hash
  ON transactions(hash);

CREATE INDEX idx_transactions_height
  ON transactions(height);

CREATE INDEX idx_transactions_time
  ON transactions(time);

CREATE INDEX idx_transactions_signer
  ON transactions(signer);

CREATE INDEX idx_transactions_receiver
  ON transactions(receiver);

-- +goose Down
DROP TABLE transactions;
