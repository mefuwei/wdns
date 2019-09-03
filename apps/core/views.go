package core

import "net"

type View struct {
	Name string `json:"name"`
	IPNet []net.IPNet `json:"ip_net"`
}

func NewView(name string, ipnet []net.IPNet) *View {
	vier := &View{
		Name:  name,
		IPNet: ipnet,
	}
	return vier
}

func List() []View {

}