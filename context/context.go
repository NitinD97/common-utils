package context

import (
	"context"
	"github.com/gin-gonic/gin"
	"sync"
)

type Context struct {
	data       *map[string]any
	mutex      *sync.RWMutex
	Context    context.Context
	GinContext *gin.Context
}

func NewContext() *Context {
	data := make(map[string]any)
	return &Context{
		data:       &data,
		mutex:      &sync.RWMutex{},
		Context:    context.Background(),
		GinContext: nil,
	}
}

func NewContextFromGinContext(ginCtx *gin.Context) *Context {
	data := make(map[string]any)
	for key, value := range ginCtx.Keys {
		data[key] = value.(string)
	}
	return &Context{
		data:       &data,
		mutex:      &sync.RWMutex{},
		Context:    context.Background(),
		GinContext: ginCtx,
	}
}

func (c *Context) Set(key string, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	(*c.data)[key] = value
}

func (c *Context) Get(key string) any {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return (*c.data)[key]
}

func (c *Context) GetContextMap() map[string]any {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return *c.data
}
