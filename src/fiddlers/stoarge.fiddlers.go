package fiddlers

import (
	"net/url"

	"github.com/patos-ufscar/quack-week/common"
)

func GetFullObjStorageUrl(objPath string) (string, error) {
	prefix := "https://"
	if !common.S3_SECURE {
		prefix = "http://"
	}
	return url.JoinPath(prefix+common.S3_ENDPOINT, common.S3_BUCKET, objPath)
}
