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
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	__version = "0.0.1"

	DumpUrl uint16 = 1 << (iota - 1)
	DumpReqHeader
	DumpReqParam
	DumpResHeader
	DumpResBody
)

var (
	_timeOut            = 30 * time.Second
	_contentJsonRegex   = "application/json"
	_enableCookie       = false
	_baseUrl              = "http://localhost:8080"	// https://wanna-shop-test.elasticbeanstalk.com:8443
	_params               = make(map[string]string)
	_headers              = make(map[string]string)
	_requestSerialization = "http" // can be http, json, xml

	DumpOption = DumpUrl | DumpReqHeader | DumpReqParam | DumpResBody | DumpResHeader

	_user string
	_pwd  string

	_perty bool

	_lastOutput string

	LastRequestCmd_ []string
)

func main() {

	fmt.Println(Color("Geeko is a CLI http tools with version:", Cyan), __version)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(Color(">>> ", Cyan))

		data, _, _ := reader.ReadLine()
		str := string(data)
		for {
			if strings.HasSuffix(str, "\\") {
				str = str[0 : len(str)-1]
				fmt.Printf(Color("  > ", Cyan))
				d, _, _ := reader.ReadLine()
				str = str + " " + string(d)
			} else {
				break
			}
		}
		fmt.Println("input command:", str)
		inStr := strings.TrimSpace(string(data))

		if inStr == "" {
			continue
		}
		if inStr == "quit" || inStr == "exit" {
			break
		}

		reg := regexp.MustCompile(`[\S]+`)
		cms := reg.FindAllString(inStr, -1)
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
		case "DO":
			result, err = doDoListCommand(cmd, cms)
		case "REMOVE":
			result, err = doRemoveCommand(cmd, cms)
		case "CO", "COPY":
			result, err = doCopyCommand(cmd, cms)
		case "GET", "POST", "PUT", "HEAD":
			result, err = doRequestCommand(cmd, cms)
		case "UPLOAD", "DOWNLOAD":
			result = "Not implemet " + Color(cmd, Red)
		default:
			result = Color("ERROR:", Red) + "\n\t" + "unknow command " + Color(cmd, Yellow) + "\n"
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

func parseInputArgs(s string) (ss []string, complete string) {

	return
}
