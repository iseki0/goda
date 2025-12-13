package goda

import "strconv"

// String returns the string representation of this year.
// Years 0-9999 are formatted as 4 digits with leading zeros (e.g., "0001", "2024").
// Years outside this range are formatted without padding.
func (y Year) String() string {
	return stringImpl(y)
}

// AppendText implements the encoding.TextAppender interface.
// It appends the year representation to b and returns the extended buffer.
func (y Year) AppendText(b []byte) ([]byte, error) {
	if y >= 0 && y <= 9999 {
		return append(b, '0'+byte(y/1000), '0'+byte((y/100)%10), '0'+byte((y/10)%10), '0'+byte(y%10)), nil
	} else if y < 0 && y >= -9999 {
		return append(b, '-', '0'+byte((-y)/1000), '0'+byte(((-y)/100)%10), '0'+byte(((-y)/10)%10), '0'+byte((-y)%10)), nil
	}
	b = append(b, strconv.FormatInt(y.Int64(), 10)...)
	return b, nil
}
