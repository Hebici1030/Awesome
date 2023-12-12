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
	sampleTime time.Duration
	handler    *pcap.Handle
	ch_packets chan gopacket.Packet
	Flows      map[gopacket.Flow]FlowMetaInfo
}

// todo 一个网络层链接的通道缓存多少合适
func (n *NetFlow) NewNetMonitor() {
	//监听网口
	handle, err := pcap.OpenLive(n.device, n.snapLen, false, n.sampleTime)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	//DecodeFragment Fragment contains all
	n.ch_packets = make(chan gopacket.Packet, 65535)
	packetSource := gopacket.NewPacketSource(handle, gopacket.DecodeFragment)
	//packetSource.DecodeOptions.NoCopy = true;
	for packet := range packetSource.Packets() {
		n.ch_packets <- packet
	}
}

func (n *NetFlow) GetFlow(flow gopacket.Flow) FlowMetaInfo {
	if n.Flows[flow] == nil {
		var once sync.Once
		once.Do(func() {
			metaFlow := MetaFlow{}
			n.Flows[flow] = &metaFlow
		})
	}
	return n.Flows[flow]
}
