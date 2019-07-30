package slave

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robfig/config"
	. "goReptile/task"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	host string
)

func init() {
	c, _ := config.ReadDefault("config/config.ini")
	host, _ = c.String("web", "host")
}

func getTasks(tag string, version int) ([]Task, error) {
	url := fmt.Sprintf(host+"v1/tasks?tag=%s&version=%d", tag, version)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("StatusCode=%d", resp.StatusCode)
		io.Copy(os.Stderr, resp.Body)
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, nil
}

func postResult(result *TaskResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	_, err = http.Post(host+"v1/results", "application/json", bytes.NewReader(data))
	return err
}
