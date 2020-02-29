package ledger

// Account represents a cryptographic identity holding value on the ledger
// 		TODO: to specify
type Account struct {
}

// Transaction reprensent an attempt at transfering value from one account to another.
// 		A "Transaction" object comes with no validity guarantees.
// 		It can represent a canonized tx, a tx waiting for confs, a mempool tx.
type Transaction struct {
	origin Account
	destination Account
	value float
}