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
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"
	"errors"
)

var ListItemMap_ map[string]*ListItem

func init() {
	ListItemMap_ = make(map[string]*ListItem)
	ListItems()
}

type ListItem struct {
	Name string
	Cmd  string
	Args []string
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
	err = os.Mkdir(p, 0777)
	if err != nil {
		return
	}
	return p, nil
}

func GeekoWorkFilePath() (path string, err error) {
	tmp, err := GeekoWorkPath()
	if err != nil {
		return
	}
	if runtime.GOOS == "windows" {
		path = tmp + "\\request.list"
	} else {
		path = tmp + "/request.list"
	}
	return
}

func NewListItem(name string, cmd string, args []string) *ListItem {
	return &ListItem{
		Name: name,
		Cmd:  cmd,
		Args: args,
	}
}

func NewListItemWithLine(line string) *ListItem {
	if line == "" {
		return nil
	}
	args := strings.Split(line, " ")
	return NewListItemWithArgs(args)
}

func NewListItemWithArgs(args []string) *ListItem {
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
	return &ListItem{
		Name: name,
		Cmd:  cmd,
		Args: args,
	}
}

func (this *ListItem) String() string {
	str := this.Name + " " + this.Cmd + " " + strings.Join(this.Args, " ")
	reg, err := regexp.Compile("[ ]+")
	if err != nil {
		return str
	}
	str = reg.ReplaceAllString(str, " ")
	return str
}

func ListLines() (line []string, err error) {
	path, err := GeekoWorkFilePath()
	if err != nil {
		return
	}

	bs, err := readFile(path)
	if err != nil {
		return
	}

	line = strings.Split(string(bs), "\n")
	return
}

func ListItems() (lm map[string]*ListItem, err error) {
	lines, err := ListLines()
	if err != nil {
		return
	}
	for _, v := range lines {
		item := NewListItemWithLine(v)
		if item == nil {
			continue
		}
		ListItemMap_[item.Name] = item
	}
	lm = ListItemMap_
	return
}

func SaveToList(item ListItem) error {
	ListItemMap_[item.Name] = &item
	return save()
}

func RemoveItemWithName(name string) error {
	if len(name) == 0 {
		return errors.New("item name can not be empty")
	}
	delete(ListItemMap_, name)
	return save()
}

func save() error{
	path, err := GeekoWorkFilePath()
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString("")
	for _, v := range ListItemMap_ {
		buf.WriteString(v.String() + "\n")
	}
	err = ioutil.WriteFile(path, buf.Bytes(), 0)
	return err
}

func ListCommands() (ls []string, err error) {
	path, err := GeekoWorkFilePath()
	if err != nil {
		return
	}
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	ls = strings.Split(string(bs), "\n")
	return
}

func readFile(filename string) ([]byte, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var n int64
	if fi, err := f.Stat(); err == nil {
		if size := fi.Size(); size < 1e9 {
			n = size
		}
	}
	return readAll(f, n+bytes.MinRead)
}

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
