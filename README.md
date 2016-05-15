diskcache
=========

A simple disk-backed cache in golang.

Usage:
------

```go
	cache := NewDiskCache()
	cache.Dir = tmpdir
	cache.CleanupSleep = time.Second * 3
	cache.MaxFiles = 1000    // larger than we'll run into
	cache.MaxBytes = 1 << 20 // 1mb cache
	cache.FileNamer = diskcache.CopyNamer // Use keys as file names
	err = cache.Start()
	// if err ...

	err = cache.Set("thekey", []byte("the value"))
	// if err ...
	b, err := cache.Get("thekey")
	// if err ...
```

### Mapping keys to file names
By default, keys are mapped to file names using the `OpportunisticNamer` that chooses between base64 or a combination of hash sums.

The `FileNamer` of `DiskCache` must have this signature: `func(key string) string`
