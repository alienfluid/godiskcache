package godiskcache

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"time"
) //import

type GoDiskCache struct {
	keys      map[string]cacheFile
	directory string
} //struct

type cacheFile struct {
	fileName string
	lifeTime int
} //struct

type Params struct {
	param map[string]interface{}
} //struct

func New(p Params) *GoDiskCache {
	var directory string = os.TempDir()

	if _, ok := p.param["directory"]; ok {
		directory = p.param["directory"].(string)
	} //if

	return &GoDiskCache{keys: make(map[string]cacheFile), directory: directory}
} //New

func (dc *GoDiskCache) Get(key string) (string, error) {
	var err error

	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
		} //if
	}() //func

	//open the cache file
	if file, err := os.Open(dc.directory + dc.keys[key].fileName); err == nil {
		//get stats about the file, need modified time
		if fi, err := file.Stat(); err == nil {
			//check that cache file is still valid
			if int(time.Since(fi.ModTime()).Seconds()) < dc.keys[key].lifeTime {
				//try reading entire file
				if data, err := ioutil.ReadAll(file); err == nil {
					return string(data), err
				} //if
			} //if
		} //if
	} //if

	return "", err
} //Get

func (dc *GoDiskCache) Set(key, data string, lifetime int) error {
	var err error

	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
		} //if
	}() //func

	//convert string to byte slice
	converted := []byte(key)

	//hash the byte slice and return the resulting string
	hasher := sha256.New()
	hasher.Write(converted)
	filename := "godiskcache_" + hex.EncodeToString(hasher.Sum(nil))

	//open the file
	if file, err := os.Create(dc.directory + filename); err == nil {
		_, err = file.Write([]byte(data))
		_ = file.Close()
	} //if

	if err == nil {
		dc.keys[key] = cacheFile{fileName: filename, lifeTime: lifetime}
	} //if

	return err
} //func
