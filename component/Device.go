package component

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"sync"
	"time"
)

var (
	Summon uint64
)

type NetFlow struct {
	device     string
	snapLen    int32
	ip         string
	sampleTime time.Duration
	handler    *pcap.Handle
	ch_packets chan gopacket.Packet
	Flows      map[gopacket.Flow]interface{}
}

func MonitorFactory(device string, snapLen int32, sampleTime time.Duration) *NetFlow {
	return &NetFlow{
		device:     device,
		snapLen:    snapLen,
		sampleTime: sampleTime,
	}
}

// todo 一个网络层链接的通道缓存多少合适
func (n *NetFlow) newNetMonitor() {
	//监听网口
	handle, err := pcap.OpenLive(n.device, n.snapLen, false, n.sampleTime)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	//DecodeFragment Fragment contains all
	n.ch_packets = make(chan gopacket.Packet, 65535)
	n.Flows = make(map[gopacket.Flow]interface{})
	packetSource := gopacket.NewPacketSource(handle, gopacket.DecodeFragment)
	//packetSource.DecodeOptions.NoCopy = true;
	for packet := range packetSource.Packets() {
		n.ch_packets <- packet
	}
}

func (n *NetFlow) GetFlow(flow gopacket.Flow) interface{} {
	if n.Flows[flow] == nil {
		var once sync.Once
		once.Do(func() {
			metaFlow := MetaFlow{
				src:       flow.Src().String(),
				dst:       flow.Dst().String(),
				in_total:  0,
				out_total: 0,
				In_Udp:    0,
				Out_udp:   0,
				In_tcp:    0,
				Out_tcp:   0,
				In_ICMP:   0,
				Out_ICMP:  0,
				status:    "",
			}
			n.Flows[flow] = &metaFlow
		})
	}
	return n.Flows[flow]
}
