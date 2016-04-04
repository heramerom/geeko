//
// cmd.go
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
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

func doHelpCommand(cmd string, cms []string) (res string, err error) {
	return
}

func doAddCommand(cmd string, cms []string) (result string, err error) {
	var headers, params string
	f := flag.NewFlagSet(cmd, flag.ContinueOnError)
	f.StringVar(&headers, "h", "", "header for each reqeust")
	f.StringVar(&headers, "header", "", "header for each request")
	f.StringVar(&params, "p", "", "params for each request")
	f.StringVar(&params, "param", "", "params for each request")
	f.StringVar(&params, "f", "", "params for each request")
	f.StringVar(&params, "form", "", "params for each request")

	err = f.Parse(cms)
	if err != nil {
		return "", err
	}
	hs, err := parseHeader(headers)
	if err != nil {
		return "", err
	}
	ps, err := parseParams(params)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBufferString("")
	if hs != nil {
		if len(hs) != 0 {
			buf.WriteString(Color("Add Header:\n", Yellow))
		}
		for k, v := range hs {
			Headers[k] = v
			buf.WriteString("\t" + k + " : " + v + "\n")
		}
	}
	if ps != nil {
		if len(ps) != 0 {
			buf.WriteString(Color("Add Params:\n", Yellow))
		}
		for k, v := range ps {
			Params[k] = v
			buf.WriteString("\t" + k + " : " + v + "\n")
		}
	}
	return buf.String(), nil
}

func doSaveCommand(cmd string, cms []string) (result string, err error) {
	var isSchema bool
	f := flag.NewFlagSet(cmd, flag.ContinueOnError)
	f.BoolVar(&isSchema, "schema", false, "save a schema")
	f.BoolVar(&isSchema, "s", false, "save a schema")
	err = f.Parse(cms)
	if err != nil {
		return
	}
	if isSchema {
		schema := NewSchema(f.Arg(0), BaseUrl, User, Pwd, TimeOut, Headers, Params)
		err = SaveSchemes(schema.Name, schema)
		if err != nil {
			return
		}
		result = Color("Save Schema:\n\t", Yellow) + schema.String()
	} else {
		cs := append(cms, LastRequestCmd...)
		item := NewCmdItemWithArgs(cs)
		err = SaveToList(*item)
		if err != nil {
			return
		}
		result = Color("Save:\n\t", Yellow) + strings.Join(cs, " ")
	}
	return
}

// requestType, timeout, baseUrl
func doSetCommand(cmd string, cms []string) (result string, err error) {
	var reqType, baseUrl string
	var timeout int64
	var perty bool
	var user, pwd string
	var enableCookie bool

	var dumpUrl, dumpReqHeader, dumpResHeader, dumpReqParam, dumpResBody bool

	f := flag.NewFlagSet(cmd, flag.ContinueOnError)
	f.StringVar(&reqType, "t", "", "request type, can be http, json, xml, default http.")
	f.StringVar(&reqType, "type", "", "request type, can be http, json, xml, default http.")
	f.StringVar(&baseUrl, "b", "", "the base url")
	f.StringVar(&baseUrl, "base", "", "the base url")
	f.Int64Var(&timeout, "timeout", 0, "request time out")
	f.BoolVar(&perty, "perty", true, "perty output")
	f.StringVar(&user, "u", "", "user name")
	f.StringVar(&user, "user", "", "user name")
	f.StringVar(&pwd, "p", "", "password")
	f.StringVar(&pwd, "pwd", "", "password")
	f.StringVar(&pwd, "password", "", "pwssword")

	f.BoolVar(&enableCookie, "enableCookie", true, "enable cookie")

	f.BoolVar(&dumpUrl, "dumpUrl", true, "dump url")
	f.BoolVar(&dumpReqHeader, "dumpReqHeader", true, "dump req header")
	f.BoolVar(&dumpReqParam, "dumpReqParam", true, "dump req params")
	f.BoolVar(&dumpResHeader, "dumpResHeader", true, "dump res header")
	f.BoolVar(&dumpResBody, "dumpResBody", true, "dump res body")

	err = f.Parse(cms)

	if err != nil {
		return
	}

	buf := bytes.NewBufferString("Set:\n")
	f.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "t", "type":
			Serialization = reqType
			buf.WriteString("\tset Request Serialization with " + Color(reqType, Cyan) + " type success")
		case "b", "base":
			BaseUrl = baseUrl
			buf.WriteString("\tset base url " + baseUrl + " success")
		case "timeout":
			TimeOut = time.Duration(timeout)
		case "u", "user":
			User = user
			buf.WriteString("\tset user " + user + " success")
		case "p", "pwd", "password":
			Pwd = pwd
			buf.WriteString("\tset password " + pwd + " success")
		case "enableCookie":
			EnableCookie = enableCookie
			s := "ON"
			if !enableCookie {
				s = "OFF"
			}
			buf.WriteString("\tset enableCookie " + s)
		case "dumpUrl":
			dump := DumpUrl
			s := "ON"
			if !dumpUrl {
				dump = ^DumpUrl
				s = "OFF"
			}
			DumpOption = DumpOption & dump
			buf.WriteString("\tset DumpUrl " + s)
		case "dumpReqHeader":
			dump := DumpReqHeader
			s := "ON"
			if !dumpReqHeader {
				dump = ^DumpReqHeader
				s = "OFF"
			}
			DumpOption = DumpOption & dump
			buf.WriteString("\tset DumpReqHeader " + s)
		case "dumpReqParam":
			dump := DumpReqParam
			s := "ON"
			if !dumpReqParam {
				dump = ^DumpReqParam
				s = "OFF"
			}
			DumpOption = DumpOption & dump
			buf.WriteString("\tset DumpReqParam " + s)
		case "dumpResHeader":
			dump := DumpResHeader
			s := "ON"
			if !dumpResHeader {
				dump = ^DumpResHeader
				s = "OFF"
			}
			DumpOption = DumpOption & dump
			buf.WriteString("\tset DumpResHeader " + s)
		case "dumpResBody":
			dump := DumpResBody
			s := "ON"
			if !dumpResBody {
				dump = ^DumpResBody
				s = "OFF"
			}
			DumpOption = DumpOption & dump
			buf.WriteString("\tset DumpResBody " + s)
		}
	})
	result = buf.String()
	return
}

func doListCommand(cmd string, cms []string) (result string, err error) {

	var isSchema bool
	f := flag.NewFlagSet(cmd, flag.ContinueOnError)
	f.BoolVar(&isSchema, "schema", false, "save a schema")
	f.BoolVar(&isSchema, "s", false, "save a schema")
	err = f.Parse(cms)
	if err != nil {
		return
	}

	if isSchema {
		schemas, err := SchemaLists()
		if err != nil {
			return "", err
		}
		fmt.Println("schema", schemas)
		for _, v := range schemas {
			result = result + v.String() + "\n"
		}
		return result, err

	} else {
		items, err := CmdItemLists()
		if err != nil {
			return "", err
		}
		if len(cms) == 0 {
			buf := bytes.NewBufferString("")
			for _, v := range items {
				buf.WriteString(Color(v.Name, Cyan))
				buf.WriteString(" ")
				buf.WriteString(Color(v.Cmd, Green))
				buf.WriteString(" ")
				buf.WriteString(Color(strings.Join(v.Args, " "), Green))
				buf.WriteString("\n")
			}
			result = buf.String()
		} else {
			item, err := FindCmdItemsWithName(cms[0])
			if err != nil {
				return "", err
			}
			buf := bytes.NewBufferString("")
			buf.WriteString(Color(item.Name, Cyan))
			buf.WriteString(" ")
			buf.WriteString(Color(item.Cmd, Green))
			buf.WriteString(" ")
			buf.WriteString(Color(strings.Join(item.Args, " "), Green))
			buf.WriteString("\n")
			result = buf.String()
		}
	}
	return
}

func doRequestCommand(method string, params []string) (result string, err error) {
	var header, param string
	f := flag.NewFlagSet(method, flag.ContinueOnError)
	f.StringVar(&header, "h", "", "header for this request")
	f.StringVar(&header, "header", "", "header for this request")
	f.StringVar(&param, "p", "", "params for this request")
	f.StringVar(&param, "param", "", "params for this request")
	f.StringVar(&param, "f", "", "params for this request")
	f.StringVar(&param, "form", "", "params for this request")

	err = f.Parse(params)
	if err != nil {
		return
	}

	urls := f.Args()
	if len(urls) != 1 {
		err = errors.New(`usage: [get|post] -h "key:value" -p "key=value" url `)
		return
	}

	req, err := NewBeegoRequest(method, urls[0], header, param, Serialization)
	if err != nil {
		return
	}
	req.Debug(true)

	// save the request
	LastRequestCmd = make([]string, 5)
	LastRequestCmd = append(LastRequestCmd, method)
	LastRequestCmd = append(LastRequestCmd, params...)

	res, err := req.Response()
	if err != nil {
		return
	}

	var dumpHeader, dumpBody []byte
	dump := req.DumpRequest()
	dps := strings.Split(string(dump), "\n")
	for i, line := range dps {
		println(line)
		if len(strings.Trim(line, "\r\n ")) == 0 {
			dumpHeader = []byte(strings.Join(dps[:i], "\n"))
			dumpBody = []byte(strings.Join(dps[i:], "\n"))
			break
		}
	}

	buf := bytes.NewBufferString("")

	if DumpOption&DumpUrl == DumpUrl {
		buf.WriteString(Color("URL:\n\t", Yellow))
		buf.WriteString(Color(strings.ToUpper(method), Cyan))
		buf.WriteString(" ")
		buf.WriteString(res.Request.URL.String())
		buf.WriteString("\n")
	}

	if DumpOption&DumpReqHeader == DumpReqHeader {
		buf.WriteString(Color("Request Header:\n", Yellow))
		lines := strings.Split(string(dumpHeader), "\n")
		fmt.Println(lines)
		for k, v := range res.Request.Header {
			buf.WriteString("\t")
			buf.WriteString(k)
			buf.WriteString(":")
			buf.WriteString(strings.Join(v, ", "))
			buf.WriteString("\n")
		}
	}

	if DumpOption&DumpReqParam == DumpReqParam {
		buf.WriteString(Color("Params:\n", Yellow))
		buf.WriteString("\t")
		if res.Request.Method == "POST" {
			if err != nil {
				buf.WriteString(err.Error())
			} else {
				buf.WriteString(string(dumpBody))
			}
		} else {
			buf.WriteString("Empty")
		}
		buf.WriteString("\n")
	}

	if DumpOption&DumpResHeader == DumpResHeader {
		buf.WriteString(Color("Response Header:\n", Yellow))
		for k, v := range res.Header {
			buf.WriteString("\t")
			buf.WriteString(k)
			buf.WriteString(" : ")
			buf.WriteString(strings.Join(v, ", "))
			buf.WriteString("\n")
		}
	}

	if DumpOption&DumpResBody == DumpResBody {
		buf.WriteString(Color("Response Body:\n", Yellow))
		buf.WriteString("\t")
		body := formatResponseBody(res, req, true)
		LastOutput = body
		buf.WriteString(ColorfulResponse(body, res.Header.Get("Content-Type")))
		buf.WriteString("\n")
	}
	result = buf.String()
	return
}

func doCopyCommand(cmd string, cms []string) (res string, err error) {
	err = clipboard.WriteAll(LastOutput)
	if err != nil {
		return
	}
	res = "\tsuccess to copy the body"
	return
}

func doUseCommand(cmd string, cms []string) (res string, err error) {
	if len(cms) != 1 {
		err = errors.New("Usage:\n\tuse [schema name]")
		return
	}
	schema := SchemaWithName(cms[0])
	if schema == nil {
		err = errors.New("Error:\n\tcan not find schema with name " + cms[0])
		return
	}
	BaseUrl = schema.Url
	TimeOut = schema.Timeout
	User = schema.User
	Pwd = schema.Pwd
	Headers = schema.Headers
	Params = schema.Params
	res = "Use schema success"
	return
}

func doRemoveCommand(cmd string, cms []string) (res string, err error) {
	if len(cms) == 0 {
		return "", errors.New("Usage:[remove|rm] list-name")
	}
	RemoveItemWithName(cms[0])
	res = "success remove " + Color(cms[0], Cyan)
	return
}

func doStateCommend(cmd string, cms []string) (result string, err error) {

	buf := bytes.NewBufferString(Color("State:\n", Yellow))
	buf.WriteString("\t")
	buf.WriteString("Base Url: ")
	buf.WriteString(Color(BaseUrl, Cyan))
	buf.WriteString("\n")

	buf.WriteString("\t")
	buf.WriteString("Request Serialization: ")
	buf.WriteString(Color(Serialization, Cyan))
	buf.WriteString("\n")

	buf.WriteString("\t")
	buf.WriteString("Time Out: ")
	buf.WriteString(Color(TimeOut.String(), Cyan))
	buf.WriteString("\n")

	if len(Headers) != 0 {
		buf.WriteString("\tHeaders:\n")
		for k, v := range Headers {
			buf.WriteString("\t")
			buf.WriteString(k)
			buf.WriteString(" : ")
			buf.WriteString(v)
		}
		buf.WriteString("\n")
	}

	if len(Params) != 0 {
		buf.WriteString("\tForms:\n")
		for k, v := range Params {
			buf.WriteString("\t")
			buf.WriteString(k)
			buf.WriteString(" : ")
			buf.WriteString(v)
		}
		buf.WriteString("\n")
	}

	result = buf.String()

	return
}
