package utils

import "os"

func GetAppURL() string {
	url := os.Getenv("APP_URL")
	if url == "" {
		return "https://sijil.io"
	}
	return url
}
