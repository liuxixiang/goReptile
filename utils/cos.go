package utils

import (
	"context"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type CosParams struct {
	AccessKeyID     string
	AccessKeySecret string
	AppID           string
	BucketName      string
	PutTimeOut      time.Duration
	Url             string
	Region          string
	UploadCosKey    string
}

var cosImageParam CosParams
var cosVideoParam CosParams

func init() {
	cosVideoParam = CosParams{}
	cosVideoParam.AccessKeyID = "AKIDGCsHPHzHI6g6GSQ1OVO0LvFI7g04iNgi"
	cosVideoParam.AccessKeySecret = "vb080as0sSnsgDTSEkzzjBcnhRu52i8H"
	cosVideoParam.BucketName = "china-style"
	cosVideoParam.AppID = "1258352729"
	cosVideoParam.Region = "chengdu"
	cosVideoParam.UploadCosKey = "video"
	cosVideoParam.PutTimeOut = 900 * time.Second
	cosVideoParam.Url = "http://" + cosVideoParam.BucketName + "-" + cosVideoParam.AppID + ".cos.ap-" + cosVideoParam.Region + ".myqcloud.com"

	cosImageParam = CosParams{}

	cosImageParam.AccessKeyID = "AKIDGCsHPHzHI6g6GSQ1OVO0LvFI7g04iNgi"
	cosImageParam.AccessKeySecret = "vb080as0sSnsgDTSEkzzjBcnhRu52i8H"
	cosImageParam.BucketName = "china-style"
	cosImageParam.AppID = "1258352729"
	cosImageParam.Region = "chengdu"
	cosImageParam.UploadCosKey = "images"
	cosImageParam.PutTimeOut = 900 * time.Second
	cosImageParam.Url = "http://" + cosImageParam.BucketName + "-" + cosImageParam.AppID + ".cos.ap-" + cosImageParam.Region + ".myqcloud.com"

}

func CosPutImage(originUrl string) (cosUrl string, err error) {
	u, err := url.Parse(originUrl)
	if err != nil {
		log4go.Error(err)
		return
	}

	ext := path.Ext(u.Path)
	if ext == "" {
		ext = ".png"
	}

	hash := Sha1(originUrl)
	return cosPutObject(originUrl, cosImageParam, hash+ext)
}

func CosPutVideo(originUrl string) (cosUrl string, err error) {
	u, err := url.Parse(originUrl)
	if err != nil {
		log4go.Error(err)
		return
	}

	ext := path.Ext(u.Path)
	if ext == "" {
		ext = ".mp4"
	}

	hash := Sha1(originUrl)
	return cosPutObject(originUrl, cosVideoParam, hash+ext)
}

func cosPutObject(originUrl string, cosParams CosParams, objectKey string) (cosUrl string, error error) {
	httpClient := http.Client{Timeout: cosParams.PutTimeOut}
	resp, err := httpClient.Get(originUrl)
	if err != nil {
		log4go.Error(err)
		return
	}
	defer resp.Body.Close()

	u, _ := url.Parse(cosParams.Url)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: cosParams.PutTimeOut,
		Transport: &cos.AuthorizationTransport{
			SecretID:  cosParams.AccessKeyID,
			SecretKey: cosParams.AccessKeySecret,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
	response, err := c.Object.Put(context.Background(), cosParams.UploadCosKey+"/"+objectKey, resp.Body, nil)
	if err != nil {
		panic(err)
		log4go.Error(err)
	}
	cosUrl = response.Request.URL.Host + response.Request.URL.Path
	return
}

func CosPutLocalVideo(filePath string, objectKey string) (cosUrl string, err error) {
	u, _ := url.Parse(cosVideoParam.Url)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: cosVideoParam.PutTimeOut,
		Transport: &cos.AuthorizationTransport{
			SecretID:  cosVideoParam.AccessKeyID,
			SecretKey: cosVideoParam.AccessKeySecret,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
		log4go.Error(err)
		return
	}
	opt := &cos.MultiUploadOptions{
		OptIni:   nil,
		PartSize: 1,
	}
	v, _, err := c.Object.MultiUpload(
		context.Background(), cosVideoParam.UploadCosKey+"/"+objectKey, f, opt,
	)
	if err != nil {
		panic(err)
		log4go.Error(err)
		return
	}
	fmt.Println(v)
	cosUrl = "http://" + v.Location
	return
}
