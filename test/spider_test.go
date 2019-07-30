package test

import (
	"goReptile/spider"
	. "goReptile/task"
	"testing"
)

func TestSpider(t *testing.T) {
	//spider := &TicketRingSpider{}

	//params := VideoResult{
	//	BaseResult: BaseResult{
	//		Origin:        "票圈长视频",
	//		OriginChannel: "推荐"},
	//}
	spider := &spider.ChinaStyleSpider{}
	params := VideoResult{
		BaseResult: BaseResult{
			Origin:        "国风",
			OriginChannel: "古装"},
	}
	videos, err := spider.GetVideoList(&params)
	if err == nil {
		t.Log(videos)
	} else {
		t.Error(err)
	}

	for _, video := range videos {
		video, err := spider.GetVideo(&video)
		if err == nil {
			t.Log(video)
		} else {
			t.Error(err)
		}
	}

}
