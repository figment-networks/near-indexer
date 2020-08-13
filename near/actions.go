package near

const (
	// Description for all types:
	// https://docs.near.org/docs/concepts/transaction

	ActionCreateAccount  = "CreateAccount"  // make a new account (for a person, company, contract, car, refrigerator, etc)
	ActionDeployContract = "DeployContract" // deploy a new contract (with its own account)
	ActionFunctionCall   = "FunctionCall"   // invoke a method on a contract (with budget for compute and storage)
	ActionTransfer       = "Transfer"       // transfer tokens from one account to another
	ActionStake          = "Stake"          // express interest in becoming a proof-of-stake validator at the next available opportunity
	ActionAddKey         = "AddKey"         // add a key to an existing account (either FullAccess or FunctionCall access)
	ActionDeleteKey      = "DeleteKey"      // delete an existing key from an account
	ActionDeleteAccount  = "DeleteAccount"  // delete an account (and transfer balance to a beneficiary account)
)

type Action struct {
	Type string
	Data interface{}
}

type CreateAccountAction struct {
}

type DeployContractAction struct {
	Code []byte `json:"code"`
}

type FunctionCallAction struct {
	Args       string `json:"args"`
	Deposit    string `json:"deposit"`
	Gas        int64  `json:"gas"`
	MethodName string `json:"method_name"`
}

type TransferAction struct {
	Deposit string `json:"deposit"`
}

type StakeAction struct {
	PublicKey string `json:"public_key"`
	Amount    string `json:"stake"`
}

type AddKeyAction struct {
	PublicKey string    `json:"public_key"`
	AccessKey AccessKey `json:"access_key"`
}

type DeleteKeyAction struct {
	PublicKey string `json:"public_key"`
}

type DeleteAccountAction struct {
	BeneficiaryID string `json:"beneficiary_id"`
}
