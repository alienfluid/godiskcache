godiskcache
===========

##Installation

Install the pkg
<pre><code>go get github.com/cbinsights/godiskcache</code></pre>

##Example

<pre><code>package main

import (
	"github.com/cbinsights/godiskcache"
	"log"
	"os"
)

func main() {
	//create new godiskcache object
	a := godiskcache.New()

	//cache data by providing your key, data to cache, and the amount of time to cache in seconds
	err := a.Set("your key here", "I would like to cache this data!", 3600)

	if err != nil {
		log.Println(err)
	} //if

	//attempt to retrieve the cached data with your key from above
	data, err := a.Get("your key here")

	if err != nil {
		log.Println(err)
	} //if

  //display data
  log.Println(data)
}//main
</code></pre>

##License

The MIT License (MIT)

Copyright (c) 2013 CB Insights

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
