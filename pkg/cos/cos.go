package cos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"github.com/ciiiii/Go2SheeP_bot/pkg/utils"
	"github.com/tencentyun/cos-go-sdk-v5"
	"strings"
)

type Client struct {
	client *cos.Client
	prefix string
}

func NewCos(bucket, region, secretId, secretKey string) *Client {
	prefix := fmt.Sprintf("http://%s.cos.%s.myqcloud.com", bucket, region)
	u, _ := url.Parse(prefix)
	b := &cos.BaseURL{BucketURL: u}
	c := Client{cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	}), prefix}
	return &c
}

func (c Client) List() ([]string, error) {
	opt := cos.BucketGetOptions{MaxKeys: 24}

	r, _, err := c.client.Bucket.Get(context.Background(), &opt)
	if err != nil {
		return nil, err
	}
	var fileList []string
	for _, o := range r.Contents {
		if o.Size != 0 {
			f := c.client.BaseURL.BucketURL.String() + "/" + o.Key
			fileList = append(fileList, f)
		}
	}
	return fileList, nil
}

func (c Client) Upload(path string) (string, error) {
	splitPath := strings.Split(path, ".")
	suffix := splitPath[len(splitPath)-1]
	key := fmt.Sprintf("%s.%s", utils.UUID(), suffix)
	_, err := c.client.Object.PutFromFile(context.Background(), key, path, nil)
	u := c.client.BaseURL.BucketURL.String() + "/" + key
	return u, err
}

func (c Client) Delete(key string) error {
	_, err := c.client.Object.Delete(context.Background(), key)
	return err
}
