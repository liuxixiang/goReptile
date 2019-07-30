package utils

import (
	"github.com/alecthomas/log4go"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/robfig/config"
	"net/http"
	"net/url"
	"path"
	"time"
)

type OssParams struct {
	EndpointInternal string
	Endpoint         string
	AccessKeyID      string
	AccessKeySecret  string
	BucketName       string
	PutTimeOut       time.Duration
}

var ossImageParam OssParams
var ossVideoParam OssParams

func init() {
	c, _ := config.ReadDefault("config/config.ini")
	endpointInternal, _ := c.String("ali-oss", "endpointInternal")

	ossVideoParam = OssParams{}
	ossVideoParam.EndpointInternal = endpointInternal
	ossVideoParam.Endpoint = "oss-cn-shanghai.aliyuncs.com"
	ossVideoParam.AccessKeyID = "LTAI51QAyHrOV894"
	ossVideoParam.AccessKeySecret = "GaTILQit6QjKEycWtJif3AX8LABEXH"
	ossVideoParam.BucketName = "xhl-video"
	ossVideoParam.PutTimeOut = 900 * time.Second

	ossImageParam = OssParams{}
	ossImageParam.EndpointInternal = endpointInternal
	ossImageParam.Endpoint = "oss-cn-shanghai.aliyuncs.com"
	ossImageParam.AccessKeyID = "LTAI51QAyHrOV894"
	ossImageParam.AccessKeySecret = "GaTILQit6QjKEycWtJif3AX8LABEXH"
	ossImageParam.BucketName = "xhl-image"
	ossImageParam.PutTimeOut = 900 * time.Second

}

func OssPutImage(originUrl string) (ossUrl string, err error) {
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
	return ossPutObject(originUrl, ossImageParam, hash+ext)
}

func OssPutVideo(originUrl string) (ossUrl string, err error) {
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
	return ossPutObject(originUrl, ossVideoParam, hash+ext)
}

func ossPutObject(originUrl string, ossParam OssParams, objectKey string) (ossUrl string, err error) {
	httpClient := http.Client{Timeout: ossParam.PutTimeOut}
	resp, err := httpClient.Get(originUrl)
	if err != nil {
		log4go.Error(err)
		return
	}
	defer resp.Body.Close()

	ossClient, err := oss.New(ossParam.EndpointInternal, ossParam.AccessKeyID, ossParam.AccessKeySecret)
	if err != nil {
		log4go.Error(err)
		return
	}

	// 获取存储空间
	bucket, err := ossClient.Bucket(ossParam.BucketName)
	if err != nil {
		log4go.Error(err)
		return
	}

	// 上传文件流
	err = bucket.PutObject(objectKey, resp.Body)
	if err != nil {
		log4go.Error(err)
		return
	}

	ossUrl = "https://" + ossParam.BucketName + "." + ossParam.Endpoint + "/" + objectKey
	return
}

func OssPutLocalVideo(filePath string, objectKey string) (ossUrl string, err error) {

	ossClient, err := oss.New(ossVideoParam.EndpointInternal, ossVideoParam.AccessKeyID, ossVideoParam.AccessKeySecret)

	if err != nil {
		log4go.Error(err)
		return
	}

	// 获取存储空间
	bucket, err := ossClient.Bucket(ossVideoParam.BucketName)
	if err != nil {
		log4go.Error(err)
		return
	}

	//上传文件
	err = bucket.PutObjectFromFile(objectKey, filePath)
	if err != nil {
		log4go.Error(err)
		return

	}

	ossUrl = "https://" + ossVideoParam.BucketName + "." + ossVideoParam.Endpoint + "/" + objectKey
	return
}
