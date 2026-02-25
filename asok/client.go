package asok

import (
	"encoding/binary"
	"net"
	"time"
)

var (
	clientTimeout = 10 * time.Second
)

type AdminSocketClient struct {
	raddr net.UnixAddr
}

func NewAdminSocketClient(socketPath string) AdminSocketClient {
	c := AdminSocketClient{
		raddr: net.UnixAddr{Name: socketPath, Net: "unix"},
	}
	return c
}

func (c *AdminSocketClient) DoRequest(command string) (string, error) {
	conn, err := net.DialUnix("unix", nil, &c.raddr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_ = conn.SetWriteDeadline(time.Now().Add(clientTimeout))
	if _, err := conn.Write([]byte(cString(command))); err != nil {
		return "", err
	}

	msgSizeRaw := [4]byte{}
	_ = conn.SetReadDeadline(time.Now().Add(clientTimeout))
	if _, err := conn.Read(msgSizeRaw[:]); err != nil {
		return "", err
	}

	msgSize := binary.BigEndian.Uint32(msgSizeRaw[:])
	buf := make([]byte, msgSize)
	n, err := conn.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func (c *AdminSocketClient) Ping() bool {
	r, err := c.DoRequest(`{"prefix":"0"}`)
	return err == nil && len(r) == 1
}

func cString(s string) string {
	return s + "\x00"
}
