package pgtype

import (
	"database/sql/driver"
	"fmt"
)

type JSONB JSON

func (dst *JSONB) Set(src interface{}) error {
	return (*JSON)(dst).Set(src)
}

func (dst JSONB) Get() interface{} {
	return (JSON)(dst).Get()
}

func (src *JSONB) AssignTo(dst interface{}) error {
	return (*JSON)(src).AssignTo(dst)
}

func (JSONB) PreferredResultFormat() int16 {
	return TextFormatCode
}

func (dst *JSONB) DecodeText(ci *ConnInfo, src []byte) error {
	return (*JSON)(dst).DecodeText(ci, src)
}

func (dst *JSONB) DecodeBinary(ci *ConnInfo, src []byte) error {
	if src == nil {
		*dst = nil
		return nil
	}

	if len(src) == 0 {
		return fmt.Errorf("jsonb too short")
	}

	if src[0] != 1 {
		return fmt.Errorf("unknown jsonb version number %d", src[0])
	}

	*dst = src[1:]
	return nil

}

func (JSONB) PreferredParamFormat() int16 {
	return TextFormatCode
}

func (src JSONB) EncodeText(ci *ConnInfo, buf []byte) ([]byte, error) {
	return (JSON)(src).EncodeText(ci, buf)
}

func (src JSONB) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	if src == nil {
		return nil, nil
	}

	buf = append(buf, 1)
	return append(buf, src...), nil
}

// Scan implements the database/sql Scanner interface.
func (dst *JSONB) Scan(src interface{}) error {
	return (*JSON)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src JSONB) Value() (driver.Value, error) {
	return (JSON)(src).Value()
}

func (src JSONB) MarshalJSON() ([]byte, error) {
	return (JSON)(src).MarshalJSON()
}

func (dst *JSONB) UnmarshalJSON(b []byte) error {
	return (*JSON)(dst).UnmarshalJSON(b)
}
