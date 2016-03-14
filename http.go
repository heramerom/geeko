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
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"

	"crypto/tls"
	"errors"
	"github.com/astaxie/beego/httplib"
	"net/url"
)

func newBeegoRequest(method string, baseUrl string, url string, headers map[string]string, params map[string]string) (req *httplib.BeegoHTTPRequest, err error) {
	u := buildUrl(baseUrl, url)
	m := joinMap(_params, params)

	method = strings.ToUpper(method)

	if method == "GET" {
		p, err := buildParams(_requestSerialization, m)
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(u, "?") {
			u = u[:len(u)-1]
		}
		u = u + "?" + p
	}
	req = httplib.NewBeegoRequest(u, method)
	req.Header("Accept-Encoding", "gzip, deflate")
	req.Header("Accept", "*/*")
	if len(_user) != 0 {
		req.GetRequest().SetBasicAuth(_user, _pwd)
	}

	h := joinMap(headers, _headers)
	for k, v := range h {
		req.Header(k, v)
	}

	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // ignore https

	if method == "POST" {
		if _requestSerialization == "json" {
			req.JSONBody(m)
		} else if _requestSerialization == "xml" {
			b, err := xmlBody(m)
			if err != nil {
				return nil, err
			}
			req.Body(b)
		} else if _requestSerialization == "http" {
			for k, v := range params {
				req.Param(k, v)
			}
		}
	}

	return
}

func xmlBody(m map[string]string) (b []byte, err error) {
	str, err := xml.Marshal(m)
	if err != nil {
		return []byte(""), err
	}
	return []byte(str), nil
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
		err := json.Indent(&output, body, "", "  ")
		if err != nil {
			log.Fatal("Response Json Indent: ", err)
		}
		return output.String()
	}

	return string(body)
}

func buildParams(reqestSerialization string, params map[string]string) (result string, err error) {

	if len(reqestSerialization) == 0 {
		err = errors.New("reqest serialization can not be nil")
		return
	}

	if len(params) == 0 {
		result = ""
		return
	}

	if reqestSerialization == "http" {
		buf := bytes.NewBufferString("")
		for k, v := range params {
			buf.WriteString(url.QueryEscape(k))
			buf.WriteString("=")
			buf.WriteString(url.QueryEscape(v))
			buf.WriteString("&")
		}
		result = buf.String()
		if len(result) > 0 {
			result = result[:len(result)-1]
		}
		return
	} else if reqestSerialization == "json" {
		buf, err := json.Marshal(params)
		if err != nil {
			return "", err
		}
		result = string(buf)
		return result, err
	} else if reqestSerialization == "xml" {
		buf, err := xml.Marshal(params)
		if err != nil {
			return "", err
		}
		result = string(buf)
		return result, nil
	} else {
		err = errors.New("unknow reqest serialization " + Color(reqestSerialization, Red))
		return "", err
	}
	return
}