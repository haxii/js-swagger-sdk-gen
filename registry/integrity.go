package registry

import (
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"
)

type Integrity struct {
	SHASum    string `json:"shasum"`
	Integrity string `json:"integrity"`
}

func MakeIntegrity(reader io.Reader) (*Integrity, error) {
	h1 := sha1.New()
	h2 := sha512.New()
	_, err := io.Copy(io.MultiWriter(h1, h2), reader)
	if err != nil {
		return nil, err
	}
	return &Integrity{
		SHASum:    hex.EncodeToString(h1.Sum(nil)),
		Integrity: "sha512-" + base64.StdEncoding.EncodeToString(h2.Sum(nil)),
	}, nil
}
