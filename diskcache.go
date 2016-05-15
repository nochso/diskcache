package diskcache

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"
)

var ErrNotFound = fmt.Errorf("Item not found")

// simple disk based cache
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
	Mapper   func(key string) string
	shutdown chan interface{}
}

// new disk cache with sensible defaults
func NewDiskCache() *DiskCache {
	return &DiskCache{
		Dir:          os.TempDir(),
		MaxBytes:     1 << 20, // 1mb
		MaxFiles:     256,
		CleanupSleep: 60 * time.Second,
		Mapper:       OpportunisticMapper,
	}
}

// Start validates and starts a DiskCache.
func (c *DiskCache) Start() error {
	if c.MaxBytes < 0 {
		return fmt.Errorf("MaxBytes cannot be < 0")
	}
	if c.MaxFiles < 0 {
		return fmt.Errorf("MaxFiles cannot be < 0")
	}
	if c.CleanupSleep <= 0 {
		return fmt.Errorf("CleanupSleep cannot be <= 0")
	}
	if c.Mapper == nil {
		return fmt.Errorf("Mapper cannot be nil")
	}
	c.shutdown = make(chan interface{}, 1)
	go func() {
		ticker := time.NewTicker(c.CleanupSleep)
		for {
			select {
			case <-ticker.C:
				err := c.cleanup()
				if err != nil {
					log.Printf("Error during cleanup: %v", err)
				}
			case <-c.shutdown:
				ticker.Stop()
				log.Printf("Stopped diskcache clean up ticker")
				return
			}
		}
	}()
	return nil
}

// Stop the clean up ticker.
func (c *DiskCache) Stop() {
	log.Printf("Stopping diskcache")
	c.shutdown <- true
}

// Read file contents from cache, returns ErrNotFound if not there
func (c *DiskCache) Get(key string) ([]byte, error) {
	p := c.keyToPath(key)

	// update timestamp
	now := time.Now()
	os.Chtimes(p, now, now)

	// read file contents
	b, err := ioutil.ReadFile(p)
	if err != nil {
		// FIXME: should do more to distinguish between file not found and other errors
		return nil, ErrNotFound
	}
	return b, nil
}

// Set a value by writing to disk.
func (c *DiskCache) Set(key string, val []byte) error {
	p := c.keyToPath(key)
	return ioutil.WriteFile(p, val, 0644)
}

// Combines base directory with the mapped key.
func (c *DiskCache) keyToPath(key string) string {
	return filepath.Join(c.Dir, c.Mapper(key))
}

func (c *DiskCache) cleanup() error {
	if c.MaxFiles == 0 && c.MaxBytes == 0 {
		return nil
	}
	files := make(FDataList, 0, 256)

	err := filepath.Walk(c.Dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, FData{
			Path:    path.Base(info.Name()),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		})
		return nil
	})
	if err != nil {
		return nil
	}

	// sort by date
	sort.Sort(files)

	// trim down until size and file count limitations are met
	s := files.TotalSize()
	fcount := len(files)
	for i := 0; i < len(files); i++ {
		if (c.MaxBytes > 0 && s > c.MaxBytes) || (c.MaxFiles > 0 && int64(fcount) > c.MaxFiles) {
			s -= files[i].Size
			os.Remove(filepath.Join(c.Dir, files[i].Path))
			fcount--
		} else {
			break
		}
	}
	return nil
}

type FData struct {
	Path    string
	ModTime time.Time
	Size    int64
}
type FDataList []FData

func (f FDataList) Len() int {
	return len(f)
}

func (f FDataList) Less(i, j int) bool {
	return f[i].ModTime.Unix() < f[j].ModTime.Unix()
}

func (f FDataList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f FDataList) TotalSize() int64 {
	ret := int64(0)
	for _, f0 := range f {
		ret += f0.Size
	}
	return ret
}
