package data

import (
	"time"
)

const (
	TypeUndefined   = RecordType(0)
	TypeCredentials = RecordType(1)
	TypeBankCard    = RecordType(2)
	TypeText        = RecordType(3)
	TypeBinary      = RecordType(4)
	TypeAny         = RecordType(5)
)

const (
	StrCredentials = "creds"
	StrBankCard    = "card"
	StrText        = "text"
	StrBinary      = "bin"
	StrAny         = "any"
)

type RecordType = int32

//go:generate easyjson -all data.go
type Credentials struct {
	Key      string `json:"key"`
	Password string `json:"password"`
}

type BankCard struct {
	Number     string `json:"number"`
	Expiration string `json:"expiration"`
	CVV        string `json:"CVV"`
	Name       string `json:"name"`
}

type Text struct {
	Data string `json:"data"`
}

type Binary struct {
	Data string `json:"data"`
}

type MetaData map[string]string

type Record struct {
	Id          int64        `json:"id"`
	RecordType  RecordType   `json:"record_type"`
	MetaData    MetaData     `json:"meta_data,omitempty"`
	Credentials *Credentials `json:"credentials,omitempty"`
	BankCard    *BankCard    `json:"bank_card,omitempty"`
	Text        *Text        `json:"text,omitempty"`
	Binary      *Binary      `json:"binary,omitempty"`
	Deleted     bool         `json:"deleted"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
