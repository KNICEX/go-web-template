package hashid

import (
	"errors"
	"github.com/sqids/sqids-go"
)

const (
	ShareID = iota
	UserID
)

var s *sqids.Sqids

func init() {
	var err error
	s, err = sqids.New()
	if err != nil {
		panic(err)
	}
}

var (
	ErrorTypeNotMatch = errors.New("mismatched ID type")
)

func HashEncode(v []uint64) (string, error) {
	id, err := s.Encode(v)
	if err != nil {
		return "", err
	}
	return id, nil
}

func HashDecode(raw string) []uint64 {
	return s.Decode(raw)
}

func HashId(id uint64, t uint64) string {
	v, _ := HashEncode([]uint64{id, t})
	return v
}

func DecodeId(id string, t int) (uint64, error) {
	v := HashDecode(id)
	if len(v) != 2 || v[1] != uint64(t) {
		return 0, ErrorTypeNotMatch
	}
	return v[0], nil
}
