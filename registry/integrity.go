package registry

import (
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
)

type Integrity struct {
	SHASum    string `json:"shasum"`
	Integrity string `json:"integrity"`
	Base64    string `json:"base_64"`
	Length    int64  `json:"length"`
}

func MakeIntegrity(reader io.Reader) (*Integrity, error) {
	sha1Hash := sha1.New()
	sha512Hash := sha512.New()
	b64Result := &strings.Builder{}
	base64Encoder := base64.NewEncoder(base64.StdEncoding, b64Result)
	size, err := io.Copy(io.MultiWriter(sha1Hash, sha512Hash, base64Encoder), reader)
	_ = base64Encoder.Close()
	if err != nil {
		return nil, err
	}
	return &Integrity{
		SHASum:    hex.EncodeToString(sha1Hash.Sum(nil)),
		Integrity: "sha512-" + base64.StdEncoding.EncodeToString(sha512Hash.Sum(nil)),
		Base64:    b64Result.String(),
		Length:    size,
	}, nil
}
