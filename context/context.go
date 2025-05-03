package context

import (
	"context"
	"sync"
)

type Context struct {
	data    map[string]string
	mutex   *sync.RWMutex
	Context context.Context
}

func NewContext() *Context {
	return &Context{
		data:    make(map[string]string),
		mutex:   &sync.RWMutex{},
		Context: context.Background(),
	}
}

func (c *Context) Set(key string, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

func (c *Context) Get(key string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.data[key]
}

func (c *Context) GetContextMap() map[string]string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.data
}
