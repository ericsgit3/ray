
import (
	"encoding/json"
	"os"
	"strings"

	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
)

type StringList []string

func NewStringList(raw []string) *StringList {
	list := StringList(raw)
	return &list
}

func (v StringList) Len() int {
	return len(v)
}

func (v *StringList) UnmarshalJSON(data []byte) error {
	var strarray []string
	if err := json.Unmarshal(data, &strarray); err == nil {
		*v = *NewStringList(strarray)
		return nil
	}

	var rawstr string
	if err := json.Unmarshal(data, &rawstr); err == nil {
		strlist := strings.Split(rawstr, ",")
		*v = *NewStringList(strlist)
		return nil
	}
	return newError("unknown format of a string list: " + string(data))
}

type Address struct {
	net.Address
}

func (v *Address) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return newError("invalid address: ", string(data)).Base(err)
	}
	if strings.HasPrefix(rawStr, "env:") {
		rawStr = os.Getenv(rawStr[4:])
	}
	v.Address = net.ParseAddress(rawStr)

	return nil
}

func (v *Address) Build() *net.IPOrDomain {
	return net.NewIPOrDomain(v.Address)
}

type Network string

func (v Network) Build() net.Network {
	switch strings.ToLower(string(v)) {
	case "tcp":
		return net.Network_TCP
	case "udp":
		return net.Network_UDP
	case "unix":
		return net.Network_UNIX
	default:
		return net.Network_Unknown
	}
}

type NetworkList []Network

func (v *NetworkList) UnmarshalJSON(data []byte) error {
	var strarray []Network
	if err := json.Unmarshal(data, &strarray); err == nil {
		nl := NetworkList(strarray)
		*v = nl
		return nil
	}

	var rawstr Network
	if err := json.Unmarshal(data, &rawstr); err == nil {
		strlist := strings.Split(string(rawstr), ",")
		nl := make([]Network, len(strlist))
		for idx, network := range strlist {
			nl[idx] = Network(network)
		}
		*v = nl
		return nil
	}
	return newError("unknown format of a string list: " + string(data))
}

func (v *NetworkList) Build() []net.Network {
	if v == nil {
		return []net.Network{net.Network_TCP}
	}

	list := make([]net.Network, 0, len(*v))
	for _, network := range *v {
		list = append(list, network.Build())
	}
	return list
}

func parseIntPort(data []byte) (net.Port, error) {
	var intPort uint32
	err := json.Unmarshal(data, &intPort)
