package storage

import (
	"net/url"
	"path"

	"github.com/patos-ufscar/quack-week/common"
)

type storagePrefix string

const (
	EVENT_BANNERS storagePrefix = "public/event-banners"
	USER_AVATARS  storagePrefix = "public/user-avatars"
)

func GetFullObjUrl(objPath string) (string, error) {
	prefix := "https://"
	if !common.S3_SECURE {
		prefix = "http://"
	}
	return url.JoinPath(prefix+common.S3_ENDPOINT, common.S3_BUCKET, objPath)
}

func GetPublicPath(p storagePrefix, filename string) string {
	return path.Join("public", string(p), filename)
}

func GetPrivatePath(p storagePrefix, filename string) string {
	return path.Join("private", string(p), filename)
}
