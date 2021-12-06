package Structs

import "net"

type Head struct {
	MagicNumber int32
	MessageType byte
	Length      int32
	Conn *net.Conn
}
