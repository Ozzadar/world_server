package client

import (
 "github.com/ozzadar/world_server/common"
 "github.com/ozzadar/world_server/database"
 "container/list"
 "net"
 "bytes"
 "strings"
)

type Client struct {
    Name string
    Authenticated bool
    Position common.Vector3
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

func (c *Client) ExecuteCommand(command common.Command, arguments []string) {


	// CUT OUT CRAP AT END OF BUFFER
	count := 0
	for _, element := range arguments[len(arguments)-1] {
		if element == 0x00 {
			break
		}
		count++
	}

	arguments[len(arguments)-1] = (arguments[len(arguments)-1])[0:count]

	//COMMAND LOGIC
	switch (command) {

		//Login to server
		case common.LOGIN:
			//Arguments username, password
			if len(arguments) != 2 {
				common.Log("Invalid arguments")
				c.Conn.Write([]byte("Invalid command."))
				return
			}

			//Authenticate User
			loggedin := database.Login(arguments[0], arguments[1])

			//Let everyone know user joined
			if loggedin {
				c.Name = arguments[0]

				send := c.Name + " has just joined the world."

				c.Outgoing <- send
			}

		case common.SAY:
			if c.Name != "Unauthenticated" {
				send := "<"+ c.Name + "> "

				for _, element := range arguments {
					send = send + element + " "
				}

				c.Outgoing <- send
			} else {
				c.Conn.Write([]byte("You must be logged in to do this."))
			}

		case common.MOVE:
			send := c.Name + " is moving in world with arguments: "

			for _, element := range arguments {
				send = send + element + ", "
			}

			common.Log(send)
			c.Outgoing <- send

		//Register new user
		case common.REGISTER:
			if len(arguments) != 2 {
				common.Log("Invalid arguments")
				c.Conn.Write([]byte("Invalid command."))
				return
			}

			newuser := database.User{
				Username: 	strings.TrimSpace(arguments[0]),
				Password: 	strings.TrimSpace(arguments[1]),
				Role: 		"admin",
				}

			database.RegisterUser(&newuser)

			common.Log("User " + arguments[0] + " Created")
			c.Conn.Write([]byte("User created. You may now Log In."))
	}
}

