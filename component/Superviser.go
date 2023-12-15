package component

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"
)

// 窗口采样
// 后续用来实现PID控制消费者数量
type Superviser struct {
	conumser_size int
	core_size     int
	max_size      int
	old           window
	Provider      *NetFlow
	Consumers     []*PacketConsumer
}
type window struct {
	last_product  int
	last_consumer int
	produce_count int
	consume_count int
	time_stamp    time.Time
}

// linux执行命令结果
type result struct {
	output []byte
	err    error
}

// iConsumer 初始消费者数量
func Start(iConsumer int) []Superviser {
	netFlows, err := MonitorFactory(65535, -1)
	if err != nil {
		return nil
	}
	res := []Superviser{}
	for i := range netFlows {
		one_superviser := Superviser{}
		netFlows[i].startMonitor()
		provider := netFlows[i]
		for i := 0; i < iConsumer; i++ {
			//添加消费者
			consumer := &PacketConsumer{
				status: true,
				net:    provider,
				ch:     provider.ch_packets,
			}
			one_superviser.Provider = provider
			one_superviser.Consumers = []*PacketConsumer{}
			one_superviser.Consumers = append(one_superviser.Consumers, consumer)
			//监督者组添加监督者
			go consumer.Consume()
		}
		res = append(res, one_superviser)
	}
	return res
}

// 增加消费者
func (s *Superviser) AddConsumer() {
	consumer := &PacketConsumer{
		status: true,
		net:    s.Provider,
		ch:     s.Provider.ch_packets,
	}
	s.Consumers = append(s.Consumers, consumer)
	go consumer.Consume()
}

// 减少消费者
func (s *Superviser) DecConsumer() {
	consumer := s.Consumers[len(s.Consumers)-1]
	consumer.Exit()
	s.Consumers = s.Consumers[:len(s.Consumers)-1]
}

func (s Superviser) PrintInfo() {
	//err := s.delOrTouchFile()
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			for _, v := range s.Provider.Flows {
				print("start PrintInfo", len(s.Provider.Flows))
				print(v.Detail())
			}
		}
	}
}

func (s Superviser) delOrTouchFile() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("Bash file is incompatible with current OS")
	}
	bashFile := "./touch_file"
	cmd := exec.Command("/bin/bash", "-c", bashFile)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Execute Shell:%s failed with error:%s", bashFile, err.Error())
	}
	log.Print("Execute Shell:%s finished with output:\n%s", bashFile, string(output))
	return nil
	//cmd := exec.Command("ls", "/var/log/netmonitor")
	//var stdout, stderr bytes.Buffer
	//cmd.Stdout = &stdout // 标准输出
	//cmd.Stderr = &stderr // 标准错误
	//err := cmd.Run()
	//_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//if err != nil {
	//	command := exec.Command("mkdir", "/usr/log/netmonitor")
	//}
	//exec.Command("")
	//exec.Command("find", "")
}
