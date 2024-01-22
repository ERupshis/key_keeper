// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package data

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData(in *jlexer.Lexer, out *Text) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "data":
			out.Data = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData(out *jwriter.Writer, in Text) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix[1:])
		out.String(string(in.Data))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Text) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Text) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Text) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Text) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData(l, v)
}
func easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData1(in *jlexer.Lexer, out *Record) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "record_type":
			out.RecordType = int32(in.Int32())
		case "meta_data":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.MetaData = make(MetaData)
				} else {
					out.MetaData = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 string
					v1 = string(in.String())
					(out.MetaData)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "credentials":
			if in.IsNull() {
				in.Skip()
				out.Credentials = nil
			} else {
				if out.Credentials == nil {
					out.Credentials = new(Credentials)
				}
				(*out.Credentials).UnmarshalEasyJSON(in)
			}
		case "bank_card":
			if in.IsNull() {
				in.Skip()
				out.BankCard = nil
			} else {
				if out.BankCard == nil {
					out.BankCard = new(BankCard)
				}
				(*out.BankCard).UnmarshalEasyJSON(in)
			}
		case "text":
			if in.IsNull() {
				in.Skip()
				out.Text = nil
			} else {
				if out.Text == nil {
					out.Text = new(Text)
				}
				(*out.Text).UnmarshalEasyJSON(in)
			}
		case "binary":
			if in.IsNull() {
				in.Skip()
				out.Binary = nil
			} else {
				if out.Binary == nil {
					out.Binary = new(Binary)
				}
				(*out.Binary).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData1(out *jwriter.Writer, in Record) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"record_type\":"
		out.RawString(prefix)
		out.Int32(int32(in.RecordType))
	}
	if len(in.MetaData) != 0 {
		const prefix string = ",\"meta_data\":"
		out.RawString(prefix)
		{
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.MetaData {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				out.String(string(v2Value))
			}
			out.RawByte('}')
		}
	}
	if in.Credentials != nil {
		const prefix string = ",\"credentials\":"
		out.RawString(prefix)
		(*in.Credentials).MarshalEasyJSON(out)
	}
	if in.BankCard != nil {
		const prefix string = ",\"bank_card\":"
		out.RawString(prefix)
		(*in.BankCard).MarshalEasyJSON(out)
	}
	if in.Text != nil {
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		(*in.Text).MarshalEasyJSON(out)
	}
	if in.Binary != nil {
		const prefix string = ",\"binary\":"
		out.RawString(prefix)
		(*in.Binary).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Record) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Record) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Record) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Record) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData1(l, v)
}
func easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData2(in *jlexer.Lexer, out *Credentials) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "key":
			out.Login = string(in.String())
		case "password":
			out.Password = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData2(out *jwriter.Writer, in Credentials) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"key\":"
		out.RawString(prefix[1:])
		out.String(string(in.Login))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Credentials) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Credentials) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Credentials) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Credentials) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData2(l, v)
}
func easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData3(in *jlexer.Lexer, out *Binary) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "data":
			out.Data = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData3(out *jwriter.Writer, in Binary) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix[1:])
		out.String(string(in.Data))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Binary) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Binary) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Binary) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Binary) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData3(l, v)
}
func easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData4(in *jlexer.Lexer, out *BankCard) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "number":
			out.Number = string(in.String())
		case "expiration":
			out.Expiration = string(in.String())
		case "CVV":
			out.CVV = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData4(out *jwriter.Writer, in BankCard) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"number\":"
		out.RawString(prefix[1:])
		out.String(string(in.Number))
	}
	{
		const prefix string = ",\"expiration\":"
		out.RawString(prefix)
		out.String(string(in.Expiration))
	}
	{
		const prefix string = ",\"CVV\":"
		out.RawString(prefix)
		out.String(string(in.CVV))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BankCard) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BankCard) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson794297d0EncodeGithubComErupshisKeyKeeperInternalCommonData4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BankCard) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BankCard) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson794297d0DecodeGithubComErupshisKeyKeeperInternalCommonData4(l, v)
}
