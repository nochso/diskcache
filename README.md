# diskcache

A simple disk-backed cache in golang.

- Map keys to file names
- Timer-based LRU (least recently used) eviction

**[Godoc documentation](https://godoc.org/github.com/nochso/diskcache)**

This is a fork of [bradleypeabody/diskcache](https://github.com/bradleypeabody/diskcache).

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

### DiskCache settings
You can modify the following settings before calling `Start`.
```go
type DiskCache struct {
	// Root directory where files will be stored
	Dir string
	// Maximum amount of bytes to keep when cleaning up.
	// Zero value ignores the limit.
	MaxBytes int64
	// Maximum amount of files to keep when cleaning up.
	// Zero value ignores the limit.
	MaxFiles int64
	// Interval between clean up jobs
	CleanupSleep time.Duration
	// Function to map keys to names
	Mapper func(key string) string
	// If true, items will expire by use age.
	// If false, items will expire by their creation age.
	ModifyOnGet bool
}
```

### Mapping keys to file names
By default, keys are mapped to file names using the `OpportunisticMapper` that chooses between base64 or a combination of hash sums.

The `Mapper` function of `DiskCache` must have this signature: `func(key string) string`
