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
}

func MakeIntegrity(reader io.Reader) (*Integrity, error) {
	h1 := sha1.New()
	h2 := sha512.New()
	b64 := &strings.Builder{}
	b1 := base64.NewEncoder(base64.StdEncoding, b64)
	_, err := io.Copy(io.MultiWriter(h1, h2, b1), reader)
	_ = b1.Close()
	if err != nil {
		return nil, err
	}
	return &Integrity{
		SHASum:    hex.EncodeToString(h1.Sum(nil)),
		Integrity: "sha512-" + base64.StdEncoding.EncodeToString(h2.Sum(nil)),
		Base64:    b64.String(),
	}, nil
}
