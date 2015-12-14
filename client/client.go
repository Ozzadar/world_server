package client

import (
 "github.com/ozzadar/world_server/common"
 "container/list"
 "net"
 "bytes"
)

type Client struct {
    Name string
    Incoming chan string
    Outgoing chan string
    Conn net.Conn
    Quit chan bool
    ClientList *list.List
}

func (c *Client) Read(buffer []byte) bool {
	bytesRead, error := c.Conn.Read(buffer)

	if error != nil {
		c.Close()
		common.Log(error)
		return false
	}

	common.Log("Read ", bytesRead, " bytes")
	return true
}

func (c *Client) Close() {
	c.Quit <- true
	c.Conn.Close()
	c.RemoveMe()
}

func (c *Client) Equal(other *Client) bool {
	if bytes.Equal([]byte(c.Name), []byte(other.Name)) {

		if c.Conn == other.Conn {
			return true
		}
	}
	return false
}

func (c *Client) RemoveMe() {
	for entry := c.ClientList.Front(); entry != nil; entry = entry.Next() {
		client := entry.Value.(Client)

		if c.Equal(&client) {
			common.Log("RemoveMe: ", c.Name)
			c.ClientList.Remove(entry)
		}
	}
}

