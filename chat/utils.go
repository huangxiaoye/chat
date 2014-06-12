package chat

import (
	"fmt"
	"github.com/huangxiaoye/goutils-genunique"
	"regexp"
)

const (
	NAME_PREFIX = "User "
)

type KVPair map[string]string

func getUniqName() string {
	return fmt.Sprintf("%s%d", NAME_PREFIX, genuniq.GetUnique())
}

func MatchAll(r *regexp.Regexp, str string) (captures []KVPair, ok bool) {
	captures = make([]KVPair, 0, 20)
	names := r.SubexpNames()
	length := len(names)
	matches := r.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		cmap := make(KVPair, length)
		for pos, val := range match {
			name := names[pos]
			if name != "" {
				cmap[name] = val
			}
		}
		captures = append(captures, cmap)
	}

	if len(captures) > 0 {
		ok = true
	}
	return
}
