//
// parse.go
//
// Create at 16/3/14
//
// Copyright (C) 2016 xhl <heramerom@163.com>
//
// Distributed under terms of the MIT license.
//

package main


import (
	"errors"
	"strings"
)

// "key:value \n  key:value"
func parseHeader(s string) (headers map[string]string, err error) {
	if len(s) == 0 {
		return
	}
	ps := strings.Split(s, "\n")
	hs := make(map[string]string)
	for _, v := range ps {
		kvs := strings.Split(v, ":")
		if len(kvs) != 2 {
			err = errors.New("error to parse headers")
			return
		}
		key := strings.TrimSpace(kvs[0])
		value := strings.TrimSpace(kvs[1])
		if len(key) == 0 || len(value) == 0 {
			err = errors.New("error to parse headers")
			return
		}
		hs[key] = value
	}
	headers = hs
	return
}

// "key=value&key=value"
func parseParams(s string) (params map[string]string, err error) {
	if len(s) == 0 {
		return
	}
	ps := strings.Split(s, "&")
	pars := make(map[string]string)
	for _, v := range ps {
		kvs := strings.Split(v, "=")
		if len(kvs) != 2 {
			err = errors.New("error to parse params")
			return
		}
		key := kvs[0]
		value := kvs[1]
		if len(key) == 0 || len(value) == 0 {
			err = errors.New("error to parse params")
			return
		}
		pars[key] = value
	}
	params = pars
	return
}

