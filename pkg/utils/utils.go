package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/lithammer/shortuuid"
	"github.com/pkg/errors"
)

func UUID() string {
	id := shortuuid.New()
	return id
}

func GetFileUrl(token, fileId string) (string, error) {
	u := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", token, fileId)
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var r struct {
		Ok     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return "", err
	}
	if r.Ok {
		return fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", token, r.Result.FilePath), nil
	} else {
		return "", errors.New("get file url failed")
	}
}

func DownloadAsTmp(fileUrl string) (string, error) {
	splitStr := strings.Split(fileUrl, ".")
	suffix := splitStr[len(splitStr)-1]
	tmpPath := fmt.Sprintf("/tmp/%s.%s", UUID(), suffix)
	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.Create(tmpPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return tmpPath, err
}

func DeleteFile(filePath string) {
	os.Remove(filePath)
}

func ParseImageKey(url string) string {
	splitStr := strings.Split(url, "/")
	return splitStr[len(splitStr)-1]
}

func CleanInvalidStr(o string) string {
	printOnly := func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}
	return strings.Map(printOnly, o)
}
