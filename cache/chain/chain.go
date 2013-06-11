package chain

import (
	"fmt"
	"github.com/pierrre/imageserver"
)

type ChainCache []imageserver.Cache

func (cache ChainCache) Get(key string) (*imageserver.Image, error) {
	for i, c := range cache {
		image, err := c.Get(key)
		if err == nil {
			if i > 0 {
				cache.setCaches(key, image, i)
			}
			return image, nil
		}
	}
	return nil, fmt.Errorf("Image not found in chained cache")
}

func (cache ChainCache) setCaches(key string, image *imageserver.Image, indexLimit int) {
	for i := 0; i < indexLimit; i++ {
		go func(i int) {
			cache[i].Set(key, image)
		}(i)
	}
}

func (cache ChainCache) Set(key string, image *imageserver.Image) (err error) {
	for _, c := range cache {
		go func(c imageserver.Cache) {
			c.Set(key, image)
		}(c)
	}
	return
}