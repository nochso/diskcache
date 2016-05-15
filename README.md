# diskcache

A simple disk-backed cache in golang.

- Map keys to file names
- Timer-based LRU (least recently used) eviction

**[Godoc documentation](https://godoc.org/github.com/nochso/diskcache)**

## Getting started

```go
cache := diskcache.NewDiskCache()
cache.Dir = "cache" // defaults to os.TempDir()
err = cache.Start()
// if err ...

err = cache.Set("thekey", []byte("the value"))
// if err ...
b, err := cache.Get("thekey")
// if err ...
```

### Mapping keys to file names
By default, keys are mapped to file names using the `OpportunisticMapper` that chooses between base64 or a combination of hash sums.

The `Mapper` function of `DiskCache` must have this signature: `func(key string) string`
