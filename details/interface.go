package details

import (
	"sync"
	"time"
)

// Ncbi 定义接口
type Ncbi interface {
	Get() error
	Json() string
	Save(string) error
}

func get(input chan Ncbi, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task, ok := <-input

		if !ok {
			break
		}

		for {
			if err := task.Get(); err != nil {
				sugar.Warn(err)
				time.Sleep(time.Second * 5)
			} else {
				break
			}
		}
	}
}
