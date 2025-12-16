package goda

import (
	"bytes"
	"database/sql/driver"
	"strconv"
)

func (y *YearMonth) UnmarshalJSON(i []byte) error {
	return unmarshalJsonImpl(y, i)
}

func (y YearMonth) MarshalJSON() ([]byte, error) {
	return marshalJsonImpl(y)
}

func (y YearMonth) MarshalText() (text []byte, err error) {
	return marshalTextImpl(y)
}

func (y *YearMonth) UnmarshalText(text []byte) (e error) {
	if len(text) == 0 {
		return nil
	}
	var i = bytes.IndexByte(text, '-')
	if i == -1 {
		return parseFailedError(text)
	}
	var year int64
	var month int64
	year, e = strconv.ParseInt(string(text[:i]), 10, 64)
	if e != nil {
		return
	}
	month, e = strconv.ParseInt(string(text[i+1:]), 10, 64)
	if e != nil {
		return
	}
	*y, e = YearMonthOf(Year(year), Month(month))
	return
}

func (y YearMonth) AppendText(b []byte) ([]byte, error) {
	if y.IsZero() {
		return b, nil
	}
	b, _ = y.Year().AppendText(b)
	b = append(b, '-')
	if y.Month() < December {
		b = append(b, '0')
	}
	b = strconv.AppendInt(b, int64(y.Month()), 10)
	return b, nil
}

func (y YearMonth) String() string {
	return stringImpl(y)
}

func (y *YearMonth) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*y = YearMonth{}
		return nil
	case string:
		return y.UnmarshalText([]byte(v))
	case []byte:
		return y.UnmarshalText(v)
	default:
		return sqlScannerDefaultBranch(v)
	}
}

func (y YearMonth) Value() (driver.Value, error) {
	if y.IsZero() {
		return nil, nil
	}
	return y.String(), nil
}
