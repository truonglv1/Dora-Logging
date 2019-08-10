// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package djson

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

func easyjson50c2aa5cDecodeGithubComDoraLogsInternalDjson(in *jlexer.Lexer, out *OsGroup) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "os_code":
			out.OsCode = int(in.Int())
		case "os_ver":
			out.OsVer = string(in.String())
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
func easyjson50c2aa5cEncodeGithubComDoraLogsInternalDjson(out *jwriter.Writer, in OsGroup) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"os_code\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.OsCode))
	}
	{
		const prefix string = ",\"os_ver\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.OsVer))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v OsGroup) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson50c2aa5cEncodeGithubComDoraLogsInternalDjson(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OsGroup) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson50c2aa5cEncodeGithubComDoraLogsInternalDjson(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OsGroup) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson50c2aa5cDecodeGithubComDoraLogsInternalDjson(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OsGroup) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson50c2aa5cDecodeGithubComDoraLogsInternalDjson(l, v)
}
func easyjson50c2aa5cDecodeGithubComDoraLogsInternalDjson1(in *jlexer.Lexer, out *ActionLog) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ip":
			out.Ip = string(in.String())
		case "os_group":
			(out.OsGroup).UnmarshalEasyJSON(in)
		case "session_id":
			out.SessionId = string(in.String())
		case "category_id":
			out.CategoryId = string(in.String())
		case "event_id":
			out.EventId = int(in.Int())
		case "time_create":
			out.TimeCreate = int64(in.Int64())
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
func easyjson50c2aa5cEncodeGithubComDoraLogsInternalDjson1(out *jwriter.Writer, in ActionLog) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Ip))
	}
	{
		const prefix string = ",\"os_group\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.OsGroup).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"session_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.SessionId))
	}
	{
		const prefix string = ",\"category_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.CategoryId))
	}
	{
		const prefix string = ",\"event_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.EventId))
	}
	{
		const prefix string = ",\"time_create\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.TimeCreate))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ActionLog) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson50c2aa5cEncodeGithubComDoraLogsInternalDjson1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ActionLog) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson50c2aa5cEncodeGithubComDoraLogsInternalDjson1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ActionLog) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson50c2aa5cDecodeGithubComDoraLogsInternalDjson1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ActionLog) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson50c2aa5cDecodeGithubComDoraLogsInternalDjson1(l, v)
}
