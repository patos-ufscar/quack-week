package storage

import (
	"net/url"
	"path"

	"github.com/patos-ufscar/quack-week/common"
)

type storageDir string

const (
	EVENT_BANNERS storageDir = "event-banners"
	USER_AVATARS  storageDir = "user-avatars"
)

func GetFullObjUrl(objPath string) (string, error) {
	return url.JoinPath(common.S3_ENDPOINT, common.S3_BUCKET, objPath)
}

func GetPublicPath(p storageDir, filename string) string {
	return path.Join("public", string(p), filename)
}

func GetPrivatePath(p storageDir, filename string) string {
	return path.Join("private", string(p), filename)
}
