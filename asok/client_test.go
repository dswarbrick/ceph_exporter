package asok

import (
	"fmt"
	"math/rand/v2"
	"net"
	"testing"
	"time"
)

func TestAdminSocketPing(t *testing.T) {
	socketPath := fmt.Sprintf("ceph-osd.%d.asok", rand.Uint32())

	// Mock admin socket listener.
	go func() {
		listener, err := net.ListenUnix("unix", &net.UnixAddr{Net: "unix", Name: socketPath})
		if err != nil {
			panic(err)
		}
		defer listener.Close()

		conn, err := listener.AcceptUnix()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		// Read (and ignore) command.
		cmd := make([]byte, 1024)
		if _, err := conn.Read(cmd); err != nil {
			panic(err)
		}

		// Blindly respond with admin socket "pong".
		if _, err := conn.Write([]byte{0x00, 0x00, 0x00, 0x01, '2'}); err != nil {
			panic(err)
		}
	}()

	// Wait for socket listener to settle.
	time.Sleep(50 * time.Millisecond)

	client := NewAdminSocketClient(socketPath)
	if resp := client.Ping(); !resp {
		t.Fatal("Admin socket ping failed")
	}
}
