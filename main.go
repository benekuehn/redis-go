package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Message struct {
	cmd         string
	arrayLength int
	args        []byte
}

func handleCmd(m Message) []byte {
	cmd := strings.ToUpper(m.cmd)
	switch cmd {
	case "PING":
		return []byte("+PONG\r\n")
	case "COMMAND":
		supportedCmds := "*1\r\n" +
			"$4\r\n" +
			"PING\r\n"
		return []byte(supportedCmds)
	default:
		return []byte("-ERR unknown command '" + cmd + "\r\n")
	}
}

func handleLine(line []byte, m Message, entriesProcessed int) (Message, int, error) {
	if line[0] == '*' && m.arrayLength == 0 {
		lengthStr := string(line[1:])
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return m, entriesProcessed, fmt.Errorf("error converting length: %w", err)
		}
		m.arrayLength = length
	}
	if entriesProcessed == 0 {
		// first element is the command
		if line[0] == '$' || line[0] == '*' {
			// currently no need to handle this
		} else {
			m.cmd = string(line)
			fmt.Printf("Command: %s\n", m.cmd)
			entriesProcessed++
		}
	} else {
		// handle arguments
		// currently no need to handle this
	}

	return m, entriesProcessed, nil

}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	reader := bufio.NewReader(c)

	m := Message{
		cmd:         "",
		arrayLength: 0,
		args:        nil,
	}

	// counter for the outermost array items
	entriesProcessed := 0

	defer c.Close()

	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			println("END OF FILE")
			break
		}

		if isPrefix {
			// true when line is longer than buffer
			// handle later
		}

		m, entriesProcessed, err = handleLine(line, m, entriesProcessed)
		if err != nil {
			fmt.Println("Error handling line:", err)
		}

		if entriesProcessed == m.arrayLength {
			c.Write(handleCmd(m))
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
