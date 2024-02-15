package models

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

func ConvertRecordTypeToString(recordType RecordType) string {
	switch recordType {
	case TypeCredentials:
		return StrCredentials
	case TypeBankCard:
		return StrBankCard
	case TypeText:
		return StrText
	case TypeBinary:
		return StrBinary
	case TypeAny:
		return StrAny
	default:
		return StrUndefined
	}
}

func getRecordValue(record *Record) interface{} {
	switch record.Data.RecordType {
	case TypeCredentials:
		if record.Data.Credentials != nil {
			return *record.Data.Credentials
		}
	case TypeBankCard:
		if record.Data.BankCard != nil {
			return *record.Data.BankCard
		}
	case TypeText:
		if record.Data.Text != nil {
			return *record.Data.Text
		}
	case TypeBinary:
		if record.Data.Binary != nil {
			return *record.Data.Binary
		}
	}

	return Invalid
}

func DeepCopyRecord(record *Record) *Record {
	var res Record
	res.ID = record.ID
	res.Data.RecordType = record.Data.RecordType

	res.Data.MetaData = make(MetaData)
	for key, val := range record.Data.MetaData {
		res.Data.MetaData[key] = val
	}

	if record.Data.Credentials != nil {
		tmpCredentials := *record.Data.Credentials
		res.Data.Credentials = &tmpCredentials
	}

	if record.Data.BankCard != nil {
		tmpBankCard := *record.Data.BankCard
		res.Data.BankCard = &tmpBankCard
	}

	if record.Data.Text != nil {
		tmpText := *record.Data.Text
		res.Data.Text = &tmpText
	}

	if record.Data.Binary != nil {
		tmpBinary := *record.Data.Binary
		res.Data.Binary = &tmpBinary
	}

	return &res
}
