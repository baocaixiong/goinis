package goinis

import (
	"path"
	"strings"
)

type _util int8

var Util _util = 0

func (u *_util) IsArrayKey(k string) bool {
	return strings.HasSuffix(k, "[]")
}

func (u *_util) IsSubKey(k string) bool {
	return strings.Index(k, ".") >= 0
}

func (u *_util) FileName(fileName string) string {
	return strings.TrimSuffix(path.Base(fileName), path.Ext(fileName))
}
