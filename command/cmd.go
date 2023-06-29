package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	path   string
	args   []string
	Input  chan string
	Output chan string
	Result chan string
	Ready  bool
	cmd    *exec.Cmd
}

func NewCommand(path string, args []string) *Command {
	return &Command{
		path:   path,
		args:   args,
		Input:  make(chan string, 1024),
		Result: make(chan string),
		Output: make(chan string),
		Ready:  false,
	}
}

func (c *Command) Start() error {
	cmd := exec.Command(c.path, c.args...)
	c.cmd = cmd
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		defer close(c.Input)
		for {
			input := <-c.Input
			inputBytes := []byte(input)
			utf8Str := string(inputBytes)
			_, err := fmt.Fprintln(stdin, utf8Str)
			if err != nil {
				fmt.Printf("input %s err: %s\n", input, err.Error())
			} else {
				fmt.Printf("input cmd %s", input)
			}
		}
	}()

	go func() {
		defer close(c.Output)
		changes := make(chan []byte)
		errors := make(chan error)
		logFile := "/tmp/analysis.log"
		for {
			_, err := os.Stat(logFile)
			if err != nil {
				time.Sleep(time.Second)
				continue
			} else {
				break
			}
		}

		go watchFileChanges(logFile, changes, errors)
		for {
			select {
			case data := <-changes:
				c.Output <- fmt.Sprintf("%s", data)
			case err := <-errors:
				fmt.Println("Error:", err)
			}
		}
	}()
	go func() {
		defer close(c.Result)
		for {
			text := <-c.Output
			if text != "" {
				if strings.Contains(text, "Started, ready to begin handling requests") {
					c.Ready = true
				}
				if strings.Contains(text, "moveInfos") {
					x := strings.Split(text, "Response:")
					if len(x) >= 1 {
						c.Result <- x[1]
					}
				}
				fmt.Println(text)
			}
		}
	}()
	return nil
}

func (c *Command) Restart() error {
	return nil
}

func watchFileChanges(filename string, changes chan<- []byte, errors chan<- error) {
	// 打开文件并跳转到文件末尾
	f, err := os.Open(filename)
	if err != nil {
		errors <- err
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		errors <- err
		return
	}

	offset := fi.Size()
	for {
		// 每隔一段时间读取文件内容，如果有变化则将变化写入通道中
		time.Sleep(time.Millisecond * 500)

		fi, err := f.Stat()
		if err != nil {
			errors <- err
			return
		}

		// 文件大小变化表示文件内容发生了变化
		if fi.Size() > offset {
			data := make([]byte, fi.Size()-offset)
			_, err := f.ReadAt(data, offset)
			if err != nil {
				errors <- err
				return
			}

			offset = fi.Size()
			changes <- data
		}
	}
}
