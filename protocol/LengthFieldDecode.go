package Protocol

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
)


type LengthFieldDecode struct {
	offset    int64
	length    int64
	maxLength int64
}

func GetLengthFieldDecode(offset int64,length int64,maxLength int64) (*LengthFieldDecode,error){
	if offset < 0 || offset > 256 {
		return nil,errors.New("illegal offset !")
	}
	if length < 0 || length > 8 {
		return nil,errors.New("illegal length !")
	}
	if maxLength < 0 {
		return nil,errors.New("illegal maxLength !")
	}
	decode := new(LengthFieldDecode)
	decode.offset = offset
	decode.length = length
	decode.maxLength = maxLength
	return decode,nil
}


func readFully(conn net.Conn) ([]byte, error) {
	// 读取所有响应数据后主动关闭连接

	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result.Bytes(), nil
}


func (lengthFieldDecode *LengthFieldDecode) Parse(conn *net.Conn) ([]byte ,error){
	messageByte, err := readBytesBlocked(conn,lengthFieldDecode.offset)
	if err != nil {
		fmt.Print(err)
		return nil,err
	}
	lengthByte, err := readBytesBlocked(conn,lengthFieldDecode.length)
	if err != nil {
		fmt.Print(err)
		return nil,err
	}

	length, err := bytesToInt64BE(append(make([]byte, 8 - lengthFieldDecode.length), lengthByte...))
	if err != nil {
		fmt.Print(err)
		return nil,err
	}
	tailByte, err := readBytesBlocked(conn,length)
	if err != nil {
		fmt.Print(err)
		return nil,err
	}
	resultBytes := append(messageByte,lengthByte...)
	resultBytes = append(resultBytes,tailByte...)
	return resultBytes,nil
}

func readBytesBlocked(conn *net.Conn,length int64) ([]byte,error) {
	conns := *conn
	result := make([]byte, 0)
	for ;; {
		if length <= 0 {
			return result,nil
		}
		tmp := make([]byte, length)
		readLength, err := conns.Read(tmp)
		if err != nil {
			if err == io.EOF {
				fmt.Println("closed connection ")
				result = append(result, result[0:readLength]...)
				return result,err
			}
			fmt.Println("read error ")
			fmt.Println(err)
			return nil, err
		}
		result = append(result, tmp[0:readLength]...)
		length -= int64(readLength)
	}
}
