package queue

import (
	"FYUO_task_manager/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// TaskMessage 放入 Redis 队列的消息体——只包含 Worker 需要的核心字段
type TaskMessage struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	UserID uint   `json:"user_id"`
}

const (
	// TaskQueueKey Redis List 的键名，API 和 Worker 通过这个名字通信
	TaskQueueKey = "task:queue"
)

// EnqueueTask 将任务序列化为 JSON 后 LPUSH 到 Redis 队列
// API Service 在创建任务后调用，投递给 Worker 异步处理
func EnqueueTask(ctx context.Context, rdb *redis.Client, task *model.Task) error {
	msg := TaskMessage{
		ID:    task.ID,
		Title: task.Title,
		//UserID: task.UserID,
	}

	data, err := json.Marshal(msg) // Go 结构体 → JSON 字节序列
	if err != nil {
		return fmt.Errorf("序列化任务消息失败: %w", err)
	}

	// LPUSH key value：将消息推入列表左端
	// Worker 端用 BRPOP 从右端取出，实现 FIFO 先进先出
	if err := rdb.LPush(ctx, TaskQueueKey, data).Err(); err != nil {
		return fmt.Errorf("任务入队失败: %w", err)
	}

	return nil
}

// DequeueTask 阻塞式从 Redis 队列取出一个任务
// timeout=0 表示永久阻塞等待，直到有任务入队
// Worker 主循环调用此函数等待新任务
func DequeueTask(ctx context.Context, rdb *redis.Client) (*TaskMessage, error) {
	// BRPOP key timeout：从列表右端弹出元素
	// timeout=0 表示永久阻塞直到有数据
	result, err := rdb.BRPop(ctx, time.Duration(0), TaskQueueKey).Result()
	if err != nil {
		return nil, fmt.Errorf("任务出队失败: %w", err)
	}

	// BRPOP 返回 [key, value]，result[0] 是键名，result[1] 是值
	var msg TaskMessage
	if err := json.Unmarshal([]byte(result[1]), &msg); err != nil {
		return nil, fmt.Errorf("反序列化任务消息失败: %w", err)
	}

	return &msg, nil
}
