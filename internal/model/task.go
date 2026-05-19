package model

type Task struct {
	// 任务包含:
	// ID
	// 标题:
	// 作者:
	// 内容:
	// 时间:{创建时间，结束时间}
	// 状态:进行中、已完成、失败
	ID      int        `json:"id"`
	Title   string     `json:"title"`
	Author  string     `json:"author"`
	Content string     `json:"content"`
	Time    Time       `json:"time"`
	Status  TaskStatus `json:"status"`
}

type Time struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type TaskStatus string

const (
	TaskStatusRunning   TaskStatus = "running"   // 进行行中
	TaskStatusCompleted TaskStatus = "completed" // 已完成
	TaskStatusFailed    TaskStatus = "failed"    // 失败
)
