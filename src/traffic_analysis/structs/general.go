package structs

import "time"

type (
	History struct {
		TransactionID   string
		Handshake       Handshake
		TransactionTime time.Time
	}
	Handshake struct {
		Request  TCPPacket
		Response TCPPacket
	}
	Packet interface {
		Unmarshal([]byte)
		LogPrint()
	}
)
