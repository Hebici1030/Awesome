package component

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

const (
	//正规表达式匹配网口
	netgape string = "^\\w*:\\sflags=[0-9]*9"
)

var SystemNet = []string{}

func DeviceFiner() error {
	compile, err := regexp.Compile(netgape)
	if err != nil {
		return fmt.Errorf("执行ifconfig查找网口失败")
	}
	cmd := exec.Command("ifconfig")
	pipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("执行ifconfig查找网口失败")
	}
	defer pipe.Close()
	reader := bufio.NewReader(pipe)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("")
		}
		if find := compile.Find(line); find != nil {
			index := bytes.IndexAny(find, ":")
			if index != -1 {
				SystemNet = append(SystemNet, string(find[0:index]))
			}
		}
	}
	if len(SystemNet) == 0 {
		return fmt.Errorf("执行ifconfig查找网口失败")
	}
	return nil
}
