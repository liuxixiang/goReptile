package utils

import (
	"net/url"
	"net/http"
	"github.com/alecthomas/log4go"
	"time"
	"io"
	"os"
	"github.com/golang/groupcache/lru"
	"github.com/grafov/m3u8"
	"errors"
	"strings"
)

type VideoEnum int

const (
	MP4   VideoEnum = 0
	M3U8  VideoEnum = 1
	OTHER VideoEnum = 5
)

type Download struct {
	URI           string
	totalDuration time.Duration
}

var USER_AGENT string
var client = &http.Client{}

func doRequest(c *http.Client, req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := c.Do(req)
	return resp, err
}

func downloadm3u8(filePath string, msChan chan *Download) {

	out, err := os.Create(filePath)

	defer out.Close()

	if err != nil {
		log4go.Error(errors.New(filePath + "创建文件报错"))
		return
	}
	log4go.Info("Downloading")

	for v := range msChan {
		req, err := http.NewRequest("GET", v.URI, nil)
		if err != nil {
			log4go.Error(err)
		}
		resp, err := doRequest(client, req)
		if err != nil {
			log4go.Error(err)
			continue
		}
		if resp.StatusCode != 200 {
			log4go.Info("Received HTTP %v for %v\n", resp.StatusCode, v.URI)
			continue
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log4go.Error(err)
		}
		resp.Body.Close()
		//log4go.Info("Downloaded %v\n", v.URI)
		//log4go.Info("Recorded %v of %v\n", v.totalDuration, v.totalDuration)
	}

	log4go.Info("DownLoaded ！")

}

func DownloadMp4(url string, tempDir string) (filePath string, err error) {

	filePath = tempDir + Sha1(url) + ".mp4";
	out, err := os.Create(filePath)

	defer out.Close()

	if err != nil {
		err = errors.New(filePath + "创建文件报错")
		log4go.Error(err)
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log4go.Error(err)
		return
	}
	resp, err := doRequest(client, req)
	if err != nil {
		log4go.Error(err)
		return

	}
	if resp.StatusCode != 200 {
		err = errors.New("调用视频地址报错" + url)
		return

	}
	_, err = io.Copy(out, resp.Body)
	resp.Body.Close()

	if err != nil {
		log4go.Error(err)
		return
	}

	log4go.Info("Downloaded other file %v\n", url+"----"+out.Name())
	return
}


func M3u8ConvertMp4(fileUrl string, tempDir string) (filePath string, err error) {
	var recDuration time.Duration = 0

	filePath = tempDir + Sha1(fileUrl) + ".mp4";
	msChan := make(chan *Download, 1024)

	//记录循环
	cache := lru.New(1024)

	playlistUrl, err := url.Parse(fileUrl)
	if err != nil {
		log4go.Error(err)
		return
	}
	log4go.Info("格式化请求链接", playlistUrl.String())

	req, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		log4go.Error(err)
		return
	}
	resp, err := doRequest(client, req)
	if resp.StatusCode != 200 {
		log4go.Error(err)
		return
	}
	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)

	if listType != m3u8.MEDIA {
		log4go.Info("视频类型非m3u8")
		err = errors.New("视频类型非m3u8");
		return
	}

	mpl := playlist.(*m3u8.MediaPlaylist)

	for _, v := range mpl.Segments {
		//防止这里死循环
		if v != nil {

			var msURI string

			msUrl, err := playlistUrl.Parse(v.URI)
			if err != nil {
				log4go.Error(err)
				continue
			}
			msURI, err = url.QueryUnescape(msUrl.String())
			if err != nil {
				log4go.Error(err)
			}

			recDuration += time.Duration(int64(v.Duration * 1000000000))

			_, hit := cache.Get(msURI)
			if !hit {
				cache.Add(msURI, nil)
				msChan <- &Download{msURI, recDuration}
			}
		}
	}

	if mpl.Closed {
		close(msChan)
	}

	downloadm3u8(filePath, msChan)

	return
}

func VideoType(vidopath string) (VideoEnum) {

	if strings.Index(vidopath, ".mp4") >= 0 {
		return MP4
	}

	if strings.Index(vidopath, ".m3u8") >= 0 {
		return M3U8
	}

	return OTHER

}
