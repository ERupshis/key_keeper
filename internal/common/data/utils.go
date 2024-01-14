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
	switch record.RecordType {
	case TypeCredentials:
		if record.Credentials != nil {
			return *record.Credentials
		}
	case TypeBankCard:
		if record.BankCard != nil {
			return *record.BankCard
		}
	case TypeText:
		if record.Text != nil {
			return *record.Text
		}
	case TypeBinary:
		if record.Binary != nil {
			return *record.Binary
		}
	}

	return DataInvalid
}

func DeepCopyRecord(record *Record) *Record {
	var res Record
	res.ID = record.ID
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
