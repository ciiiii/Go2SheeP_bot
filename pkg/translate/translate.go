package translate

import (
	"net/http"
	"regexp"
	"strings"
	"time"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"errors"
)

type Translator struct {
	appId  string
	key    string
	client *http.Client
}

type respBody struct {
	Result []struct{
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
	From string `json:"from"`
	To string `json:"to"`
}

func NewTranslator(appId, key string) Translator {
	return Translator{
		appId:  appId,
		key:    key,
		client: &http.Client{},
	}
}

func (t *Translator) Translate(from, to, source string) (string, error) {
	reg, _ := regexp.Compile(`[\s~!@<>#$%\+\^\&\*\(\)\./\\]+`)
	source = reg.ReplaceAllString(strings.ToLower(source), ",")
	salt := time.Now().Unix()
	sign := fmt.Sprintf("%s%s%d%s", t.appId, source, salt, t.key)
	sign = fmt.Sprintf("%x", md5.Sum([]byte(sign)))
	url := fmt.Sprintf("http://api.fanyi.baidu.com/api/trans/vip/translate?q=%s&appid=%s&salt=%d&from=%s&to=%s&sign=%s&_=%d",
		source,
		t.appId,
		salt,
		from,
		to,
		sign,
		time.Now().Unix(),
	)
	resp, err := t.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result respBody
	json.Unmarshal(body, &result)
	if len(result.Result) >= 1 {
		return result.Result[0].Dst, nil
	}
	return "", errors.New("translate failed")
}
