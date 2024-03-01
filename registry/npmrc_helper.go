package registry

import (
	"errors"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
)

type DefaultNpmRc struct {
	Path, URL, Token string
}

func GetDefaultNpmRC() (*DefaultNpmRc, error) {
	rcs, err := GetNpmRC()
	if err != nil {
		return nil, err
	}
	if len(rcs) == 0 {
		return nil, errors.New(".npmrc not found")
	}
	d := &DefaultNpmRc{}
	for _, rc := range rcs {
		if len(rc.DefaultURL) > 0 {
			d.Path = rc.FilePath
			d.URL = rc.DefaultURL
			break
		}
	}
	if len(d.URL) == 0 {
		return nil, errors.New("default registry not found in .npmrc")
	}
	d.Token = FindToken(d.URL, rcs)
	return d, nil
}

// GetNpmRC find the .npmrc file in system and parse it
// will find from `pwd`, then $HOME folder
func GetNpmRC() ([]*NpmRC, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	rcs := make([]*NpmRC, 0)
	for _, loc := range []string{wd, home} {
		p := filepath.Join(loc, ".npmrc")
		if rc, parseErr := ParseNpmRc(p); err != nil && !errors.Is(parseErr, os.ErrNotExist) {
			return nil, parseErr
		} else if rc != nil {
			rcs = append(rcs, rc)
		}
	}

	return rcs, nil
}

type NpmRCToken struct {
	URL, Token string
}

type NpmRC struct {
	FilePath   string
	DefaultURL string
	Tokens     []NpmRCToken
}

// ParseNpmRc parse main reg url and token
// npmrc declares in ini format https://docs.npmjs.com/cli/v10/configuring-npm/npmrc
// we only get default registry and all the tokens auth
// scoped settings and other auth method is not yet implemented
func ParseNpmRc(p string) (*NpmRC, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	iniFile, err := ini.Load(f)
	if err != nil {
		return nil, err
	}
	rc := &NpmRC{FilePath: p, Tokens: make([]NpmRCToken, 0)}
	for _, section := range iniFile.Sections() {
		for _, key := range section.Keys() {
			if key.Name() == "registry" {
				rc.DefaultURL = NormRegURL(key.Value())
			} else if strings.HasPrefix(key.Value(), "_authToken") {
				type authToken struct {
					Token string `ini:"_authToken"`
				}
				t := authToken{}
				if err = ini.MapTo(&t, []byte(key.Value())); err != nil {
					continue
				}
				rc.Tokens = append(rc.Tokens, NpmRCToken{
					URL:   NormRegURL(key.Name()),
					Token: t.Token,
				})
			}
		}
	}
	return rc, nil
}

func FindToken(regURL string, rcs []*NpmRC) string {
	regURL = NormRegURL(regURL)
	for _, rc := range rcs {
		for _, token := range rc.Tokens {
			if strings.HasSuffix(regURL, token.URL) {
				return token.Token
			}
		}
	}
	return ""
}

// NormRegURL make sure reg url has a `/` suffix
func NormRegURL(u string) string {
	if strings.HasSuffix(u, "/") {
		return u
	}
	return u + "/"
}
