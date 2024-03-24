package iptables

import (
	configs "Mointor/config"
	"bufio"
	"fmt"
	"github.com/go-co-op/gocron"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

var mapIpExist = make(map[string]bool)

type NetDetail struct {
	pkts        string
	bytes       string
	proto       string
	opt         string
	in          string
	out         string
	source      string
	destination string
}

const IPTABLES_RES_PATH = "/tmp/iptable-res"
const CMD_IPTABLES = "iptables -n -v -L -t filter -x | grep -E 'Chain OUTPUT|Chain INPUT' >> /tmp/iptable-res"

var Schedule *gocron.Scheduler

func StartRecord() {
	Schedule := gocron.NewScheduler(time.UTC)
	_, err := Schedule.Every(configs.Config.Time).Second().Do(func() {
		getAddress()
		getNetFlowInfo()
	})
	if err != nil {
		return
	}
	Schedule.StartAsync()
}
func StopRecord() bool {
	if Schedule == nil {
		return false
	}
	Schedule.Stop()
	return true
}

// 添加监控规则
func AppendIpRule(ip string) bool {
	if mapIpExist[ip] {
		fmt.Sprintf("%v has been append")
	}
	mapIpExist[ip] = true
	//设置输入监听
	_, err := exec.Command("iptables", "-L", "INPUT", "-d", ip).Output()
	if err != nil {
		log.Fatal("Cmd Run With %v /", err)
		return false
	}
	_, err = exec.Command("iptables", "-L", "OUTPUT", "-s", ip).Output()
	if err != nil {
		log.Fatal("Cmd Run With %v /", err)
		return false
	}
	return true
}
func getAddress() {
	// 获取本机的所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error fetching network interfaces:", err)
		return
	}

	// 遍历所有网络接口
	for _, iface := range interfaces {
		// 忽略没有名称的接口和回环接口
		if iface.Name == "" || iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 获取接口的所有地址
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf("Error fetching addresses for interface %s: %v\n", iface.Name, err)
			continue
		}
		// 遍历所有IP地址
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				ip = ip.To16()
			}
			AppendIpRule(ip.String())
		}
	}
}
func getNetFlowInfo() (nets []*NetDetail) {
	cmd := exec.Command("bash", "-c", CMD_IPTABLES)
	_, err := cmd.Output()
	if err != nil {
		log.Fatal("[FIELD]:", CMD_IPTABLES)
		return nil
	}
	file, err := os.Open(IPTABLES_RES_PATH)
	defer file.Close()
	if err != nil {
		log.Fatal("[FIELD]:Connot Open File ", IPTABLES_RES_PATH)
		return nil
	}
	reader := bufio.NewReader(file)
	nets = make([]*NetDetail, 0, 10)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		fields := strings.Fields(line)
		if err != nil {
			return nil
		}
		detail := &NetDetail{
			pkts:        fields[0],
			proto:       fields[1],
			opt:         fields[2],
			in:          fields[3],
			out:         fields[4],
			source:      fields[5],
			destination: fields[6],
		}
		nets = append(nets, detail)
	}
	return nets
}
