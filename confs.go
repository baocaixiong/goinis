package goinis

import (
	"fmt"
	"os"
	"sync"
)

type Configs struct {
	configs map[string]*ConfigFile
}

func NewConfigs(fileNames ...string) (*Configs, error) {
	cs := new(Configs)
	for _, fileName := range fileNames {
		c, err := newConfigFile(fileName)
		if err != nil {
			return nil, err
		}
		cs.configs[Util.FileName(fileName)] = c
	}

	return cs, nil
}

func (cs *Configs) GetConfig(configName string) *ConfigFile {
	return cs.configs[configName]
}

type ConfigFile struct {
	lock     sync.RWMutex
	fileName string
	Key      string

	sections map[string]*Section

	BlockMode bool

	Comment
}

func newConfigFile(fileName string) (*ConfigFile, error) {
	c := new(ConfigFile)
	c.fileName = fileName
	c.Key = Util.FileName(fileName)

	if err := c.loadFile(fileName); err != nil {
		return nil, err
	}

	c.sections = make(map[string]*Section)

	c.BlockMode = true
	return c, nil
}

func (c *ConfigFile) GetSection(key string) (*Section, error) {
	var s *Section
	var has bool
	if s, has = c.sections[key]; has {
		return nil, getError{ErrSectionNotFound, key}
	}

	return s, nil
}

func (c *ConfigFile) SetSection(s *Section) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.HasSectionKey(s.Title) || s == nil {
		return false
	}

	fmt.Println(c.sections, s, "<><>000")

	c.sections[s.Title] = s
	return true
}

func (c *ConfigFile) DeleteSection(section string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if section exists.
	if _, ok := c.sections[section]; !ok {
		return false
	}

	delete(c.sections, section)

	return true
}

func (c *ConfigFile) HasSectionKey(k string) bool {
	_, has := c.sections[k]
	return has
}

func (c *ConfigFile) HasSection(s *Section) bool {
	for _, _s := range c.sections {
		if s == _s {
			return true
		}
	}

	return false
}

func (c *ConfigFile) loadFile(fileName string) (err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.read(f)
}

func (c *ConfigFile) GetSectionList() []*Section {
	list := make([]*Section, len(c.sections))
	for _, s := range c.sections {
		list = append(list, s)
	}
	return list
}

type Comment struct {
	comment string
}

func (c *Comment) AddComment(str string) *Comment {
	c.comment = c.comment + str

	return c
}

func (c *Comment) SetCommnet(str string) *Comment {
	c.comment = str

	return c
}
