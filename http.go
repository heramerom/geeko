//
// http.go
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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"crypto/tls"
	"errors"
	"net/url"

	"github.com/astaxie/beego/httplib"
)

func NewBeegoRequest(method string, url string, header string, param string, serialization string) (req *httplib.BeegoHTTPRequest, err error) {
	url = buildUrl(BaseUrl, url)
	if len(url) == 0 {
		err = errors.New("request url can not be empty")
		return
	}

	hs, err := parseHeader(header)
	if err != nil {
		return
	}
	hs = joinMap(Headers, hs)

	ps, err := parseParams(param)
	if err != nil {
		return
	}
	ps = joinMap(Params, ps)

	method = strings.ToUpper(method)

	if method == "GET" {
		var query string
		if serialization == "http" {
			query = param
		} else if serialization == "form" {
			query = BuildFormPrams(ps)
		}
		if len(query) != 0 {
			url = url + "?" + query
		}
	}
	req = httplib.NewBeegoRequest(url, method)
	req.Header("Accept-Encoding", "gzip, deflate")
	req.Header("Accept", "*/*")
	req.SetEnableCookie(EnableCookie)
	if len(User) != 0 {
		req.GetRequest().SetBasicAuth(User, Pwd)
	}

	for k, v := range hs {
		req.Header(k, v)
	}

	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // ignore https

	if method == "POST" {
		if serialization == "json" {
			req.JSONBody(ps)
		} else if serialization == "form" {
			for k, v := range ps {
				req.Param(k, v)
			}
		} else if serialization == "http" {
			req.Body(param)
		}
	}
	return
}

func BuildFormPrams(m map[string]string) (body string) {
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(url.QueryEscape(k))
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(v))
		buf.WriteString("&")
	}
	body = buf.String()
	return
}

func buildUrl(baseUrl string, url string) string {
	if len(baseUrl) == 0 && len(url) == 0 {
		return ""
	}
	if len(baseUrl) == 0 {
		return url
	}
	if len(url) == 0 {
		return baseUrl
	}
	if strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl[:len(baseUrl)-1]
	}
	if strings.HasPrefix(url, "/") {
		url = url[1:]
	}
	return baseUrl + "/" + url
}

func joinMap(m1, m2 map[string]string) (m map[string]string) {
	m = make(map[string]string)
	if m1 != nil {
		for k, v := range m1 {
			m[k] = v
		}
	}
	if m2 != nil {
		for k, v := range m2 {
			m[k] = v
		}
	}
	return
}

func formatResponseBody(res *http.Response, httpreq *httplib.BeegoHTTPRequest, pretty bool) string {
	body, err := httpreq.Bytes()
	if err != nil {
		log.Fatalln("can't get the url", err)
	}
	fmt.Println("")
	if pretty && strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		var output bytes.Buffer
		err := json.Indent(&output, body, "\t", "    ")
		if err != nil {
			log.Fatal("Response Json Indent: ", err)
		}
		return output.String()
	}

	return string(body)
}
