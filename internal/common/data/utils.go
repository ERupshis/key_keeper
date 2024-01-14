package data

func ConvertStringToRecordType(str string) RecordType {
	switch str {
	case StrCredentials:
		return TypeCredentials
	case StrBankCard:
		return TypeBankCard
	case StrText:
		return TypeText
	case StrBinary:
		return TypeBinary
	case StrAny:
		return TypeAny
	default:
		return TypeUndefined
	}
}

func DeepCopyRecord(record *Record) *Record {
	var res Record
	res.Id = record.Id
	res.RecordType = record.RecordType

	res.MetaData = make(MetaData)
	for key, val := range record.MetaData {
		res.MetaData[key] = val
	}

	if record.Credentials != nil {
		tmpCredentials := *record.Credentials
		res.Credentials = &tmpCredentials
	}

	if record.BankCard != nil {
		tmpBankCard := *record.BankCard
		res.BankCard = &tmpBankCard
	}

	if record.Text != nil {
		tmpText := *record.Text
		res.Text = &tmpText
	}

	if record.Binary != nil {
		tmpBinary := *record.Binary
		res.Binary = &tmpBinary
	}

	return &res
}
