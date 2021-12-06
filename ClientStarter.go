package main

import (
	"fmt"
	"net"
	Protocol "startClient/protocol"
	Structs "startClient/structs"
	"time"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	out:
	for ; ; {
		//主动连接服务器
		conn, err :=  net.Dial ("tcp", ": 8848" ) //服务器的ip地址和端口
		if err != nil {
			fmt.Println ( "connection err = " , err)
			continue out
		}
		register(conn)
		Protocol.MessageDecode(conn,messageHandler)
		time.Sleep(time.Second * 10)
	}

}

func messageHandler(message *Structs.Head) {
	switch message.MessageType {
	case Protocol.HEART_BEAT:
		heartBeat(message)
		break
	case Protocol.STARTUP:
		fmt.Println("start ")
		startup()
		break
	case Protocol.REBOOT:
		fmt.Println("reboot ")
		break
	}
}

func heartBeat(message *Structs.Head) {
	conn := *message.Conn
	b := []byte{0x08,0x08,0x04,0x08,Protocol.HEART_BEAT,0x00,0x00,0x00,0x00}
	conn.Write(b)
}

func register(conn net.Conn) {
	b := []byte{0x08,0x08,0x04,0x08,Protocol.REGISTER_STARTER,0x00,0x00,0x00,0x00}
	conn.Write(b)
}


func startup()  {
	err := rpio.Open()
	if err != nil {
		return
	}
	pin := rpio.Pin(10)

	pin.Output()       // Output mode
	pin.High()         // Set pin High
	pin.Low()          // Set pin Low
	pin.Toggle()       // Toggle pin (Low -> High -> Low)

	pin.Input()        // Input mode
	res := pin.Read()  // Read state from pin (High / Low)

	fmt.Println(res)

	pin.Mode(rpio.Output)   // Alternative syntax
	pin.Write(rpio.High)    // Alternative syntax




}


