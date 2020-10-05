package model

type PublicDaemon struct {
	Public string `json:"public_key"`
}
type Trans struct {
	Txn string `json:"txn"`
}
type TransID struct {
	ID int64 `json:"id"`
}

type TransIDs struct {
	IDs []int64 `json:"ids"`
}

type Signature struct {
	Message   []Trans `json:"message"`
	Signature string  `json:"signature"`
}
