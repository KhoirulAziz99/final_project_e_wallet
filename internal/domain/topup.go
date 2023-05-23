package domain

import "time"

type TopUp struct {
	ID        int       `json:"id"`
	WalletId  Wallet    `json:"wallet_id"`
	Amount    int64     `json:"amount"`
	Timestamp time.Time `json:"time"`
}
