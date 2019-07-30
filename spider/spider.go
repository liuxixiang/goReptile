package spider

import . "goReptile/task"

type VideoSpider interface {
	GetVideoList(params *VideoResult) ([]VideoResult, error)
	GetVideo(params *VideoResult) (*VideoResult, error)
}

var (
	VideoSpiders map[string]VideoSpider
)

func init() {
	VideoSpiders = map[string]VideoSpider{
		"小年糕:推荐":   &XiaoNianGaoSpider{},
		"小年糕:开心":   &XiaoNianGaoSpider{},
		"小年糕:广场舞":  &XiaoNianGaoSpider{},
		"小年糕:祝福":   &XiaoNianGaoSpider{},
		"小年糕:健康":   &XiaoNianGaoSpider{},
		"小年糕:妙招":   &XiaoNianGaoSpider{},
		"票圈长视频:推荐": &TicketRingSpider{},
		"票圈长视频:音乐": &TicketRingSpider{},
		"票圈长视频:综艺": &TicketRingSpider{},
		"票圈长视频:搞笑": &TicketRingSpider{},
		"票圈长视频:祝福": &TicketRingSpider{},
		"票圈长视频:舞蹈": &TicketRingSpider{},
		"票圈长视频:旅行": &TicketRingSpider{},
		"票圈长视频:百态": &TicketRingSpider{},
		"票圈长视频:健康": &TicketRingSpider{},
		"票圈长视频:科技": &TicketRingSpider{},
		"票圈长视频:妙招": &TicketRingSpider{},
		"票圈长视频:影视": &TicketRingSpider{},
		"票圈长视频:美食": &TicketRingSpider{},
		"票圈长视频:时尚": &TicketRingSpider{},
		"票圈长视频:运动": &TicketRingSpider{},
		"票圈长视频:游戏": &TicketRingSpider{},
		"国风:古装": &ChinaStyleSpider{},
	}
}
