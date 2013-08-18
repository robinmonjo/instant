package tcp_server

import (
  "net"
  "os"
  "fmt"
  "strings"
  "crypto/sha1"
  "io"
  "encoding/base64"
)
 
const BUF_LEN = 2048

func StartServer() {
	listener, err := net.Listen("tcp", "0.0.0.0:1337")
	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("TCP server listening on 1337")
 
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accept:", err.Error())
			return
		}
		go handleConnection(conn)
	}
}
 
func handleConnection(conn net.Conn) {
  fmt.Println("Incoming Connection, reading data")

  buffer := make([]byte, BUF_LEN)

  handShaked := false

  for {
    buffer = make([]byte, BUF_LEN)
    _, err := conn.Read(buffer)
    if err != nil {
      fmt.Println("Error reading:", err.Error())
      return
    }
    fmt.Println("received : ", buffer)
    if !handShaked {
      parts := strings.Split(string(buffer), "\n")
      secret := ""
      for _, element := range parts {
        if strings.Contains(element, "Sec-WebSocket-Key") {
          secret = strings.Split(element, ":")[1]
          secret = strings.TrimSpace(secret)
          break
        }
      }
      fmt.Println("Secret is :", secret)
      secretToHash := secret + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
      fmt.Println("Secret to hash is :", secretToHash)
      h := sha1.New()
      io.WriteString(h, secretToHash)
      hashedSecret := h.Sum(nil)
      fmt.Println("hashed secret is :", hashedSecret)
      str := base64.StdEncoding.EncodeToString(hashedSecret)

      response := "HTTP/1.1 101 Switching Protocols\r\n"
      response += "Upgrade: websocket\r\n"
      response += "Connection: Upgrade\r\n"
      response += "Sec-WebSocket-Accept: " + str + "\r\n\r\n"
      fmt.Println(response)
      _, err = conn.Write([]byte(response))
      if err != nil {
        fmt.Println("Failed to write response")
      }
      conn.Write([]byte("\n"))
      handShaked = true
    } else {
      decode(buffer)
    }

  } 
  fmt.Println("Connection closed")
}

func decode (rawBytes []byte) string {
  secondByte := rawBytes[1]

  length := secondByte & 127 // may not be the actual length in the two special cases

  indexFirstMask := 2          // if not a special case

  if length == 126 {
    indexFirstMask = 4
  } else if length == 127 {
    indexFirstMask = 10
  }

  masks := make([]byte, 4)
  masks[0] = rawBytes[indexFirstMask]
  masks[1] = rawBytes[indexFirstMask + 1]
  masks[2] = rawBytes[indexFirstMask + 2]
  masks[3] = rawBytes[indexFirstMask + 3]
 
  indexFirstDataByte := indexFirstMask + 4 // four bytes further

  // decoded = new array

  // length := bytes.length - indexFirstDataByte // length of real data
j := 0
for i := indexFirstDataByte; i < len(rawBytes); i++ {
    fmt.Println (string(rawBytes[i] | masks[j % 4]))
    j++
}
return "hello"

// keep going here : http://stackoverflow.com/questions/8125507/how-can-i-send-and-receive-websocket-messages-on-the-server-side

// now use "decoded" to interpret the received data
}


