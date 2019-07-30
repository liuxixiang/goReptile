package slave

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/robfig/config"
	"goReptile/spider"
	."goReptile/task"
	"sync"
	"time"
)

var (
	taskPoolSize int
)

func init() {
	c, _ := config.ReadDefault("config/config.ini")
	taskPoolSize, _ = c.Int("task", "taskPoolSize")
}

func Run(tag string, version int) {
	for {
		tasks, err := getTasks(tag, version)
		log4go.Debug("getTask")

		if err != nil {
			log4go.Error(err)
			time.Sleep(time.Second)
			continue
		}

		doTasks(tasks)
	}
}

func doTasks(tasks []Task) {
	var waitGroup sync.WaitGroup
	var poolChan = make(chan int, taskPoolSize)

	for i, task := range tasks {
		waitGroup.Add(1)
		poolChan <- i

		go func(task Task) {
			defer func() {
				waitGroup.Done()
				<-poolChan
			}()

			var taskResult TaskResult
			taskResult.ID = task.ID
			taskResult.Name = task.Name

			switch task.Name {
			case TaskNameVideoList, TaskNameVideo:
				var params VideoResult
				err := json.Unmarshal([]byte(task.Params), &params)
				if nil != err {
					taskResult.Error = err.Error()
					postResult(&taskResult)
					return
				}

				s := spider.VideoSpiders[params.Origin+":"+params.OriginChannel]
				if nil == s {
					taskResult.Error = "Unknown origin: " + params.Origin
					postResult(&taskResult)
					return
				}

				if TaskNameVideoList == task.Name {
					doVideoListTask(s, &task, &params, &taskResult)
				} else {
					doVideoTask(s, &task, &params, &taskResult)
				}

			default:
				taskResult.Error = fmt.Sprint("Unknown task: " + task.Name)
				postResult(&taskResult)
				return
			}
		}(task)
	}

	waitGroup.Wait()
}

func doVideoListTask(s spider.VideoSpider, task *Task, params *VideoResult, taskResult *TaskResult) error {
	videos, err := s.GetVideoList(params)
	if nil != err {
		taskResult.Error = err.Error()
	} else if nil != videos && len(videos) > 0 {
		taskResult.SubTasks = make([]Task, 0, len(videos))
		for _, video := range videos {
			var t Task
			t.GroupID = task.GroupID + 1
			t.Name = TaskNameVideo
			t.Nonce = video.Nonce
			t.Version = task.Version

			v, err := json.Marshal(video)
			if nil != err {
				taskResult.Error = err.Error()
				continue
			}
			t.Params = string(v)

			taskResult.SubTasks = append(taskResult.SubTasks, t)
		}
	} else {
		taskResult.Error = "[]"
	}

	return postResult(taskResult)
}

func doVideoTask(s spider.VideoSpider, task *Task, params *VideoResult, taskResult *TaskResult) error {
	video, err := s.GetVideo(params)
	if nil != err {
		taskResult.Error = err.Error()
	} else {
		v, err := json.Marshal(video)
		if nil != err {
			taskResult.Error = err.Error()
		} else {
			taskResult.Result = string(v)
		}
	}

	return postResult(taskResult)
}
