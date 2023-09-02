package dialer

import (
	"context"
	"fmt"
	"net"
	"net/netip"

	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

const MTU = 1420

type Dialer struct {
	net *netstack.Net
}

func NewDialer(in Interface, peers ...Peer) (*Dialer, error) {
	addr, err := netip.ParseAddr(in.Address)
	if err != nil {
		return nil, err
	}
	dns_addresses := []netip.Addr{}
	for _, raw := range in.Dns {
		dns, err := netip.ParseAddr(raw)
		if err != nil {
			return nil, err
		}
		dns_addresses = append(dns_addresses, dns)
	}
	tun, tnet, err := netstack.CreateNetTUN([]netip.Addr{addr}, dns_addresses, MTU)
	if err != nil {
		return nil, err
	}

	dev := device.NewDevice(tun, conn.NewDefaultBind(), device.NewLogger(device.LogLevelVerbose, ""))

	key, err := base64KeyToHex(in.PrivateKey)
	if err != nil {
		return nil, err
	}
	ipcString := fmt.Sprintf("private_key=%s\n", key)

	for _, peer := range peers {
		str, err := peer.toIpcString()
		if err != nil {
			return nil, err
		}
		ipcString = fmt.Sprintf("%s%s", ipcString, str)
	}

	logrus.Debugf("ipc_string: %s", ipcString)

	err = dev.IpcSet(ipcString)
	if err != nil {
		return nil, err
	}

	return &Dialer{
		net: tnet,
	}, nil
}

func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return d.net.DialContext(ctx, network, address)
}
