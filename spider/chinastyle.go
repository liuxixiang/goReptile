package spider

import (
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/jeanphorn/log4go"
	. "goReptile/task"
	. "goReptile/utils"
	"net/url"
	"os"
	"runtime"
	"strings"
)

var chinaStyleChannelMap = map[string]string{
	"古装": "?categoryJson=%7B%22categoryId%22%3A55%7D&pageNo=1&pageSize=6&sortField=2&versionCode=96&appType=5",
}

func init() {

	switch runtime.GOOS {
	case "darwin":
		TempDir = "/tmp/"
	case "linux":
		TempDir = "/tmp/"
	case "windows":
		TempDir = os.Getenv("TMP") + "\\"
	}
}

type ChinaStyleListRsp struct {
	statusCode int `json:"status_code"`
	Data       [] struct {
		Data struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			CreateTime  int64  `json:"create_time"`

			Video struct {
				VideoId string `json:"video_id"`
				Cover   struct {
					Url [] string `json:"url_list"`
				} `json:"cover"`
			} `json:"video"`

			Author struct {
				NickName   string `json:"nickname"`
				AvatarThum struct {
					Url [] string `json:"url_list"`
				} `json:"avatar_thumb"`
			} `json:"author"`
		} `json:"data"`
	} `json:"data"`
	Extra struct {
	} `json:"extra"`
}

type ChinaStyleSpider struct {
}

func (s *ChinaStyleSpider) GetVideoList(params *VideoResult) (res []VideoResult, err error) {

	if params.Origin != "国风" {
		err = errors.New("origin invalid")
		return
	}

	channel := params.OriginChannel

	reqBody := chinaStyleChannelMap[channel]

	if reqBody == "" {
		err = errors.New("channel invalid")
		return
	}

	resLoc := make([]VideoResult, 0)

	apiurl := "https://hotsoon-hl.snssdk.com/hotsoon/hashtag/1593654619601933/items/?count=10"

	reqBody = url.QueryEscape(reqBody)

	req := httplib.Post(apiurl)

	listRsp := new(ChinaStyleListRsp)

	req.ToJSON(&listRsp)

	for _, item := range listRsp.Data {

		d := VideoResult{}
		d.Source = SourceWechatMiniApp
		d.Origin = params.Origin
		d.OriginChannel = params.OriginChannel
		videoPath := "https://api-hl.huoshan.com/hotsoon/item/video/_playback/?video_id=" + item.Data.Video.VideoId
		params := strings.Split(videoPath, "?") //去除?之后的临时参数
		newUrl := params[0]
		d.Nonce = Sha1(newUrl)
		d.OriginUrl = videoPath
		d.Cover = item.Data.Video.Cover.Url[0]
		d.Title = item.Data.Title
		d.Author = item.Data.Author.NickName
		d.Avatar = item.Data.Author.AvatarThum.Url[0]

		d.OriginTime = item.Data.CreateTime

		tags := make([]string, 0)
		tags = append(tags, channel)
		d.OriginTags = tags

		resLoc = append(resLoc, d)

	}
	res = resLoc
	return

}

func (s *ChinaStyleSpider) GetVideo(params *VideoResult) (res *VideoResult, err error) {
	res = params
	url := params.OriginUrl
	log4go.Info("开始处理爬虫视频链接", url)

	var finalFilePath string

	videoType := VideoType(url)

	switch videoType {

	case MP4:
		finalFilePath, err = DownloadMp4(url, TempDir)
		break
	case M3U8:
		finalFilePath, err = FfmpegM3u8ConverMp4(url, TempDir)
		break
	case OTHER:
		//err = errors.New("暂不支持MP4与M3U8的视频处理")
		finalFilePath, err = DownloadMp4(url, TempDir)
		break
	}

	if err != nil {
		log4go.Error(err)
		return
	}

	mediaInfo, err := FFMpegGetMediaInfo(finalFilePath)

	if err != nil {
		log4go.Error(err)
		return
	}

	res.Width = int(mediaInfo.Width)
	res.Height = int(mediaInfo.Height)
	res.Duration = int(mediaInfo.Duration)

	cosCoverUrl, err := CosPutImage(res.Cover)
	cosAvatarUrl, err := CosPutImage(res.Avatar)
	cosVideoUrl, err := CosPutLocalVideo(finalFilePath, finalFilePath[strings.LastIndex(finalFilePath, "/")+1:len(finalFilePath)])

	if err != nil {
		log4go.Error(err)
		return
	}

	res.Avatar = cosAvatarUrl
	res.Cover = cosCoverUrl
	res.Url = cosVideoUrl

	log4go.Info(res.OriginUrl, "oos上传地址->", res.Url)

	return
}
