package component

import (
	"Awesome/component/model"
	"Awesome/utils"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Comsumer interface {
	Consume()
	Exit()
}
type PacketConsumer struct {
	status bool
	net    *NetFlow
	ch     chan gopacket.Packet
	////数据解析曾
	//all_layer []gopacket.DecodingLayer
}

func (consumer *PacketConsumer) Consume() {
	var ip4 layers.IPv4
	var ip6 layers.IPv6
	var tcp layers.TCP
	var udp layers.UDP
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &ip4, &ip6, &tcp, &udp)
	layerData := []gopacket.LayerType{}
	for {
		if !consumer.status {
			return
		}
		fmt.Printf("%v,%v,%v,%v,%v \n", consumer.net, consumer.ch, cap(consumer.ch), len(consumer.ch), consumer.status)
		packet := <-consumer.ch
		//TODO
		//layer := packet.NetworkLayer()
		//flow := layer.NetworkFlow()
		//metaInfo := receiver.parent.GetFlow(flow)
		//meta := metaInfo.(*model.MetaFlow)
		//b_InData := (receiver.parent.device.Addresses == meta.Dst)
		err := parser.DecodeLayers(packet.Data(), &layerData)
		if err != nil {
			print("解析协议层失败，继续解析下一个数据包\n")
			continue
		}
		slice2Map := utils.ConvertSlice2Map(layerData)
		var meta *model.MetaFlow
		var metaInfo model.FlowMetaInfo

		if utils.InMap(slice2Map, layers.LayerTypeIPv4) {
			flow := ip4.NetworkFlow()
			metaInfo = consumer.net.GetMetaInfoByFlow(flow)
		} else {
			flow := ip6.NetworkFlow()
			metaInfo = consumer.net.GetMetaInfoByFlow(flow)
		}
		meta = metaInfo.(*model.MetaFlow)
		for _, ltype := range layerData {
			switch ltype {
			case layers.LayerTypeUDP:
				meta.In_total++
			case layers.LayerTypeTCP:
				meta.In_total++
			}
		}
	}
}

func (receiver *PacketConsumer) Exit() {
	receiver.status = false
}
