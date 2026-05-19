package handler

import (
	"FYUO_task_manager/internal/database"
	"FYUO_task_manager/internal/model"
	"FYUO_task_manager/internal/queue"
	"FYUO_task_manager/pkg/log"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// 创建新任务
// 内容

func CreateTask(c *gin.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil { //接收请求体并绑定到结构体
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	task := model.Task{
		Title:   req.Title,
		Author:  "占位",
		Content: req.Content,
		Time: model.Time{
			StartTime: time.Now().Format("2006-01-02 15:04:05"),
			EndTime:   "",
		},
		Status: model.TaskStatusRunning,
	}

	// c.JSON(http.StatusOK, gin.H{ //返回响应
	// 	"code": 200,
	// 	"data": task,
	// })

	// TODO: 推入Redis List
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置超时上下文，防止Redis操作阻塞过久
	defer cancel()
	queue.EnqueueTask(ctx, database.RDB, &task) //将任务推入Redis队列，供Worker异步处理

	length, err := database.RDB.LLen(ctx, queue.TaskQueueKey).Result()
	if err != nil {
		//fmt.Errorf("获取队列长度失败: %w", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{ //返回响应
		"code":         200,
		"data":         task,
		"queue_length": length,
	})

}

// 测试手动从队列取出任务
func PopTasks(c *gin.Context) {
	task, err := queue.DequeueTask(context.Background(), database.RDB)
	if err != nil {
		log.Logger.Error("[Handler]", "从队列取出任务失败: ", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) //设置超时上下文，防止Redis操作阻塞过久
	defer cancel()
	length, err := database.RDB.LLen(ctx, queue.TaskQueueKey).Result()
	if err != nil {
		//fmt.Errorf("获取队列长度失败: %w", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{ //返回响应
		"code":         200,
		"data":         task,
		"queue_length": length,
	})
}
