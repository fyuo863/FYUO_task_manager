package handler

import (
	"FYUO_task_manager/internal/model"
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

	c.JSON(http.StatusOK, gin.H{ //返回响应
		"code": 200,
		"data": task,
	})

	// TODO: 推入Redis List

}

func ListTasks(c *gin.Context) {

}
