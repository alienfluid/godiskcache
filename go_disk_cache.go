package godiskcache

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/golang/groupcache/lru"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"
) //import

type GoDiskCache struct {
	mutex       sync.RWMutex
	cachePrefix string
	memCache    *lru.Cache
} //struct

type Params struct {
	Directory string
	MemItems  int
} //struct

type DataWrapper struct {
	Ts   time.Time
	Data string
} // struct

func New(p *Params) *GoDiskCache {
	var directory string = os.TempDir()
	var items int = 10000

	if len(p.Directory) > 0 {
		directory = p.Directory
		err := os.MkdirAll(directory, 0744)

		if err != nil {
			log.Println(err)
		} //if
	} //if

	if p.MemItems != 0 {
		items = p.MemItems
	}

	dc := &GoDiskCache{}
	dc.cachePrefix = path.Join(directory, "godiskcache_")
	dc.mutex = sync.RWMutex{}
	dc.memCache = lru.New(items)

	return dc
} //New

func NewParams() *Params {
	return &Params{}
} //NewParams

func (dc *GoDiskCache) Get(key string, lifetime int) (string, error) {
	var err error

	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
		} //if
	}() //func

	// Take the reader lock
	dc.mutex.RLock()
	defer dc.mutex.RUnlock()

	// check the in-memory cache first
	if val, ok := dc.memCache.Get(key); ok {
		dw := val.(DataWrapper)
		if int(time.Since(dw.Ts).Seconds()) < lifetime {
			return string(dw.Data), err
		}
	}

	//open the cache file
	if file, err := os.Open(dc.buildFileName(key)); err == nil {
		defer file.Close()
		//get stats about the file, need modified time
		if fi, err := file.Stat(); err == nil {
			//check that cache file is still valid
			if int(time.Since(fi.ModTime()).Seconds()) < lifetime {
				//try reading entire file
				if data, err := ioutil.ReadAll(file); err == nil {
					// update the cache with this value
					dc.memCache.Add(key, DataWrapper{Ts: fi.ModTime(),
						Data: string(data)})

					return string(data), err
				} //if
			} //if
		} //if
	} //if

	return "", err
} //Get

func (dc *GoDiskCache) Set(key, data string) error {
	var err error

	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
		} //if
	}() //func

	// Take the writer lock
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	//open the file
	if file, err := os.Create(dc.buildFileName(key)); err == nil {
		_, err = file.Write([]byte(data))

		// store it in the in-memory cache
		if fi, err := file.Stat(); err == nil {
			ts := fi.ModTime()
			dc.memCache.Add(key, DataWrapper{Ts: ts, Data: data})
		}

		_ = file.Close()
	} //if

	return err
} //func

func (dc *GoDiskCache) buildFileName(key string) string {
	//hash the byte slice and return the resulting string
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return dc.cachePrefix + hex.EncodeToString(hasher.Sum(nil))
} //buildFileName
