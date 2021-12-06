package Protocol

import (
	"fmt"
	"net"
	Structs "startClient/structs"
)

const MAGIC_NUMBER = 0x08080408

const (
	HEART_BEAT = 0x01 + iota
	REGISTER_STARTER
	REGISTER_MONITOR
	STARTUP
	REBOOT
)

func MessageDecode(conn net.Conn,service func(message *Structs.Head)) *Structs.Head {
	defer conn.Close()
	decode, err := GetLengthFieldDecode(5, 4, 1024)
	if err != nil {
		return nil
	}
	for ; ; {
		parse, err := decode.Parse(&conn)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		magicNumberB := parse[0:4]
		magicNumber, err := bytesToInt32BE(magicNumberB)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if magicNumber != MAGIC_NUMBER {
			fmt.Println("magicNumber error")
			return nil
		}
		messageType := parse[4]
		lengthB := parse[5:9]
		length, err := bytesToInt32BE(lengthB)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		head := new(Structs.Head)
		head.Length = length
		head.MagicNumber = magicNumber
		head.MessageType = messageType
		head.Conn = &conn
		service(head)
	}
}




