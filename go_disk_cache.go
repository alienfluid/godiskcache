package godiskcache

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
) //import

type GoDiskCache struct {
	cachePrefix string
} //struct

type Params struct {
	Directory string
} //struct

func New(p *Params) *GoDiskCache {
	var directory string = os.TempDir()

	if len(p.Directory) > 0 {
		directory = path.Clean(p.Directory)
		err := os.MkdirAll(directory, 0744)

		if err != nil {
			log.Println(err)
		} //if
	} //if

	dc := &GoDiskCache{}
	dc.cachePrefix = path.Join(directory, "godiskcache_")

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

	//open the cache file
	if file, err := os.Open(dc.buildFileName(key)); err == nil {
		//get stats about the file, need modified time
		if fi, err := file.Stat(); err == nil {
			//check that cache file is still valid
			if int(time.Since(fi.ModTime()).Seconds()) < lifetime {
				//try reading entire file
				if data, err := ioutil.ReadAll(file); err == nil {
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

	//open the file
	if file, err := os.Create(dc.buildFileName(key)); err == nil {
		_, err = file.Write([]byte(data))
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
