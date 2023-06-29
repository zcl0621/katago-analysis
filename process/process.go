package process

import (
	"ktago-server-go/command"
	"ktago-server-go/logger"
	"ktago-server-go/queue"
	"ktago-server-go/syncLock"
	"os"
	"time"
)

var RunMode string = retrieveEnvOrDefault("RUN_MODE", "debug")

var (
	katagoDebugPath = "/home/zhangchaolong/katago-server-go/ktago-server-go/katago"
	katagoDebugArgs = []string{
		"analysis",
		"-model /home/zhangchaolong/katago-server-go/ktago-server-go/best.bin.gz",
		"-config /home/zhangchaolong/katago-server-go/ktago-server-go/analysis.cfg",
	}
	katagoPath = "/usr/bin/katago"
	katagoArgs = []string{
		"analysis",
		"-model /project/best.bin.gz",
		"-config /project/analysis.cfg",
	}
	CMD      *command.Command
	Lock     *syncLock.ChanMutex
	JobQueue = queue.NewQueue(1024)
)

func Init() {
	Lock = syncLock.NewChanMutex()
	Lock.Lock()
	if RunMode == "debug" {
		CMD = command.NewCommand(katagoDebugPath, katagoDebugArgs)
	} else {
		CMD = command.NewCommand(katagoPath, katagoArgs)
	}
	err := CMD.Start()
	if err != nil {
		panic(err)
	}
	for {
		if CMD.Ready {
			Lock.Unlock()
			break
		} else {
			logger.Logger("process.Init", logger.INFO, nil, "waiting for katago ready")
			time.Sleep(time.Second)
		}
	}
	go startQueue()
}

func retrieveEnvOrDefault(key string, defaultValue string) string {
	result := os.Getenv(key)
	if len(result) == 0 {
		result = defaultValue
	}
	return result
}
func startQueue() {
	var j *queue.Iterm
	go func() {
		for {
			select {
			case result := <-CMD.Result:
				if j != nil {
					j.Result = result
				}
			}
		}
	}()
	for {
		job := JobQueue.Pop()
		j = job
		CMD.Input <- job.Cmd
	}
}
