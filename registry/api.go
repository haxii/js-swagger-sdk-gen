package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type API struct {
	url    *url.URL
	token  string // registry url & auth token
	client *http.Client
}

func NewAPI(u string, token string) (*API, error) {
	regURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(token, "Bearer ") {
		token = fmt.Sprintf("Bearer %s", token)
	}
	return &API{
		url: regURL, token: token,
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				DialContext: (&net.Dialer{Timeout: time.Second * 10}).DialContext,
			},
			Timeout: 5 * time.Minute,
		},
	}, nil
}

type PackageInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	// raw message of package info
	raw map[string]json.RawMessage
}

func ParsePackageInfo(packageJSON json.RawMessage) (*PackageInfo, error) {
	pkgInfo := &PackageInfo{}
	if err := json.Unmarshal(packageJSON, &pkgInfo); err != nil {
		return nil, err
	}
	pkgInfo.raw = make(map[string]json.RawMessage)
	if err := json.Unmarshal(packageJSON, &pkgInfo.raw); err != nil {
		return nil, err
	}
	return pkgInfo, nil
}

func makePublishBody() (io.ReadCloser, error) {
	// TODO
	return nil, nil
}

func (api *API) Publish(tarball io.Reader, packageJSON json.RawMessage) error {
	pkg, err := ParsePackageInfo(packageJSON)
	if err != nil {
		return err
	}
	resp, err := api.client.Do(&http.Request{
		Method: http.MethodPut,
		URL:    api.url.JoinPath(pkg.Name),
		Header: http.Header{
			"Content-Type":  []string{"application/json"},
			"Authorization": []string{api.token},
		},
		Body: nil,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("invalid status %d on uploading publish: %s", resp.StatusCode, body)
}
