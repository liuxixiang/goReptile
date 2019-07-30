package spider

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/corona10/goimagehash"
	"github.com/jeanphorn/log4go"
	. "goReptile/task"
	. "goReptile/utils"
	"os"
	"runtime"
	"strings"
	"time"
)

var XiaoNianGaoChannelMap = map[string]string{
	"推荐":  "{\"log_params\":{\"page\":\"discover_rec\",\"common\":\"\"},\"qs\":\"jpg\",\"token\":\"7b1456e041d1f3b7cff56f5de18af35a\"}",
	"开心":  "{\"topic_id\":4,\"log_params\":{\"page\":\"discover_happy\",\"common\":\"\"},\"qs\":\"jpg\",\"token\":\"7b1456e041d1f3b7cff56f5de18af35a\"}",
	"广场舞": "{\"topic_id\":5,\"log_params\":{\"page\":\"discover_squareDancing\",\"common\":\"\"},\"qs\":\"jpg\",\"token\":\"7b1456e041d1f3b7cff56f5de18af35a\"}",
	"祝福":  "{\"tag_id\":7,\"log_params\":{\"page\":\"discover_bless\",\"common\":\"\"},\"qs\":\"jpg\",\"token\":\"7b1456e041d1f3b7cff56f5de18af35a\"}",
	"健康":  "{\"topic_id\":8,\"log_params\":{\"page\":\"discover_health\",\"common\":\"\"},\"qs\":\"jpg\",\"token\":\"7b1456e041d1f3b7cff56f5de18af35a\"}",
	"妙招":  "{\"topic_id\":7,\"log_params\":{\"page\":\"discover_trick\",\"common\":\"\"},\"qs\":\"jpg\",\"token\":\"7b1456e041d1f3b7cff56f5de18af35a\"}",
}

type XiaoNianGaoSpider struct {
}

var TempDir string

var TemplateHashList []*goimagehash.ExtImageHash

func init() {

	switch runtime.GOOS {
	case "darwin":
		TempDir = "/tmp/"
	case "linux":
		TempDir = "/tmp/"
	case "windows":
		TempDir = os.Getenv("TMP") + "\\"
	}

	initTemplateHashList()
}

func (s *XiaoNianGaoSpider) GetVideoList(params *VideoResult) (res []VideoResult, err error) {

	resLoc := make([]VideoResult, 0)

	if params.Origin != "小年糕" {
		err = errors.New("origin invalid")
		return
	}

	tag := params.OriginChannel

	body := XiaoNianGaoChannelMap[tag]
	if body == "" {
		err = errors.New("channel invalid")
		return
	}

	type Resp struct {
		Ret  int `json:"ret"`
		Data struct {
			List []struct {
				Title    string `json:"title"`
				Producer string `json:"producer"`
				CoverURL string `json:"url"`
				VideoURL string `json:"v_url"`
				User     struct {
					Hurl string `json:"hurl"`
					Nick string `json:"nick"`
				} `json:"user"`
				T int64 `json:"t"`
			} `json:"list"`
		} `json:"data"`
	}

	url := "https://api.xiaoniangao.cn/trends/get_recommend_trends"
	req := httplib.Post(url)

	req.Header("Content-Type", "application/json; charset=utf-8")
	req.SetTimeout(time.Duration(30)*time.Second, time.Duration(30)*time.Second)

	if body != "" {
		req.Body(body)
	}

	resp := new(Resp)
	err = req.ToJSON(&resp)

	if err != nil {
		log4go.Error(err)
		return
	}

	if resp.Data.List != nil {
		for _, item := range resp.Data.List {

			d := VideoResult{}
			d.Source = SourceWechatMiniApp
			d.Origin = params.Origin
			d.OriginChannel = params.OriginChannel

			params := strings.Split(item.VideoURL, "?") //去除?之后的临时参数
			newUrl := params[0]
			d.Nonce = Sha1(newUrl)
			d.OriginUrl = item.VideoURL
			d.Cover = item.CoverURL
			d.Title = item.Title
			d.Author = item.Producer
			d.Avatar = item.User.Hurl

			d.OriginTime = item.T

			tags := make([]string, 0)
			tags = append(tags, tag)
			d.OriginTags = tags

			resLoc = append(resLoc, d)
		}
	}
	res = resLoc
	return
}

func (s *XiaoNianGaoSpider) GetVideo(params *VideoResult) (videoResult *VideoResult, err error) {

	res := *params

	url := params.OriginUrl

	fmt.Println(url)

	tempParams := strings.Split(url, "?") //去除?之后的临时参数
	newUrl := tempParams[0]
	fileName := Sha1(newUrl)

	originVideoFileName := TempDir + fileName + "-old" + ".mp4"
	newVideoFileName := TempDir + fileName + "-new" + ".mp4"
	newVideoCuttedFileName := TempDir + fileName + "-new-cut" + ".mp4"
	newVideoTailImage := TempDir + fileName + "-tail" + ".jpg"
	objectKey := fileName + ".mp4"

	//删除可能存在的旧文件
	RemoveFile(originVideoFileName)
	RemoveFile(newVideoFileName)
	RemoveFile(newVideoCuttedFileName)
	RemoveFile(newVideoTailImage)

	//函数结束时再删一次
	defer RemoveFile(originVideoFileName)
	defer RemoveFile(newVideoFileName)
	defer RemoveFile(newVideoCuttedFileName)
	defer RemoveFile(newVideoTailImage)

	req := httplib.Get(url)
	req.SetTimeout(time.Duration(30)*time.Second, time.Duration(30)*time.Second)

	err = req.ToFile(originVideoFileName)

	if err != nil {
		log4go.Error(err)
		return nil, err
	} else {
		mediaInfo, err := FFMpegGetMediaInfo(originVideoFileName)
		if err != nil {
			log4go.Error(err)
			return nil, err
		}

		res.Width = int(mediaInfo.Width)
		res.Height = int(mediaInfo.Height)
		res.Duration = int(mediaInfo.Duration)

		bitRate := mediaInfo.BitRate
		w := uint64(92)
		h := uint64(50)

		x := mediaInfo.Width - (w + 19)
		y := mediaInfo.Height - (h + 18)

		err = FFMpegDelogo(originVideoFileName, newVideoFileName, x, y, w, h, bitRate)
		if err != nil {
			log4go.Error(err)
			return nil, err
		}

		if IsFileExist(newVideoFileName) {

			newFileName := newVideoFileName

			RemoveFile(originVideoFileName)

			//files, _ := ListDir("./images", "xng", "")

			secTail := int64(res.Duration/1000) - 1 //最后一秒

			bRemoveTailFlag := false
			bSnapShotPicSuccess := false

			_, err := FFMpegVideScreenShot(newVideoFileName, newVideoTailImage, secTail)
			if err == nil { //生成片尾截图成功
				if IsFileExist(newVideoTailImage) {
					bSnapShotPicSuccess = true
					bRemoveTailFlag = IsNeedToRemoveTail(newVideoTailImage)
				}
			}

			if bRemoveTailFlag { //需要截断尾部，默认为8秒
				secCutEnd := int64(res.Duration/1000) - 8 //截掉最后8秒
				_, err := FFMpegVideoCut(newVideoFileName, newVideoCuttedFileName, 0, secCutEnd)

				if err == nil && IsFileExist(newVideoCuttedFileName) {
					newFileName = newVideoCuttedFileName
				}
			}

			if !bSnapShotPicSuccess { //视频截图失败，不好处理，丢弃
				return nil, errors.New("video process failed")
			}

			fmt.Println(objectKey)

			ossVideoUrl, err := OssPutLocalVideo(newFileName, objectKey)
			ossCoverUrl, err := OssPutImage(res.Cover)
			ossAvatarUrl, err := OssPutImage(res.Avatar)

			/*
			RemoveFile(newVideoFileName)
			RemoveFile(newVideoCuttedFileName)
			RemoveFile(newVideoTailImage)
			*/

			if err != nil {
				return nil, err
			}

			res.Avatar = ossAvatarUrl
			res.Cover = ossCoverUrl
			res.Url = ossVideoUrl

		}

	}

	if res.Url != "" {
		return &res, nil
	} else {
		return nil, errors.New("video process failed")
	}

}

func initTemplateHashList(){

	TemplateHashList = make([]*goimagehash.ExtImageHash, 0)

	files, _ := ListDir("./images", "xng", "")
	width, height := 8, 8

	for _, file := range files {

		img, err := GetImageFromFile(file)
		if err != nil || img == nil {
			continue
		}

		hash1, err := goimagehash.ExtAverageHash(img, width, height)
		if err != nil || hash1 == nil {
			continue
		}
		TemplateHashList = append(TemplateHashList, hash1)
	}

}

func IsNeedToRemoveTail(tailImageFile string)(res bool){

	width, height := 8, 8
	res = false

	img, err := GetImageFromFile(tailImageFile)
	if err != nil || img == nil {
		res = true //安全起见，解析图片异常默认截取
		return
	}
	hash, _ := goimagehash.ExtAverageHash(img, width, height)

	for _, item := range TemplateHashList {

		dis, err := hash.Distance(item)
		if err == nil && dis <= 20 {
			res = true //与某个模板匹配上，需要截断尾部
			break
		}
	}

	return

}
