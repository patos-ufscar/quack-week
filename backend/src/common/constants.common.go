package common

import (
	"time"
)

const (
	// TIMESTAMP_STR_FORMAT string = "yyyy-mm-ddThh:mm:ssZhh:mm"
	// TIMESTAMP_STR_FORMAT string = "2006-01-02T15:04:05-07:00"
	TIMESTAMP_STR_FORMAT        string = time.RFC3339
	DEFAULT_TIMEZONE            string = "GMT-3"
	GIN_CTX_JWT_CLAIM_KEY_NAME  string = "jwtClaims"
	JWT_TIMEOUT_SECS            int    = 30 * 60
	OTP_LEN                     int    = 128
	ORG_INVITE_TIMEOUT_DAYS     int    = 15
	PASSWORD_RESET_TIMEOUT_DAYS int    = 1
	MAX_REQUEST_SIZE            int64  = 5 * 1024 * 1024 // 5MB default
)

var (
	PROJECT_NAME                           string = GetEnvVarDefault("PROJECT_NAME", "patos-app")
	NOREPLY_EMAIL                          string = GetEnvVarDefault("NOREPLY_EMAIL", "no-reply@example.com")
	APP_HOST_URL                           string = GetEnvVarDefault("APP_HOST_URL", "http://127.0.0.1:8080/")
	API_HOST_URL                           string = GetEnvVarDefault("API_HOST_URL", "http://127.0.0.1:8080/")
	JWT_COOKIE_NAME                        string = PROJECT_NAME + "_jwt"
	PASSWORD_RESET_TIMEOUT_JWT_COOKIE_NAME string = PROJECT_NAME + "_pwreset_jwt"
	S3_ENDPOINT                            string = GetEnvVarDefault("S3_ENDPOINT", "https://br-se1.magaluobjects.com")
	S3_REGION                              string = GetEnvVarDefault("S3_REGION", "br-se1")
	S3_BUCKET                              string = GetEnvVarDefault("S3_BUCKET", PROJECT_NAME+"-gopherbase")
)
