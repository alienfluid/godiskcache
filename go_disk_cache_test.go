package godiskcache

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"sync"
	"testing"
)

var _MAX_ITEMS int = 5000
var _MAX_ITERATIONS int = 5
var _MAX_KEY_LEN = 128
var _MAX_VALUE_LEN = 1024

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

func BenchmarkCacheReads(b *testing.B) {
	var p = Params{Directory: "godiskcache/"}
	gc := New(&p)

	// Generate the test data
	kv := buildRandomKeyValuePairs()

	// Store the data in the cache concurrently
	var w sync.WaitGroup
	for k, v := range kv {
		w.Add(1)
		go addToCache(gc, k, v, &w)
	}
	w.Wait()

	// Get a list of all keys so that we can randomly pick them
	keys := make([]string, len(kv))
	for k, _ := range kv {
		keys = append(keys, k)
	}

	// Get a list of random indices since we don't want to generate them in the loop
	inds := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		inds = append(inds, mrand.Intn(len(kv)))
	}

	// Reset the timer to begin the read test
	b.ResetTimer()

	// Test loop
	for i := 0; i < b.N; i++ {
		_, _ = gc.Get(keys[inds[i]], 3600)
	}
}

func BenchmarkCacheWrites(b *testing.B) {
	var p = Params{Directory: "godiskcache/"}
	gc := New(&p)

	// Generate the test data
	kv := buildRandomKeyValuePairs()

	// Store the key/value pairs in lists for each access
	keys := make([]string, len(kv))
	values := make([]string, len(kv))

	for k, v := range kv {
		keys = append(keys, k)
		values = append(values, v)
	}

	// Get a list of random indices since we don't want to generate them in the loop
	inds := make([]int, b.N)
	for i := 0; i < b.N; i++ {
		inds = append(inds, mrand.Intn(len(kv)))
	}

	// Reset the timer to begin the read test
	b.ResetTimer()

	// Test loop
	for i := 0; i < b.N; i++ {
		_ = gc.Set(keys[inds[i]], values[inds[i]])
	}
}
