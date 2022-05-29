package tencent

import (
	"api/common"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/vars"
	"net/http"
	"net/url"
	"time"
)

type Service struct {
	*common.Inject
	Vars *vars.Service
}

func (x *Service) Client(ctx context.Context) (client *cos.Client, err error) {
	secretId, err := x.Vars.GetTencentSecretId(ctx)
	if err != nil {
		return
	}
	secretKey, err := x.Vars.GetTencentSecretKey(ctx)
	if err != nil {
		return
	}
	bucket, err := x.Vars.GetTencentCosBucket(ctx)
	if err != nil {
		return
	}
	region, err := x.Vars.GetTencentCosRegion(ctx)
	if err != nil {
		return
	}
	var u *url.URL
	if u, err = url.Parse(
		fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`, bucket, region),
	); err != nil {
		return
	}
	client = cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})
	return
}

// CosPresigned 对象存储预签名
func (x *Service) CosPresigned(ctx context.Context) (data interface{}, err error) {
	secretId, err := x.Vars.GetTencentSecretId(ctx)
	if err != nil {
		return
	}
	secretKey, err := x.Vars.GetTencentSecretKey(ctx)
	if err != nil {
		return
	}
	bucket, err := x.Vars.GetTencentCosBucket(ctx)
	if err != nil {
		return
	}
	expired, err := x.Vars.GetTencentCosExpired(ctx)
	if err != nil {
		return
	}
	date := time.Now()
	keyTime := fmt.Sprintf(`%d;%d`, date.Unix(), date.Add(expired).Unix())
	key := fmt.Sprintf(`%s/%s/%s`,
		x.Values.Namespace,
		date.Format("20060102"),
		helper.Uuid(),
	)
	policy := map[string]interface{}{
		"expiration": date.Add(expired).Format("2006-01-02T15:04:05.000Z"),
		"conditions": []interface{}{
			map[string]interface{}{"bucket": bucket},
			[]interface{}{"starts-with", "$key", key},
			map[string]interface{}{"q-sign-algorithm": "sha1"},
			map[string]interface{}{"q-ak": secretId},
			map[string]interface{}{"q-sign-time": keyTime},
		},
	}
	var policyText []byte
	if policyText, err = jsoniter.Marshal(policy); err != nil {
		return
	}
	signKeyHash := hmac.New(sha1.New, []byte(secretKey))
	signKeyHash.Write([]byte(keyTime))
	signKey := hex.EncodeToString(signKeyHash.Sum(nil))
	stringToSignHash := sha1.New()
	stringToSignHash.Write(policyText)
	stringToSign := hex.EncodeToString(stringToSignHash.Sum(nil))
	signatureHash := hmac.New(sha1.New, []byte(signKey))
	signatureHash.Write([]byte(stringToSign))
	signature := hex.EncodeToString(signatureHash.Sum(nil))
	return gin.H{
		"key":              key,
		"policy":           policyText,
		"q-sign-algorithm": "sha1",
		"q-ak":             secretId,
		"q-key-time":       keyTime,
		"q-signature":      signature,
	}, nil
}

func (x *Service) ImageInfo(ctx context.Context, url string) (result map[string]interface{}, err error) {
	var client *cos.Client
	if client, err = x.Client(ctx); err != nil {
		return
	}
	var response *cos.Response
	if response, err = client.CI.Get(ctx, url, "imageInfo", nil); err != nil {
		return
	}
	if err = jsoniter.NewDecoder(response.Body).Decode(&result); err != nil {
		return
	}
	return
}
