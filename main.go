//server.go
/*

  FIRST COMMIT IS QUITE LITERALLY A COPY-PASTE OF TCP CODE FOUND ONLINE

*/
package main

import (  
    "net"
    "github.com/ozzadar/world_server/client"
    "container/list"
    "bytes"
    "github.com/ozzadar/world_server/common"
)

const (  
    CONN_HOST = ""
    CONN_PORT = "1337"
    CONN_TYPE = "tcp"
)



func main() {  
  common.Log ("Hello Server!")

  clientList := list.New()
  in := make(chan string)

  go IOHandler(in, clientList)

  service := ":1337"
  tcpAddr, error := net.ResolveTCPAddr("tcp", service)

  if error != nil {
      common.Log(error)
  } else {
    netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String())

    if error != nil {
      common.Log(error)
    } else {
      defer netListen.Close()

      for {
        common.Log("Waiting for clients")

        connection, error := netListen.Accept()

        if error != nil {
          common.Log(error)
        } else {
          go ClientHandler(connection, in, clientList)
        }
      }
    }
  }
}

func IOHandler(Incoming <-chan string, clientList *list.List) {
  for {
    common.Log("IOHandler: Waiting for input")
    input := <-Incoming
    common.Log("IOHandler: Handling ", input)

    for e := clientList.Front(); e != nil; e = e.Next() {
      cli := e.Value.(client.Client)
      cli.Incoming <-input
    }
  }
}

func ClientReader(cli *client.Client) {
  buffer := make([]byte,2048)

  for cli.Read(buffer) {
    if bytes.Equal(buffer, []byte("/quit")) {
      cli.Close()
      break
    }

    common.Log("ClientReader received ", cli.Name, "> ", string(buffer))
    send := cli.Name + "> "+ string(buffer)
    cli.Outgoing <- send

    for i := 0; i < 2048; i++ {
      buffer[i] = 0x00
    }
  }

  cli.Outgoing <- cli.Name + " has left chat"
  common.Log("ClientReader stopped for ", cli.Name)

}

func ClientSender(cli *client.Client) {
  for {
    select {
      case buffer := <-cli.Incoming:
        common.Log("ClientSender sending \" ", string(buffer), " \" to ", cli.Name)
        count := 0

        for i:=0; i<len(buffer); i++ {
          if buffer[i] == 0x00 {
            break
          }
          count++
        }
        common.Log("Send size: ", count)
        cli.Conn.Write([]byte(buffer)[0:count])
      case <-cli.Quit:
        common.Log("Client ", cli.Name, " quitting")
        cli.Conn.Close()
        break
    }
  }
}

func ClientHandler(conn net.Conn, ch chan string, clientList *list.List) {
  buffer := make([]byte, 1024)
  bytesRead, error := conn.Read(buffer)

  if error != nil {
    common.Log("Client connection error: ", error)
  }

  name := string(buffer[0:bytesRead])
  newClient := &client.Client{name, make(chan string), ch, conn, make(chan bool), clientList}

  go ClientSender(newClient)
  go ClientReader(newClient)

  clientList.PushBack(*newClient)

  ch <-string(name + " has joined the chat")
}

