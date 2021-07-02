UPDATE transactions
SET fee = $2,
    signature = $3,
    public_key = $4,
    outcome = $5,
    receipt = $6
WHERE hash = $1
