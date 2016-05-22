# diskcache

A simple disk-backed cache in golang.

- Map keys to file names
- Timer-based eviction of files
    - Least recently used
    - Least recently accessed

**[Godoc documentation](https://godoc.org/github.com/nochso/diskcache)**

- [Installation](#installation)
- [Getting started](#getting-started)
- [DiskCache settings](#diskcache-settings)
    - [New() Defaults](#new-defaults)
- [Mapping keys to file names](#mapping-keys-to-file-names)
- [Credits](#credits)
- [License](#license)

# Installation
```
go get github.com/nochso/diskcache
```

# Getting started

```go
// Create a new DiskCache. See below for the default settings.
cache := diskcache.New()
// Validate and start the clean up handler.
err := cache.Start()
err = cache.Set("thekey", []byte("the value"))
b, err := cache.Get("thekey")
// If you don't need the cache anymore during runtime, stop the clean up ticker.
cache.Stop()
```

# DiskCache settings
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

## New() Defaults
A DiskCache created with `diskcache.New()` defaults to:

```go
DiskCache{
    Dir:          os.TempDir(),
    MaxBytes:     1024 * 1024 * 512, // 512MB
    MaxFiles:     0,
    CleanupSleep: time.Minute * 5,
    ModifyOnGet:  false,
    Mapper:       OpportunisticMapper,
}
```

# Mapping keys to file names
Available mappers:
* `OpportunisticMapper` chooses between base64 or a combination of hash sums based on the length of the key.
* `CopyMapper` uses the key as a file path as it is.

You can use your own mapper by creating a function with this signature: `func(key string) string`

# Credits
This is a fork of [bradleypeabody/diskcache](https://github.com/bradleypeabody/diskcache).

# License
Just like the upstream repo, this library is released under the Apache 2.0 license. See [LICENSE](LICENSE) for the
full license text.