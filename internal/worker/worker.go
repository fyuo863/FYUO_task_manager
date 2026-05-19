package worker

import (
	"FYUO_task_manager/internal/model"
	"FYUO_task_manager/internal/queue"
	"FYUO_task_manager/pkg/log"
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

// Pool 管理一组固定数量的 Goroutine，通过 Channel 分发任务
type Pool struct {
	workers  int              // 工作协程数量
	taskChan chan *model.Task // 任务 Channel，带缓冲以削峰
	wg       sync.WaitGroup   // 等待所有 Worker 完成
}

// NewPool 创建协程池并指定并发数
func NewPool(numWorkers int, bufferSize int) *Pool {
	return &Pool{
		workers:  numWorkers,
		taskChan: make(chan *model.Task, bufferSize), // 带缓冲 Channel，防止短暂积压时阻塞入队
	}
}

// Start 启动所有工作协程，每个协程从 Channel 中竞争领取任务
// 使用 sync.WaitGroup 追踪所有协程，确保优雅停机时等待全部完成
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i) // i 是 Worker 编号，用于日志区分
	}
	log.Logger.Info("[Worker] 协程池已启动", "Worker", p.workers)
}

func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	log.Logger.Info("[Worker] 已启动", "ID", id)
	for task := range p.taskChan { //worker从taskChan取数据
		select {
		case <-ctx.Done(): //如果上下文触发退出
			log.Logger.Info("[Worker] 停止", "ID", id)
			return
		default:
			// 无事发生
		}
		//TODO: 处理task
		p.processTaskTest(ctx, task, id)

	}
}

// --- 搬运工部分：负责把 Redis 里的东西搬到协程池里 ---
func (p *Pool) ListenQueue(ctx context.Context, rdb *redis.Client) {
	log.Logger.Info("[Dispatcher] 开始监听 Redis 队列...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// 1. 从 Redis 阻塞读取 (超时时间设为 0)
			msg, err := queue.DequeueTask(ctx, rdb)
			if err != nil {
				log.Logger.Error("出队失败", "err", err)
				continue
			}

			// 2. 转换为 model.Task (根据需要从数据库查一遍或直接转换)
			task := &model.Task{
				ID:    msg.ID,
				Title: msg.Title,
			}

			// 3. 投递到协程池的 Channel 中
			p.Submit(task)
		}
	}
}

func (p *Pool) Submit(task *model.Task) { //将任务发送到Worker Pool的Channel，供Worker处理
	p.taskChan <- task
}

// func (p *Pool) processTask(ctx context.Context, task *model.Task, workerID int) { //从消息队列取出Task
// 	log.Logger.Info("[Worker] 处理任务", "WorkerID", workerID, "TaskID", task.ID)
// 	task, err := queue.DequeueTask(0, database.RDB)
// 	if err != nil {
// 		log.Logger.Error("[Handler]", "从队列取出任务失败: ", err)
// 		return
// 	}
// }

func (p *Pool) processTaskTest(ctx context.Context, task *model.Task, workerID int) { //从消息队列取出Task
	select {
	case <-ctx.Done():
		return
	default:
		log.Logger.Info("[Worker] 处理任务", "WorkerID", workerID, "TaskID", task.ID)
		log.Logger.Info("[Worker]", "任务详情", task)
	}
}

func (p *Pool) Stop() {
	close(p.taskChan) // 关闭 Channel，Worker 的 range 循环检测到关闭后退出
	p.wg.Wait()       // 等待所有 Worker 的 goroutine 完全退出
	log.Logger.Info("[Worker] 所有 Worker 已退出")
}
