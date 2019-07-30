package spider

import (
	. "goReptile/task"
	"github.com/astaxie/beego/httplib"
	"errors"
	"strings"
	. "goReptile/utils"
	"github.com/jeanphorn/log4go"
	"runtime"
	"os"
	"net/url"
)

var TicketRingChannelMap = map[string]string{
	"推荐": "?categoryJson=%7B%22categoryId%22%3A55%7D&pageNo=1&pageSize=6&sortField=2&versionCode=96&appType=5",
	"音乐": "?categoryJson={\"categoryId\":3}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"综艺": "?categoryJson={\"categoryId\":2}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"搞笑": "?categoryJson={\"categoryId\":5}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"祝福": "?categoryJson={\"categoryId\":85}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"舞蹈": "?categoryJson={\"categoryId\":4}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"旅行": "?categoryJson={\"categoryId\":9}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"百态": "?categoryJson={\"categoryId\":7}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"健康": "?categoryJson={\"categoryId\":49}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"科技": "?categoryJson={\"categoryId\":11}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"妙招": "?categoryJson={\"categoryId\":45}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"影视": "?categoryJson={\"categoryId\":1}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"美食": "?categoryJson={\"categoryId\":8}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"时尚": "?categoryJson={\"categoryId\":10}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"运动": "?categoryJson={\"categoryId\":12}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
	"游戏": "?categoryJson={\"categoryId\":6}&pageNo=1&pageSize=6&sortField=0&versionCode=96",
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

type ListRsp struct {
	status_code int
	Msg  string
	Data []struct {
		Id        int
		Status    int
		Uid       int
		VideoPath string
		Title     string

		User struct {
			NickName  string
			AvatarUrl string
		}

		CoverImg struct {
			CoverImgPath string
		}

		GmtCreateTimestamp int64
	}
}

type TicketRingSpider struct {
}

func (s *TicketRingSpider) GetVideoList(params *VideoResult) (res []VideoResult, err error) {

	if params.Origin != "票圈长视频" {
		err = errors.New("origin invalid")
		return
	}

	channel := params.OriginChannel

	reqBody := TicketRingChannelMap[channel]

	if reqBody == "" {
		err = errors.New("channel invalid")
		return
	}

	resLoc := make([]VideoResult, 0)

	apiurl := "https://longvideoapi.qingqu.top/longvideoapi/video/distribute/category/videoList/v2?categoryJson=%7B%22categoryId%22%3A55%7D&pageNo=1&pageSize=6&sortField=2&versionCode=96&appType=5"

	reqBody = url.QueryEscape(reqBody)

	req := httplib.Post(apiurl)

	listRsp := new(ListRsp)

	req.ToJSON(&listRsp)

	for _, item := range listRsp.Data {

		d := VideoResult{}
		d.Source = SourceWechatMiniApp
		d.Origin = params.Origin
		d.OriginChannel = params.OriginChannel

		params := strings.Split(item.VideoPath, "?") //去除?之后的临时参数
		newUrl := params[0]
		d.Nonce = Sha1(newUrl)
		d.OriginUrl = item.VideoPath
		d.Cover = item.CoverImg.CoverImgPath
		d.Title = item.Title
		d.Author = item.User.NickName
		d.Avatar = item.User.AvatarUrl

		d.OriginTime = item.GmtCreateTimestamp

		tags := make([]string, 0)
		tags = append(tags, channel)
		d.OriginTags = tags

		resLoc = append(resLoc, d)

	}
	res = resLoc
	return

}

func (s *TicketRingSpider) GetVideo(params *VideoResult) (res *VideoResult, err error) {
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
		err = errors.New("暂不支持MP4与M3U8的视频处理")
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

	ossCoverUrl, err := OssPutImage(res.Cover)
	ossAvatarUrl, err := OssPutImage(res.Avatar)
	ossVideoUrl, err := OssPutLocalVideo(finalFilePath, finalFilePath[strings.LastIndex(finalFilePath, "/")+1:len(finalFilePath)])

	if err != nil {
		log4go.Error(err)
		return
	}

	res.Avatar = ossAvatarUrl
	res.Cover = ossCoverUrl
	res.Url = ossVideoUrl

	log4go.Info(res.OriginUrl, "oos上传地址->", res.Url)

	return
}
