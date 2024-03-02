package registry

import (
	"encoding/json"
	"fmt"
	"github.com/haxii/js-swagger-sdk-gen/model"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
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
	regURL = regURL.JoinPath("/")
	if !strings.HasPrefix(token, "Bearer ") {
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
	model.PackageJSON
	raw rawMessage
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

func (api *API) Publish(tarball io.Reader, packageJSON json.RawMessage) error {
	// parse package info
	pkg, err := ParsePackageInfo(packageJSON)
	if err != nil {
		return err
	}
	// parse integrity
	integrity, err := MakeIntegrity(tarball)
	if err != nil {
		return err
	}
	pkgTarURL := pkg.URL(api.url)

	// build dist in version tag, merge it into info's raw
	verDist := &VersionDist{
		Integrity: integrity.Integrity,
		SHASum:    integrity.SHASum,
		Tarball:   pkgTarURL,
	}
	if pkg.raw["dist"], err = json.Marshal(verDist); err != nil {
		return err
	}
	pkg.raw["_id"] = []byte(strconv.Quote(pkg.NpmID()))

	// build publish info
	info := PublishInfo{
		Name:        pkg.Name,
		Description: pkg.Description,
		Versions:    map[string]rawMessage{pkg.Version: pkg.raw},
		Attachments: map[string]PublishAttachments{pkgTarURL: {
			ContentType: "application/octet-stream",
			Data:        integrity.Base64,
			Length:      integrity.Length,
		}},
	}
	info.DistTags.Latest = pkg.Version

	// make a pipe and write package info into it
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(false)
		if encodeErr := enc.Encode(&info); encodeErr != nil {
			fmt.Println("fail to encode publish info", err)
		}
	}()

	// publish info to npm registry
	resp, err := api.client.Do(&http.Request{
		Method: http.MethodPut,
		URL:    api.url.JoinPath(pkg.Name),
		Header: http.Header{
			"content-type":  []string{"application/json"},
			"authorization": []string{api.token},
			"accept":        []string{"*/*"},
		},
		Body: r,
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
