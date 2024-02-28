package registry

import "encoding/json"

type VersionDist struct {
	Integrity string `json:"integrity"` // sha512 in b64
	SHASum    string `json:"shasum"`    // sha1 in hex
	Tarball   string `json:"tarball"`   // tar's location url in registry
} // should in "dist" key of PublishInfo.Versions

type PublishAttachments struct {
	ContentType string `json:"content_type"` // application/octet-stream
	Data        string `json:"data"`         // b64 encoded tar
	Length      int64  `json:"length"`       // tar size
}

type rawMessage map[string]json.RawMessage

type PublishInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DistTags    struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Versions    map[string]rawMessage         `json:"versions"`     // key is version
	Attachments map[string]PublishAttachments `json:"_attachments"` // key is tar name
}
