package check_passport

import (
	"io/fs"
	"os"
	"sync"
)

type dataCache struct {
	sync.RWMutex
	items map[string]*dataItemCache
}

type dataItemCache struct {
	file *os.File
	data []byte
}

func (i *dataItemCache) CloseFile() (err error) {
	if err = i.Save(); err != nil {
		return
	}
	err = i.file.Sync()
	_ = i.file.Close()
	return
}

func (i *dataItemCache) Save() (err error) {
	if len(i.data) == 0 {
		return
	}

	if _, err = i.file.Write(i.data); err != nil {
		return
	}

	i.data = []byte{}
	return
}

func (i *dataItemCache) FileName() string {
	if i.file == nil {
		return ""
	}
	return i.file.Name()
}

func NewDataCache() *dataCache {
	return &dataCache{
		items: map[string]*dataItemCache{},
	}
}

func (c *dataCache) OpenFile(filepath string) (file *os.File, err error) {
	c.Lock()
	defer c.Unlock()

	if f, ok := c.items[filepath]; ok {
		file = f.file
		return
	}

	file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.ModeAppend|filePermission)
	if err == nil {
		c.items[filepath] = &dataItemCache{
			file: file,
			data: []byte{},
		}
	}
	return
}

func (c *dataCache) CloseFile(filepath string) (err error) {
	c.Lock()
	defer c.Unlock()

	var (
		item *dataItemCache
		ok   bool
	)

	if item, ok = c.items[filepath]; !ok {
		return
	}

	if err = item.CloseFile(); err == nil {
		delete(c.items, filepath)
	}
	return
}

func (c *dataCache) CloseAllFile() (errs map[string]error) {
	c.Lock()
	defer c.Unlock()

	newItems := map[string]*dataItemCache{}
	errs = map[string]error{}

	for k, item := range c.items {
		if e := item.CloseFile(); e != nil {
			errs[item.FileName()] = e
			newItems[k] = item
			continue
		}
	}

	c.items = newItems
	return
}

func (c *dataCache) AddData(filepath string, b []byte) (err error) {
	c.Lock()
	defer c.Unlock()

	var (
		item *dataItemCache
		ok   bool
	)

	if item, ok = c.items[filepath]; !ok {
		return
	}

	item.data = append(item.data, b...)
	if len(item.data) >= 8192 {
		err = item.Save()
	}
	return
}
