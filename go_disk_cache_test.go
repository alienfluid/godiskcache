package godiskcache

import (
	"crypto/rand"
	"fmt"
	"sync"
	"testing"
)

var _MAX_ITEMS int = 10000
var _MAX_ITERATIONS int = 10
var _MAX_KEY_LEN = 10
var _MAX_VALUE_LEN = 10

func generateRandomString(nchars int) string {
	var alphanum string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, nchars)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func buildRandomKeyValuePairs() map[string]string {
	var results = make(map[string]string)
	for i := 0; i < _MAX_ITEMS; i++ {
		results[generateRandomString(_MAX_KEY_LEN)] = generateRandomString(_MAX_VALUE_LEN)
	}
	return results
}

func addToCache(gc *GoDiskCache, k string, v string, w *sync.WaitGroup) {
	gc.Set(k, v)
	w.Done()
}

func TestConcurrentCacheWrite(t *testing.T) {
	var p = Params{Directory: "godiskcache/"}
	gc := New(&p)

	for run := 0; run < _MAX_ITERATIONS; run++ {
		fmt.Println("Run ", run+1)

		// Generate test data
		kv := buildRandomKeyValuePairs()

		// Store the data into the cache using multiple go routines
		var w sync.WaitGroup
		for k, v := range kv {
			w.Add(1)
			go addToCache(gc, k, v, &w)
		}
		w.Wait()

		// Retrieve the results and make sure they are consistent
		var count int = 0
		for k, v := range kv {
			data, err := gc.Get(k, 3600)

			if err != nil || v != data {
				t.Error("Inconsistent value ", err, " key:", k, "Actual: ", v, "Returned: ", data)
				count++
			}
		}

		fmt.Println("Total: ", len(kv), " Failed: ", count)
	}
}
