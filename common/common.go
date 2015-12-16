package common

import (
	"fmt"
	"strings"
)


type Vector3 struct {
	X float32
	Y float32
	Z float32
}


type Command int

const (
	LOGIN Command = iota
	MOVE
	REGISTER
	SAY
	NOCOMMAND
)

var commands = [...]string{
	"LOGIN",
	"MOVE",
	"REGISTER",
	"SAY",
}

func Log(v ...interface{}) {
    fmt.Println(v...)
}

func IsCommand(text string) (Command, []string ,bool) {
	isCommand := false
	command := NOCOMMAND

	texts := strings.Split(text, " ")

	for index, element := range commands {
		if element == texts[0] {
			//Valid command
			isCommand = true
			command = Command(index)
			break
		}
	}

	return command, texts[1:len(texts)], isCommand
	
}
