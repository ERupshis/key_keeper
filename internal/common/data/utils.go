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

	tmpBankCard := *record.BankCard
	res.BankCard = &tmpBankCard

	tmpCredentials := *record.Credentials
	res.Credentials = &tmpCredentials

	tmpText := *record.Text
	res.Text = &tmpText

	tmpBinary := *record.Binary
	res.Binary = &tmpBinary

	return &res
}
