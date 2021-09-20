package netcp

// #cgo LDFLAGS: -lcrypto
// #include <stdio.h>
// #include <openssl/sha.h>
// #include <stdlib.h>
// #include <string.h>
import "C"

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

//Server log
func ServerLog(m string) { //could be async
	fmt.Println("ServerLog:" + m) //[" + time.Now().String() + "]" + "\n:
}

//Network establishment
func CheckAndResolveDialAddress(RemoteAddress string, RemotePort string) *net.TCPConn {

	//check tcp addr and Resolve
	tcpaddr, err := net.ResolveTCPAddr("tcp4", RemoteAddress+":"+RemotePort)
	if err != nil {
		fmt.Println("Error: Remote Adder+Port resolve fail.")
		fmt.Println(RemoteAddress + RemotePort)
		fmt.Println(err)
	}
	tcpcon, err := net.DialTCP("tcp4", nil, tcpaddr)
	if err != nil {
		fmt.Println("Error: Dialign Remote Adder+ Port fail.")
		fmt.Println(RemoteAddress + RemotePort)
		fmt.Println(err)
	}

	return tcpcon
}

func CheckAndListeningOnPort(localAddress string, ListiningPort string) *net.TCPListener {

	//check tcp addr and Resolve
	remoteTcpAdder, err := net.ResolveTCPAddr("tcp4", localAddress+ListiningPort)
	if err != nil {
		fmt.Println("Error: Remote Adder+Port resolve fail.")
		fmt.Println(localAddress + ListiningPort)
		fmt.Println(err)
		return nil
	}
	localTCPListener, err := net.ListenTCP("tcp4", remoteTcpAdder)
	if err != nil {
		fmt.Println("Error: Dialign Remote Adder+ Port fail.")
		fmt.Println(localAddress + ListiningPort)
		fmt.Println(err)
		return nil
	}

	return localTCPListener
}

//tcp io
//RECIVE  Excact targetDataLength byte message
func ReciveConstBytes(op *net.TCPConn, targetDataLength int) (res []byte, err error) {
	res = nil
	for {
		//read from tcp buffer
		tempReadData := make([]byte, targetDataLength)
		readbyte, err := op.Read(tempReadData)
		if err != nil && err != io.EOF {
			fmt.Println("Error: In reciveConstBytes:")
			fmt.Println(err)
			return nil, err
		}

		res = append(res[:], tempReadData[:readbyte]...)
		//check if read exact wanted
		if readbyte == targetDataLength {
			if int(targetDataLength) != len(res) {
				fmt.Println("Error: In reciveConstBytes: didnt read exact!!!")
				fmt.Println("Read : ", readbyte)
			}
			return res, nil
			//if didnt read exact wanted then recalculate targetlength and redo
		} else if readbyte < targetDataLength {
			targetDataLength -= readbyte
		}
	}
}

//RECIVE  [Header]:[Data] with headerSpeace
func ReciveConstHeaderData(op *net.TCPConn, headerSpace int) ([]byte, []byte, error) {
	rawheader, err := ReciveConstBytes(op, headerSpace)
	header := rawheader
	if err != nil {
		fmt.Println("In func ReciveConstHeaderData Error:", err)
		return nil, nil, err
	}
	if headerSpace <= 4 {
		blank := []byte{0}
		for i := 0; i < (4 - int(headerSpace)); i++ {
			header = append(header[:], blank...)
		}
	}
	dataLength := binary.LittleEndian.Uint32(header)
	data, err := ReciveConstBytes(op, int(dataLength))
	if err != nil {
		fmt.Println("In func ReciveConstHeaderData Error:", err)
		fmt.Println("Error: In reciveConstBytes: data")
		fmt.Println(err)
		return nil, nil, err
	}
	return rawheader, data, nil
}
