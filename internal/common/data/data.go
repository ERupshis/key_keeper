package data

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

type MetaData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Record struct {
	Id          int64        `json:"id"`
	MetaData    []MetaData   `json:"meta_data,omitempty"`
	Credentials *Credentials `json:"credentials,omitempty"`
	BankCard    *BankCard    `json:"bank_card,omitempty"`
	Text        *Text        `json:"text,omitempty"`
	Binary      *Binary      `json:"binary,omitempty"`
}
