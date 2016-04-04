//
// complete.go
//
// Create at 2016-03-15
//
// Copyright (C) 2016 xhl <heramerom@163.com>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"strings"

	"github.com/chzyer/readline"
)

type ListComplete struct {
	prefixComplete *readline.PrefixCompleter
}

func NewListComplete(pc *readline.PrefixCompleter) *ListComplete {
	return &ListComplete{pc}
}

func (this *ListComplete) Do(line []rune, pos int) (newLine [][]rune, length int) {
	newLine, length = this.prefixComplete.Do(line, pos)
	s := strings.TrimSpace(string(line))
	cms := []string{"list ", "ls ", "l ", "do ", "remove ", "save "}
	var cmd string
	for _, v := range cms {
		if strings.HasPrefix(s, v) {
			cmd = v
		}
	}
	if len(cmd) != 0 {
		sub := string(line[len(cmd):])
		for k, _ := range CmdItems {
			if strings.HasPrefix(k, sub) {
				newLine = append(newLine, []rune(k[pos-len(cmd):]))
			}
		}
	}
	return
}
