package domain

type RedeemAllocation struct {
	CreditPoolID    string
	SourceID        string
	Amount          int64
	RemainingAmount int64
}

// var allowedSubTypes = map[common.TransactionType][]common.TransactionSubType{
// 	common.TransactionType_CREDIT: {
// 		common.TransactionSubType_EARNING,
// 		common.TransactionSubType_ADJUSTMENT,
// 		common.TransactionSubType_TRANSFER_IN,
// 	},
// 	common.TransactionType_DEBIT: {
// 		common.TransactionSubType_REDEEM,
// 		common.TransactionSubType_EXPIRY,
// 		common.TransactionSubType_TRANSFER_OUT,
// 	},
// }

// func IsValidSubType(t common.TransactionType, sub common.TransactionSubType) bool {
// 	for _, s := range allowedSubTypes[t] {
// 		if s == sub {
// 			return true
// 		}
// 	}
// 	return false
// }

type MetaDebit struct {
	LedgerEntryID string `json:"ledger_entry_id"`
	Amount        int64  `json:"amount"`
}
