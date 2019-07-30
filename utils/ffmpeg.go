package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"github.com/alecthomas/log4go"
)

const (
	WINDOWS = 1
	LINUX   = 2
	MACOS   = 3
)

type MediaInfo struct {
	Width    uint64
	Height   uint64
	Duration uint64
	BitRate  uint64
}

var OSType int

func init() {
	switch runtime.GOOS {
	case "darwin":
		OSType = MACOS
	case "windows":
		OSType = WINDOWS
	case "linux":
		OSType = LINUX
	}

}

func FFMpegGetMediaInfo(locFilePath string) (mediaInfo MediaInfo, err error) {

	res := MediaInfo{}

	info := ""
	cmdParam := fmt.Sprintf("ffmpeg -i %s", locFilePath)

	cmdRes, _ := ExecFFmpegCmd(cmdParam)
	info = string(cmdRes)

	if info == "" {
		err = errors.New("invalid info")
		return res, err
	}

	regSize := regexp.MustCompile(`[1-9]+[0-9]*x[1-9]+[0-9]*`)
	vSizeStr := regSize.FindString(info)

	vSize := strings.Split(vSizeStr, "x")

	width := uint64(0)
	height := uint64(0)

	if len(vSize) == 2 {
		width, _ = strconv.ParseUint(vSize[0], 10, 64)
		height, _ = strconv.ParseUint(vSize[1], 10, 64)
	} else {
		err = errors.New("invalid vSize from info")
		return res, err
	}

	if width == 0 || height == 0 {
		err = errors.New("invalid width or height")
		return res, err
	} else {
		res.Width = width
		res.Height = height
	}

	regDuration := regexp.MustCompile(`[0-9]{2}:[0-9]{2}:[0-9]{2}\.[0-9]{2}`)
	vDurationStr := regDuration.FindString(info)

	vDuration := strings.Split(vDurationStr, ":")
	if len(vDuration) == 3 {
		hour, _ := strconv.ParseUint(vDuration[0], 10, 64)
		min, _ := strconv.ParseUint(vDuration[1], 10, 64)

		secs := strings.Split(vDuration[2], ".")
		if len(secs) == 2 {
			sec, _ := strconv.ParseUint(secs[0], 10, 64)
			msec, _ := strconv.ParseUint(secs[1], 10, 64)

			duration := hour*3600*1000 + min*60*1000 + sec*1000 + msec*10
			res.Duration = duration

			mediaInfo = res

		} else {
			err = errors.New("video duration parse failed")
		}

	} else {
		err = errors.New("invalid video duration")
		return res, err
	}

	regBitRate := regexp.MustCompile(`[1-9]+[0-9]* kb/s`)
	vBitRateStr := regBitRate.FindAllString(info, -1)

	videoBitRate := uint64(0)
	if len(vBitRateStr) == 3 {
		//windows & linux
		videoBitRateStr := strings.Split(vBitRateStr[1], " ")
		videoBitRate, _ = strconv.ParseUint(videoBitRateStr[0], 10, 64)
		res.BitRate = videoBitRate
	} else if len(vBitRateStr) == 2 {
		//mac
		videoBitRateStr := strings.Split(vBitRateStr[0], " ")
		videoBitRate, _ = strconv.ParseUint(videoBitRateStr[0], 10, 64)
		res.BitRate = videoBitRate
	} else {
		err = errors.New("invalid bitrate")
		return res, err
	}

	return res, nil
}


func FfmpegM3u8ConverMp4(fileUrl string, tempDir string) (filePath string, err error) {

	filePath = tempDir + Sha1(fileUrl) + ".mp4";
	cmdParam := fmt.Sprintf("ffmpeg -i \"%s\" -vcodec copy -acodec copy -absf aac_adtstoasc  %s", fileUrl, filePath)

	cmdRes, err := ExecFFmpegCmd(cmdParam)

	if err != nil {
		return
	}
	info := string(cmdRes)
	log4go.Info("m3u8在线地址", fileUrl, "通过ffmpeg工具转换mp4返回结果", info)
	return
}

func FFMpegDelogo(oldFilePath string, newFilePath string, x uint64, y uint64, w uint64, h uint64, bitRate uint64) (err error) {

	cmdParam := fmt.Sprintf("ffmpeg -i %s -b:v %dk -filter_complex delogo=x=%d:y=%d:w=%d:h=%d:show=0 %s", oldFilePath, bitRate, x, y, w, h, newFilePath)

	_, err = ExecFFmpegCmd(cmdParam)

	return err
}

func FFMpegVideoCut(oldFilePath string, newFilePath string, secStart int64, secEnd int64)(res []byte, err error){

	cmdParam := fmt.Sprintf("ffmpeg -ss %s -accurate_seek -i %s -to %s  -codec copy -avoid_negative_ts 1 %s", GetHHMMSS(secStart), oldFilePath, GetHHMMSS(secEnd), newFilePath)
	res, err = ExecFFmpegCmd(cmdParam)
	return
}

func FFMpegVideScreenShot(videoPath string, imagePath string, secTime int64)(res []byte, err error){
	cmdParam := fmt.Sprintf("ffmpeg -ss %s -i %s -frames:v 1 -f image2 %s",GetHHMMSS(secTime), videoPath, imagePath)
	res, err = ExecFFmpegCmd(cmdParam)
	return
}

func GetHHMMSS(secInt int64)(res string){

	hour := int64(secInt / 3600)
	min := int64((secInt % 3600) / 60)
	sec := int64(secInt % 60)

	res = fmt.Sprintf("%02d:%02d:%02d",hour, min, sec)
	return

}

func ExecFFmpegCmd(cmdStr string)(res []byte, err error){

	if OSType == WINDOWS {
		cmd := exec.Command("cmd", "/C", cmdStr)
		res, err = cmd.CombinedOutput()

	}

	if OSType == LINUX || OSType == MACOS {
		cmd := exec.Command("sh", "-c", cmdStr)
		res, err = cmd.CombinedOutput()
	}

	return
}