package goinis

import (
	"bufio"
	// "bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	// Default section name.
	DEFAULT_SECTION = "DEFAULT"
)

var LineBreak = "\n"

var SECTCRE, _ = regexp.Compile(`\[(?P<header>[^]]+)\]`)

var OPTCRE, _ = regexp.Compile(`(?P<option>[^:=\s][^:=]*)\s*(?P<vi>[:=])\s*(?P<value>.*)$`)

func init() {
	if runtime.GOOS == "windows" {
		LineBreak = "\r\n"
	}
}

func (c *ConfigFile) read(reader io.Reader) error {
	buf := bufio.NewReader(reader)

	var currentSection *Section = NewSection(c, DEFAULT_SECTION)
	c.SetSection(currentSection)

	var currentKeyValue *KeyValue = nil

	// inKeyValue := false
	//var lastLineObj interface{}

	var lineno = 0

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		lineLengh := len(line)
		if err != nil {
			if err != io.EOF {
				return err
			}

			if lineLengh == 0 {
				break
			}
		}

		lineno += 1

		switch {
		case lineLengh == 0:
			continue

		case line[0] == ' ' && currentKeyValue != nil: // continuation line?
			currentKeyValue.AddValue(line)
		case SECTCRE.Match([]byte(line)):
			titles := SECTCRE.FindStringSubmatch(line)
			title := titles[1]
			section, err := c.GetSection(title)
			if err == nil {
				currentSection = section
			} else {
				currentSection = NewSection(c, title)
				c.SetSection(currentSection)
			}

			currentKeyValue = nil

			continue
		case OPTCRE.Match([]byte(line)):
			matches := OPTCRE.FindStringSubmatch(line)[1:]
			key, value := matches[0], matches[2]
			if currentSection == nil {
				return &getError{ErrParser, "parser error"}
			}
			keyValue, has := currentSection.GetKeyValue(key)
			if !has {
				keyValue = NewKeyValue(key, value)
				currentSection.SetKeyValue(keyValue)
			} else {
				keyValue.SetValue(value)
			}

			currentKeyValue = keyValue

		default:
			continue
		}

		if err == io.EOF {
			break
		}
	}
	return nil
}

func LoadFromData(data []byte) (c *ConfigFile, err error) {
	tmpName := path.Join(os.TempDir(), "goinis", fmt.Sprintf("%d", time.Now().Nanosecond()))
	os.MkdirAll(path.Dir(tmpName), os.ModePerm)
	if err = ioutil.WriteFile(tmpName, data, 0655); err != nil {
		return nil, err
	}

	return NewConfigFile(tmpName)
}

func LoadConfigFile(fileName string, moreFiles ...string) (cs *Configs, err error) {
	fileNames := make([]string, 1, len(moreFiles)+1)
	fileNames[0] = fileName
	if len(moreFiles) > 0 {
		fileNames = append(fileNames, moreFiles...)
	}

	cs, err = NewConfigs(fileNames...)

	return cs, err
}

type readError struct {
	Reason  int
	Content string
}

func (err readError) Error() string {
	switch err.Reason {
	case ErrBlankSectionName:
		return "empty section name not allowed"
	case ErrCouldNotParse:
		return fmt.Sprintf("could not parse line: %s", string(err.Content))
	}
	return "invalid read error"
}
