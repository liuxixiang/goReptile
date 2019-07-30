package task

type Source uint

const (
	SourceUnknown Source = iota
	SourceWebSite
	SourceWechatPublic
	SourceWechatMiniApp
	SourceWeibo
	SourceApp

	TaskNameVideoList = "video-list"
	TaskNameVideo     = "video"
)

type Task struct {
	ID      int    `json:"id"`
	GroupID int    `json:"groupId"`
	Name    string `json:"name"`
	Params  string `json:"params"`
	Nonce   string `json:"nonce"`
	Version int    `json:"version"`
}

type TaskResult struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Result   string `json:"result"`
	Error    string `json:"error"`
	SubTasks []Task `json:"subTasks"`
}

type BaseResult struct {
	Source        Source `json:"source"`
	Origin        string `json:"origin"`
	OriginChannel string `json:"originChannel"`
	Nonce         string `json:"nonce"`
}

type VideoResult struct {
	BaseResult
	Title      string   `json:"title"`
	Cover      string   `json:"cover"`
	Author     string   `json:"author"`
	Avatar     string   `json:"avatar"`
	OriginUrl  string   `json:"originUrl"`
	Url        string   `json:"url"`
	OriginTime int64    `json:"originTime"`
	OriginTags []string `json:"originTags"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Duration   int      `json:"duration"`
	// TODO
}
