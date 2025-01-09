package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	packet := make([]byte, 0)
	tmp := make([]byte, 4096)
	defer c.Close()

	for {
		n, err := c.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			println("END OF FILE")
			break
		}

		packet = append(packet, tmp[:n]...)

		command := string(tmp[:n])
		fmt.Printf("Received: %s\n", command)

		if strings.Contains(command, "COMMAND") {
			response := "*1\r\n" + // Array of 1 element
				"$4\r\n" + // String length 4
				"PING\r\n" // The command name
			c.Write([]byte(response))
		} else if strings.Contains(command, "PING") {
			response := "+PONG\r\n"
			c.Write([]byte(response))
		} else {
			response := "-ERR unknown command\r\n"
			c.Write([]byte(response))
		}
	}

}

func main() {
	l, err := net.Listen("tcp4", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
