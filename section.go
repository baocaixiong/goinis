package goinis

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const (
	ErrSectionNotFound = iota + 1
	ErrKeyNotFound
	ErrBlankSectionName
	ErrCouldNotParse
	ErrParser
)

type Section struct {
	lock    sync.RWMutex // Go map is not safe.
	Title   string
	content map[string]*keyValue

	subSections map[string]*Section

	configFile *ConfigFile
}

func NewSection(config *ConfigFile, title string) *Section {
	s := new(Section)
	s.content = make(map[string]*keyValue)
	s.subSections = make(map[string]*Section)
	s.configFile = config
	s.Title = title
	return s
}

func (s *Section) SetKeyValue(kv *keyValue) *Section {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.content[kv.K] = kv

	return s
}

func (s *Section) SetValue(key, value string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Check if key exists.
	kv, ok := s.content[key]
	if ok { // 已经包含了。
		kv.setValue(value)
		// @ZHANGMING 向KeyValue中添加值
	} else {
		s.content[key] = newKeyValue(key, value)
	}
	return true
}

func (s *Section) DeleteKey(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.content[key]; ok {
		delete(s.content, key)
	}
}

func (s *Section) GetSubSection(key string) (*Section, error) {
	if !Util.IsSubKey(key) {
		return s.subSections[key], nil
	}

	keys := strings.SplitN(key, ".", 2)

	if j, has := s.subSections[keys[0]]; has {
		return j.GetSubSection(keys[1])
	} else {
		return nil, &getError{ErrSectionNotFound, key}
	}
}

func (s *Section) SetSubSection(subs *Section) *Section {
	s.subSections[subs.Title] = subs
	return s
}

func (s *Section) GetValue(key string) (interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var err error
	if Util.IsSubKey(key) {
		s, err = s.GetSubSection(key)
		if err != nil {
			return nil, &getError{ErrKeyNotFound, key}
		}
	} else {
		value, ok := s.content[key]
		if ok {
			return value.getValue(), nil
		} else {
			return nil, &getError{ErrKeyNotFound, key}
		}
	}

	key = key[strings.LastIndex(key, "."):]
	return s.GetValue(key)
}

func (s *Section) Bool(key string) (bool, error) { // default is false
	value, err := s.GetValue(key)
	if err != nil {
		return false, err
	}
	switch v := value.(type) {
	case string:
		return strconv.ParseBool(v)
	default:
		return false, &getError{ErrCouldNotParse, key}
	}
}

func (s *Section) Float64(key string) (float64, error) {
	value, err := s.GetValue(key)
	if err != nil {
		return 0.0, err
	}
	switch v := value.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0.0, &getError{ErrCouldNotParse, key}
	}
}

func (s *Section) Int(key string) (int, error) { // default is 0
	value, err := s.GetValue(key)
	if err != nil {
		return 0, err
	}
	switch v := value.(type) {
	case string:
		return strconv.Atoi(v)
	default:
		return 0, &getError{ErrCouldNotParse, key}
	}
}

func (s *Section) Int64(key string) (int64, error) {
	value, err := s.GetValue(key)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, &getError{ErrCouldNotParse, key}
	}
}

func (s *Section) MustStringValue(key string, defaultVal ...string) string {
	value, err := s.GetValue(key)
	switch v := value.(type) {
	case string:
		if err != nil && len(defaultVal) > 0 {
			return defaultVal[0]
		} else {
			return v
		}
	default:
		if len(defaultVal) > 0 {
			return defaultVal[0]
		} else {
			return ""
		}
	}
}

func (s *Section) MustStringValueRange(key, defaultVal string, candidates []string) string {
	val, err := s.GetValue(key)
	if err != nil {
		return defaultVal
	}

	switch v := val.(type) {
	case string:
		for _, cand := range candidates {
			if v == cand {
				return v
			}
		}
	default:
		return defaultVal
	}

	return defaultVal
}

func (s *Section) MustBool(key string, defaultVal ...bool) bool {
	value, err := s.Bool(key)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

func (s *Section) MustFloat64(key string, defaultVal ...float64) float64 {
	value, err := s.Float64(key)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

func (s *Section) MustInt(key string, defaultVal ...int) int {
	value, err := s.Int(key)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

func (s *Section) MustInt64(key string, defaultVal ...int64) int64 {
	value, err := s.Int64(key)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

func (s *Section) GetKeyList() []string {
	list := make([]string, len(s.content)-1)
	for key, _ := range s.content {
		list = append(list, key)
	}
	return list
}

func (s *Section) GetKeyValue(key string) (*keyValue, bool) {
	kv, has := s.content[key]
	return kv, has
}

type getError struct {
	Reason int
	Name   string
}

func (err *getError) Error() string {
	switch err.Reason {
	case ErrSectionNotFound:
		return fmt.Sprintf("section '%s' not found", err.Name)
	case ErrKeyNotFound:
		return fmt.Sprintf("key '%s' not found", err.Name)
	case ErrParser:
		return fmt.Sprintf("parse configure file error.")
	}
	Util.Println(err.Reason, err.Name)
	return "invalid get error"
}
