package data

import (
	"fmt"
	"strings"
	"time"
)

const (
	TypeUndefined = RecordType(0)

	TypeCredentials = RecordType(1)
	TypeBankCard    = RecordType(2)
	TypeText        = RecordType(3)
	TypeBinary      = RecordType(4)
	TypeAny         = RecordType(5)
)

const (
	StrUndefined = "undefined"

	StrCredentials = "creds"
	StrBankCard    = "card"
	StrText        = "text"
	StrBinary      = "bin"
	StrAny         = "any"
)

const (
	DataInvalid = "INVALID"
)

type RecordType = int32

//go:generate easyjson -all data.go
type Credential struct {
	Login    string `json:"login"`
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
	Name            string `json:"name"`
	SecuredFileName string `json:"file"`
}

type MetaData map[string]string

type Record struct {
	ID          int64       `json:"id"`
	RecordType  RecordType  `json:"record_type"`
	MetaData    MetaData    `json:"meta_data,omitempty"`
	Credentials *Credential `json:"credentials,omitempty"`
	BankCard    *BankCard   `json:"bank_card,omitempty"`
	Text        *Text       `json:"text,omitempty"`
	Binary      *Binary     `json:"binary,omitempty"`
	Deleted     bool        `json:"deleted"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (r Record) String() string {
	formatBuilder := strings.Builder{}
	formatBuilder.WriteString("{ID: %d,")

	switch r.RecordType {
	case TypeCredentials:
		formatBuilder.WriteString(" Credential: %+v,")
	case TypeBankCard:
		formatBuilder.WriteString(" BankCard: %+v,")
	case TypeText:
		formatBuilder.WriteString(" Text: %+v,")
	case TypeBinary:
		formatBuilder.WriteString(" Binary: %+v,")
	default:
	}
	formatBuilder.WriteString(" MetaData: %s}")

	return fmt.Sprintf(
		formatBuilder.String(),
		r.ID,
		getRecordValue(&r),
		r.MetaData,
	)
}

func (r Record) TabString() string {
	formatBuilder := strings.Builder{}
	formatBuilder.WriteString("\tID: %d")

	switch r.RecordType {
	case TypeCredentials:
		formatBuilder.WriteString("\tCredential: %+v")
	case TypeBankCard:
		formatBuilder.WriteString("\tBankCard: %+v")
	case TypeText:
		formatBuilder.WriteString("\tText: %+v")
	case TypeBinary:
		formatBuilder.WriteString("\tBinary: %+v")
	default:
	}

	formatBuilder.WriteString("\tMetaData: %s")

	return fmt.Sprintf(
		formatBuilder.String(),
		r.ID,
		getRecordValue(&r),
		r.MetaData,
	)
}
