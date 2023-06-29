package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"ktago-server-go/logger"
	"ktago-server-go/process"
	"ktago-server-go/queue"
	"ktago-server-go/utils"
	"ktago-server-go/wq"
	"net/http"
	"sync"
	"time"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, "running")
}

func analysis(c *gin.Context) {
	var request analysisScoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusOK, &analysisScoreResponse{
			Id:      request.Id,
			Code:    1,
			Data:    "",
			Message: "数目参数错误",
		})
		return
	}
	cmdline, err := utils.UnzipString(request.Data)
	if err != nil {
		c.JSON(http.StatusOK, &analysisScoreResponse{
			Id:      request.Id,
			Code:    1,
			Data:    "",
			Message: "数目命令错误",
		})
		return
	}
	locked := process.Lock.TryLockWithTimeout(time.Second * 10)
	if !locked {
		c.JSON(http.StatusOK, &analysisScoreResponse{
			Id:      request.Id,
			Code:    1,
			Data:    "",
			Message: "当前数目请求繁忙",
		})
		return
	}
	j := queue.Iterm{
		C:      c,
		Cmd:    cmdline,
		Result: "",
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(c *gin.Context, request *analysisScoreRequest, wg *sync.WaitGroup, j *queue.Iterm) {
		defer wg.Done()
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.JSON(http.StatusOK, &analysisScoreResponse{
					Id:      request.Id,
					Code:    1,
					Data:    "",
					Message: "数目超时",
				})
				return
			default:
				if j.Result != "" {
					var katagoResult katagoAnalysisResult
					e := json.Unmarshal([]byte(j.Result), &katagoResult)
					if e != nil {
						c.JSON(http.StatusOK, &analysisScoreResponse{
							Id:      request.Id,
							Code:    1,
							Data:    "",
							Message: "数目失败",
						})
						return
					}
					resultMap := make(map[string]interface{})
					resultMap["id"] = katagoResult.Id
					resultMap["ownership"] = katagoResult.Ownership
					resultMap["rootInfo"] = katagoResult.RootInfo
					resultBytes, err := json.Marshal(resultMap)
					if err != nil {
						c.JSON(http.StatusOK, &analysisScoreResponse{
							Id:      request.Id,
							Code:    1,
							Data:    "",
							Message: "数目失败",
						})
						return
					}
					c.JSON(http.StatusOK, &analysisScoreResponse{
						Id:      request.Id,
						Code:    0,
						Data:    utils.ZipString(fmt.Sprintf("%s", resultBytes)),
						Message: "success",
					})
					return
				}
			}
		}
	}(c, &request, &wg, &j)
	logger.Logger("analysis", logger.INFO, nil, "set cmdline")
	process.JobQueue.Push(&j)
	wg.Wait()
	process.Lock.Unlock()
	return
}

func demo(c *gin.Context) {
	locked := process.Lock.TryLockWithTimeout(time.Second * 10)
	if !locked {
		c.JSON(http.StatusBadRequest, "wait")
		return
	}
	var request demoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		process.Lock.Unlock()
		c.JSON(http.StatusBadRequest, "sgf error")
		return
	}
	b, _, err := wq.Load(request.SGF)
	if err != nil {
		process.Lock.Unlock()
		c.JSON(http.StatusBadRequest, "sgf error")
		return
	}
	board := b.Board()
	anyData := board.GetKataGoAnalysisData(1)
	cmdline, err := json.Marshal(anyData)
	if err != nil {
		process.Lock.Unlock()
		c.JSON(http.StatusBadRequest, "sgf error")
		return
	}
	j := queue.Iterm{
		C:      c,
		Cmd:    fmt.Sprintf("%s", cmdline),
		Result: "",
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(c *gin.Context, wg *sync.WaitGroup, j *queue.Iterm) {
		defer wg.Done()
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.JSON(http.StatusBadRequest, "timeout")
				return
			default:
				if j.Result != "" {
					var katagoResult katagoAnalysisResult
					e := json.Unmarshal([]byte(j.Result), &katagoResult)
					if e != nil {
						c.JSON(http.StatusBadRequest, "failed")
						return
					}
					resultMap := make(map[string]interface{})
					resultMap["id"] = katagoResult.Id
					resultMap["ownership"] = katagoResult.Ownership
					resultMap["rootInfo"] = katagoResult.RootInfo
					if err != nil {
						c.JSON(http.StatusBadRequest, "failed")
						return
					}
					c.JSON(http.StatusOK, resultMap)
					return
				}
			}
		}
	}(c, &wg, &j)
	logger.Logger("analysis", logger.INFO, nil, "set cmdline")
	process.JobQueue.Push(&j)
	wg.Wait()
	process.Lock.Unlock()
}
