package workerpool

import (
	"fmt"
	logx "github.com/sirupsen/logrus"
	"sync"
)

// Task 任务的基本行为
type Task interface {
	Process() error
}

// TaskPool 表示任务池结构
type TaskPool struct {
	numConsumers int
	taskChannel  chan Task
	consumerWg   sync.WaitGroup
	isRunning    bool
	mutex        sync.Mutex
}

// Config 定义任务池的配置选项
type Config struct {
	NumConsumers int
	BufferSize   int
}

// NewTaskPool 创建一个新的任务池实例
func NewTaskPool(config Config) (*TaskPool, error) {
	if config.NumConsumers <= 0 {
		logx.Errorf("消费者数量必须大于0")
		return nil, fmt.Errorf("消费者数量必须大于0")
	}
	if config.BufferSize < 0 {
		logx.Errorf("缓冲区大小不能为负数")
		return nil, fmt.Errorf("缓冲区大小不能为负数")
	}

	return &TaskPool{
		numConsumers: config.NumConsumers,
		taskChannel:  make(chan Task, config.BufferSize),
	}, nil
}

// Start 启动任务池
func (tp *TaskPool) Start() error {
	tp.mutex.Lock()
	if tp.isRunning {
		tp.mutex.Unlock()
		logx.Errorf("任务池已经在运行")
		return fmt.Errorf("任务池已经在运行")
	}
	tp.isRunning = true
	tp.mutex.Unlock()

	// 启动消费者
	for i := 0; i < tp.numConsumers; i++ {
		tp.consumerWg.Add(1)
		go tp.startConsumer(i)
	}
	return nil
}

// Stop 停止任务池并等待所有任务完成
func (wp *TaskPool) Stop() {
	wp.mutex.Lock()
	if !wp.isRunning {
		wp.mutex.Unlock()
		return
	}
	wp.mutex.Unlock()
	// 关闭任务通道
	close(wp.taskChannel)
	// 等待所有消费者完成
	wp.consumerWg.Wait()
	wp.mutex.Lock()
	wp.isRunning = false
	wp.mutex.Unlock()
}

// Submit 提交任务到任务池
func (wp *TaskPool) Submit(task Task) error {
	wp.mutex.Lock()
	if !wp.isRunning {
		wp.mutex.Unlock()
		logx.Errorf("任务池未运行")
		return fmt.Errorf("任务池未运行")
	}
	wp.mutex.Unlock()
	wp.taskChannel <- task
	return nil
}

// startConsumer 启动单个消费者
func (wp *TaskPool) startConsumer(id int) {
	defer wp.consumerWg.Done()
	for task := range wp.taskChannel {
		if err := task.Process(); err != nil {
			logx.Errorf("消费者 %d 处理任务时发生错误: %v\n", id, err)
		}
	}
}
