package component

type FlowMetaInfo interface {
	Detail() string
	Refresh() bool
}
type MetaFlow struct {
	src       string
	dst       string
	In_total  uint64
	Out_total uint64
	In_Udp    int64
	Out_udp   int64
	In_tcp    int64
	Out_tcp   int64
	In_ICMP   int64
	Out_ICMP  int64
	status    string
}

func (m *MetaFlow) Detail() string {
	return "nil"
}
func (m *MetaFlow) Refresh() bool {
	return true
}

//func MetaFlowFactory(flow gopacket.Flow) MetaFlow {
//	return
//}
