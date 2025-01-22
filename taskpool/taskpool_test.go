package workerpool

import (
	"fmt"
	"sync"
	"testing"
)

// DemoTask 接口的具体任务
type DemoTask struct {
	ProducerID int
	TaskID     int
	Data       string
}

// Process 实现 Task 接口的方法
func (t *DemoTask) Process() error {
	return nil
}

// Producer 模拟生产者
func Producer(id int, pool *TaskPool, wg *sync.WaitGroup) {
	defer wg.Done()
	// 模拟动态产生任务
	for i := 0; i < 5; i++ {
		task := &DemoTask{
			ProducerID: id,
			TaskID:     i,
			Data:       fmt.Sprintf("数据-%d", i),
		}
		if err := pool.Submit(task); err != nil {
			fmt.Printf("生产者 %d 提交任务失败: %v\n", id, err)
			return
		}
	}
}

func TestTaskPool(t *testing.T) {
	config := Config{
		NumConsumers: 3,
		BufferSize:   100,
	}
	// 创建任务池
	pool, err := NewTaskPool(config)
	if err != nil {
		fmt.Printf("创建工作池失败: %v\n", err)
		return
	}
	// 启动任务池
	if err := pool.Start(); err != nil {
		fmt.Printf("启动工作池失败: %v\n", err)
		return
	}
	var producerWg sync.WaitGroup
	numProducers := 4
	for i := 0; i < numProducers; i++ {
		producerWg.Add(1)
		go Producer(i, pool, &producerWg)
	}
	producerWg.Wait()
	pool.Stop()
	fmt.Println("所有任务处理完成")
}
