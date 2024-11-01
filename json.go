package pgtype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type JSON []byte

func (dst *JSON) Set(src interface{}) error {
	if src == nil {
		*dst = nil
		return nil
	}

	if value, ok := src.(interface{ Get() interface{} }); ok {
		value2 := value.Get()
		if value2 != value {
			return dst.Set(value2)
		}
	}

	switch value := src.(type) {
	case string:
		*dst = []byte(value)
	case *string:
		if value == nil {
			*dst = nil
		} else {
			*dst = []byte(*value)
		}
	case []byte:
		if value == nil {
			*dst = nil
		} else {
			*dst = value
		}
	// Encode* methods are defined on *JSON. If JSON is passed directly then the
	// struct itself would be encoded instead of Bytes. This is clearly a footgun
	// so detect and return an error. See https://github.com/jackc/pgx/issues/350.
	case JSON:
		return errors.New("use pointer to pgtype.JSON instead of value")
	// Same as above but for JSONB (because they share implementation)
	case JSONB:
		return errors.New("use pointer to pgtype.JSONB instead of value")

	default:
		buf, err := json.Marshal(value)
		if err != nil {
			return err
		}
		*dst = buf
	}

	return nil
}

func (dst JSON) Get() interface{} {
	if dst == nil {
		return nil
	}

	var i interface{}
	err := json.Unmarshal(dst, &i)
	if err != nil {
		return dst
	}
	return i
}

func (src *JSON) AssignTo(dst interface{}) error {
	switch v := dst.(type) {
	case *string:
		if src != nil && *src != nil {
			*v = string(*src)
		} else {
			return fmt.Errorf("cannot assign non-present status to %T", dst)
		}
	case **string:
		if src != nil {
			s := string(*src)
			*v = &s
			return nil
		} else {
			*v = nil
			return nil
		}
	case *[]byte:
		if src != nil {
			*v = nil
		} else {
			buf := make([]byte, len(*src))
			copy(buf, *src)
			*v = buf
		}
	default:
		data := *src
		if data == nil || src == nil || *src == nil {
			data = []byte("null")
		}

		p := reflect.ValueOf(dst).Elem()
		p.Set(reflect.Zero(p.Type()))

		return json.Unmarshal(data, dst)
	}

	return nil
}

func (JSON) PreferredResultFormat() int16 {
	return TextFormatCode
}

func (dst *JSON) DecodeText(_ *ConnInfo, src []byte) error {
	if src == nil {
		*dst = nil
		return nil
	}

	*dst = src
	return nil
}

func (dst *JSON) DecodeBinary(ci *ConnInfo, src []byte) error {
	return dst.DecodeText(ci, src)
}

func (JSON) PreferredParamFormat() int16 {
	return TextFormatCode
}

func (src JSON) EncodeText(_ *ConnInfo, buf []byte) ([]byte, error) {
	if src == nil {
		return nil, nil
	}

	return append(buf, src...), nil
}

func (src JSON) EncodeBinary(ci *ConnInfo, buf []byte) ([]byte, error) {
	return src.EncodeText(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *JSON) Scan(src interface{}) error {
	if src == nil {
		*dst = nil
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.DecodeText(nil, []byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return dst.DecodeText(nil, srcCopy)
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src JSON) Value() (driver.Value, error) {
	if src == nil {
		return nil, nil
	}

	return src, nil
}

func (src JSON) MarshalJSON() ([]byte, error) {
	if src == nil {
		return []byte("null"), nil
	}

	if len(src) == 0 {
		return []byte("[]"), nil
	}

	return src, nil
}

func (dst *JSON) UnmarshalJSON(b []byte) error {
	if b == nil || string(b) == "null" {
		*dst = nil
	} else {
		bCopy := make([]byte, len(b))
		copy(bCopy, b)
		*dst = bCopy
	}
	return nil
}
