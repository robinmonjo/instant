package web_sockets

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
      break
    }

    data := buffer[0:n]

    if !handShaked {
      resp := handshakeResp(data)
      _, err = conn.Write([]byte(resp))
      if err != nil {
        fmt.Println("Failed to write handshake payload")
        break
      }
      handShaked = true
    } else {
      clearTxt := decode(data)
      fmt.Println(clearTxt)

      _, err = conn.Write([]byte(encode(clearTxt)))
      if err != nil {
        fmt.Println("Failed to write response")
        return
      }
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

  //2: add magick string + sha1 hash + base64 enconding
  secretToHash := secret + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
  cypher := sha1.New()
  io.WriteString(cypher, secretToHash)
  hashedSecret := cypher.Sum(nil)
  encoded := base64.StdEncoding.EncodeToString(hashedSecret)

  //3: clean response string
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

func encode (message string) (result []byte) {
  rawBytes := []byte(message)
  var idxData int

  length := byte(len(rawBytes))
  if len(rawBytes) <= 125 { //one byte to store data length
    result = make([]byte, len(rawBytes) + 2)
    result[1] = length
    idxData = 2
  } else if len(rawBytes) >= 126 && len(rawBytes) <= 65535 { //two bytes to store data length
    result = make([]byte, len(rawBytes) + 4)
    result[1] = 126 //extra storage needed
    result[2] = ( length >> 8 ) & 255
    result[3] = ( length      ) & 255
    idxData = 4
  } else {
    result = make([]byte, len(rawBytes) + 10)
    result[1] = 127
    result[2] = ( length >> 56 ) & 255
    result[3] = ( length >> 48 ) & 255
    result[4] = ( length >> 40 ) & 255
    result[5] = ( length >> 32 ) & 255
    result[6] = ( length >> 24 ) & 255
    result[7] = ( length >> 16 ) & 255
    result[8] = ( length >>  8 ) & 255
    result[9] = ( length       ) & 255
    idxData = 10
  }

  result[0] = 129 //only text is supported

  // put raw data at the correct index
  for i, b := range rawBytes {
    result[idxData + i] = b
  }
  return
}


