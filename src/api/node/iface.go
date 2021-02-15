package chord

import (
	"errors"
	"net"
)

func getIfaces() (map[string]string, error) {
	ifaces , err := net.Interfaces()
	if err != nil{
		return nil, err
	}

	ifaceMap := make(map[string]string)
	var ip net.IP
	for _, iface := range ifaces{
		addrs, err := iface.Addrs()

		if err != nil{
			return nil, err
		}

		for _, addr := range addrs{
			switch v:= addr.(type){
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}


			if ip.To4() != nil{
				ifaceMap[iface.Name]= ip.String()
			}
		}
	}


	return ifaceMap, nil
}


func IfaceIPv4Addr(ifaceName string)(*string, error){
	ifaceMap, err := getIfaces()

	if err != nil{
		return nil, err
	}

	ipv4Addr, ok := ifaceMap[ifaceName]

	if !ok{
		return nil, errors.New("Interface Name not found")
	}

	return &ipv4Addr, nil
}
