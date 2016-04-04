//
// geeko.go
//
// Create at 16/3/14
//
// Copyright (C) 2016 xhl <heramerom@163.com>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

const (
	Version = "0.0.1"

	DumpUrl uint16 = 1 << (iota - 1)
	DumpReqHeader
	DumpReqParam
	DumpResHeader
	DumpResBody
)

var cmdCompleter = readline.NewPrefixCompleter(
	readline.PcItem("st"),
	readline.PcItem("state"),
	readline.PcItem("set",
		readline.PcItem("-b"),
		readline.PcItem("-baseUrl"),
		readline.PcItem("-t", readline.PcItem("http"), readline.PcItem("json"), readline.PcItem("form")),
		readline.PcItem("-type", readline.PcItem("http"), readline.PcItem("json"), readline.PcItem("form")),
		readline.PcItem("-timeout"),
		readline.PcItem("-enableCookie"),
		readline.PcItem("-dumpReqParam", readline.PcItem("true"), readline.PcItem("false")),
		readline.PcItem("-dumpReqHeader", readline.PcItem("true"), readline.PcItem("false")),
		readline.PcItem("-dumpResHeader", readline.PcItem("true"), readline.PcItem("false")),
		readline.PcItem("-dumpResBody", readline.PcItem("true"), readline.PcItem("false"))),
	readline.PcItem("save", readline.PcItem("-s=true")),
	readline.PcItem("ls", readline.PcItem("-s=true")),
	readline.PcItem("get"),
	readline.PcItem("post"))

var (
	TimeOut       = 30 * time.Second
	BaseUrl       = "http://localhost:8080"
	Params        = make(map[string]string)
	Headers       = make(map[string]string)
	Serialization = "form" // can be from, http, json, default form
	EnableCookie  = true

	DumpOption = DumpUrl | DumpReqHeader | DumpReqParam | DumpResBody | DumpResHeader

	User string
	Pwd  string

	LastOutput string

	LastRequestCmd []string
)

func main() {

	fmt.Println(Color("Geeko is a cURL-like http tools with version:", Yellow), Version)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       Color(">>> ", Cyan),
		HistoryFile:  "/tmp/geeko.history",
		AutoComplete: NewListComplete(cmdCompleter),
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	buf := bytes.NewBufferString("")
	for {

		line, e := rl.Readline()
		if err != nil {
			fmt.Println(Color("Error:\n", Red), "\t", e.Error())
			continue
		}

		app, complete := AppendInput(buf.String(), line)
		buf.WriteString(app)
		if !complete {
			rl.SetPrompt(Color("--> ", Cyan))
			continue
		}

		ins := strings.TrimSpace(buf.String())
		rl.SetPrompt(Color(">>> ", Cyan))
		buf.Reset()

		if ins == "" {
			continue
		}
		if ins == "quit" || ins == "exit" {
			break
		}

		cms, e := ParseInputArgs(ins)
		if e != nil {
			fmt.Println(Color("Error:\n", Red), "\t", e.Error())
		}

		cmd := cms[0]
		cms = cms[1:]

		var result string
		var err error

		switch strings.ToUpper(cmd) {
		case "STATE", "ST":
			result, err = doStateCommend(cmd, cms)
		case "ADD":
			result, err = doAddCommand(cmd, cms)
		case "SET":
			result, err = doSetCommand(cmd, cms)
		case "SAVE", "S":
			result, err = doSaveCommand(cmd, cms)
		case "LIST", "LS", "L":
			result, err = doListCommand(cmd, cms)
		case "HELP", "H":
			result, err = doHelpCommand(cmd, cms)
		case "USE":
			result, err = doUseCommand(cmd, cms)
		case "DO":
			result, err = doListCommand(cmd, cms)
		case "REMOVE":
			result, err = doRemoveCommand(cmd, cms)
		case "CO", "COPY":
			result, err = doCopyCommand(cmd, cms)
		case "GET", "POST", "PUT", "HEAD":
			result, err = doRequestCommand(cmd, cms)
		case "UPLOAD", "DOWNLOAD":
			result = "Not implemet " + Color(cmd, Red)
		default:
			err = errors.New("unknow command" + Color(cmd, Yellow))
		}
		if err != nil {
			fmt.Println(Color("ERROR:\n\t", Red) + err.Error() + "\n")
			continue
		}
		if len(result) != 0 {
			fmt.Println(result)
		}
	}
}

func AppendInput(org, app string) (args string, complete bool) {
	var s string
	s = " " + app
	complete = true
	if strings.HasSuffix(s, "\\") && !strings.HasSuffix(s, "\\\\") {
		complete = false
		s = s[:len(s)-1]

	}
	tmp := org + s
	c := strings.Count(tmp, `"`) - strings.Count(tmp, `\"`)
	if c%2 != 0 {
		complete = false
	}
	args = s
	return
}

func ParseInputArgs(line string) (args []string, err error) {
	as := make([]string, 0)
	rs := []rune(line)
	res := make([][]rune, 0)
	sub := make([]rune, 0)
	quote := false
	tm := false
	for _, v := range rs {
		if v == '\\' {
			if tm {
				sub = append(sub, v)
				tm = false
			} else {
				tm = true
			}
		} else if v == '"' {
			if tm {
				sub = append(sub, v)
				tm = false
			} else {
				quote = !quote
				if !quote {
					res = append(res, sub)
					sub = make([]rune, 0)
				}
			}
		} else if v == ' ' {
			if tm {
				sub = append(sub, v)
				tm = false
			} else if quote {
				sub = append(sub, v)
			} else {
				if len(sub) != 0 {
					res = append(res, sub)
					sub = make([]rune, 0)
				}
			}
		} else {
			if quote {
				if tm {
					if v == 't' {
						sub = append(sub, []rune{' ', ' ', ' ', ' '}...)
					} else if v == 'n' {
						sub = append(sub, '\n')
					} else {
						sub = append(sub, v)
					}
					tm = false
				} else {
					sub = append(sub, v)
				}
			} else {
				sub = append(sub, v)
			}
		}
	}
	if len(sub) != 0 {
		res = append(res, sub)
	}

	for _, v := range res {
		as = append(as, string(v))
	}
	args = as
	return
}
