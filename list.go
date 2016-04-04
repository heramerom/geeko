//
// list.go
//
// Create at 16/3/14
//
// Copyright (C) 2016 xhl <heramerom@163.com>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"encoding/gob"
	"errors"
	"io"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"
)

var CmdItems = make(map[string]*CmdItem)
var Schemas = make(map[string]*Schema)

func init() {

	items, _ := CmdItemLists()
	for k, v := range items {
		CmdItems[k] = v
	}

	schemas, _ := SchemaLists()
	for k, v := range schemas {
		Schemas[k] = v
	}

}

type CmdItem struct {
	Name string
	Cmd  string
	Args []string
}

type Schema struct {
	Name    string
	Url     string
	User    string
	Pwd     string
	Timeout time.Duration
	Headers map[string]string
	Params  map[string]string
}

func UserHomePath() (path string, err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	path = u.HomeDir
	return
}

func CheckFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GeekoWorkPath() (path string, err error) {
	home, err := UserHomePath()
	if err != nil {
		return
	}
	var p string
	if runtime.GOOS == "windows" {
		p = home + "\\.geeko"
	} else {
		p = home + "/.geeko"
	}
	if CheckFileExist(p) {
		return p, nil
	}
	err = os.Mkdir(p, 0666)
	if err != nil {
		return
	}
	return p, nil
}

func GeekoWorkReqFilPath() (path string, err error) {
	tmp, err := GeekoWorkPath()
	if err != nil {
		return
	}
	if runtime.GOOS == "windows" {
		path = tmp + "\\req.gob"
	} else {
		path = tmp + "/req.gob"
	}
	return
}

func GeekoWorkSchemaFilePath() (path string, err error) {
	tmp, err := GeekoWorkPath()
	if err != nil {
		return
	}
	if runtime.GOOS == "windows" {
		path = tmp + "\\schema.gob"
	} else {
		path = tmp + "/schema.gob"
	}
	return
}

func NewCmdItem(name string, cmd string, args []string) *CmdItem {
	return &CmdItem{
		Name: name,
		Cmd:  cmd,
		Args: args,
	}
}

func NewSchema(name string, url string, user string, pwd string, timeout time.Duration, headers map[string]string, params map[string]string) *Schema {
	return &Schema{
		Name:    name,
		Url:     url,
		User:    user,
		Pwd:     pwd,
		Timeout: timeout,
		Headers: headers,
		Params:  params,
	}
}

func NewCmdItemWithArgs(args []string) *CmdItem {
	if len(args) < 2 {
		return nil
	}
	name := args[0]
	cmd := args[1]
	if len(args) > 2 {
		args = args[2:]
	} else {
		args = nil
	}
	return &CmdItem{
		Name: name,
		Cmd:  cmd,
		Args: args,
	}
}

func (this *Schema) String() string {
	return this.Name + " " + this.Url
}

func (this *CmdItem) String() string {
	return this.Name + " " + this.Cmd + " " + strings.Join(this.Args, " ")
}

func CmdItemLists() (lm map[string]*CmdItem, err error) {
	path, err := GeekoWorkReqFilPath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	enc := gob.NewDecoder(file)
	err = enc.Decode(&lm)
	return
}

func SchemaLists() (lists map[string]*Schema, err error) {
	path, err := GeekoWorkSchemaFilePath()
	if err != nil {
		return
	}
	file, err := os.Open(path)
	if err != nil {
		return
	}
	var sm map[string]*Schema
	dec := gob.NewDecoder(file)

	e := dec.Decode(&sm)
	if e != nil && e != io.EOF {
		err = e
		return
	}
	lists = sm
	return
}

func FindCmdItemsWithName(name string) (item *CmdItem, err error) {
	if len(name) == 0 {
		err = errors.New("item name can not be nil!")
		return
	}
	for k, v := range CmdItems {
		if k == name {
			item = v
			return
		}
	}
	err = errors.New("can not found list item with name: " + Color(name, Magenta))
	return
}

func SaveSchemes(name string, s *Schema) error {
	Schemas[name] = s
	return saveScheme(Schemas)
}

func SaveToList(item CmdItem) error {
	CmdItems[item.Name] = &item
	return saveCmdItems(CmdItems)
}

func RemoveItemWithName(name string) error {
	if len(name) == 0 {
		return errors.New("item name can not be empty")
	}
	delete(CmdItems, name)
	return saveCmdItems(CmdItems)
}

func SchemaWithName(name string) *Schema {
	return Schemas[name];
}

func saveCmdItems(items map[string]*CmdItem) error {
	path, err := GeekoWorkReqFilPath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(items)
	return err
}

func saveScheme(schemcs map[string]*Schema) error {
	path, err := GeekoWorkSchemaFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(schemcs)
	return err
}
