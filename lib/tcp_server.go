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
 
const BUF_LEN = 1024

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
  buffer := make([]byte, BUF_LEN)
  handShaked := false
  
  for {
    n, err := conn.Read(buffer)
    
    if err != nil {
      fmt.Println("Error reading:", err.Error())
      return
    }

    data := buffer[0:n]

    if !handShaked {
      resp := handshakeResp(data)
      _, err = conn.Write([]byte(resp))
      if err != nil {
        fmt.Println("Failed to write response")
        return
      }
      handShaked = true
    } else {
      clearTxt := decode(data)
      fmt.Println(clearTxt)
    }

    //reset buffer
    for i := 0; i < n; i++ {
      buffer[i] = 0
    }
  } 
  fmt.Println("Connection closed")
}

func handshakeResp (rawBytes []byte) string {
  //1: find secret
  parts := strings.Split(string(rawBytes), "\n")
  secret := ""
  for _, element := range parts {
    if strings.Contains(element, "Sec-WebSocket-Key") {
      secret = strings.Split(element, ":")[1]
      secret = strings.TrimSpace(secret)
      break
    }
  }

  //2: add magick string + sha1 hash
  secretToHash := secret + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
  cypher := sha1.New()
  io.WriteString(cypher, secretToHash)
  hashedSecret := cypher.Sum(nil)

  //3: base64 encode
  encoded := base64.StdEncoding.EncodeToString(hashedSecret)

  //4: clean response string
  response := "HTTP/1.1 101 Switching Protocols\r\n"
  response += "Upgrade: websocket\r\n"
  response += "Connection: Upgrade\r\n"
  response += "Sec-WebSocket-Accept: " + encoded + "\r\n\r\n"
  return response
}

//big thx: http://stackoverflow.com/questions/8125507/how-can-i-send-and-receive-websocket-messages-on-the-server-side
func decode (rawBytes []byte) string {
  var idxMask int
  if rawBytes[1] == 126 {
    idxMask = 4
  } else if rawBytes[1] == 127 {
    idxMask = 10
  } else {
    idxMask = 2
  }

  masks := rawBytes[idxMask:idxMask + 4]
  data := rawBytes[idxMask + 4:len(rawBytes)]
  decoded := make([]byte, len(rawBytes) - idxMask + 4)

  for i, b := range data {
    decoded[i] = b ^ masks[i % 4]
  }
  return string(decoded)
}


