package diskcache

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
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
	Dir          string
	MaxBytes     int64
	MaxFiles     int64
	CleanupSleep time.Duration
	Shutdown     chan interface{}
	FileNamer    func(fname string) string
}

// new disk cache with sensible defaults
func NewDiskCache() *DiskCache {
	return &DiskCache{
		Dir:          os.TempDir(),
		MaxBytes:     1 << 20, // 1mb
		MaxFiles:     256,
		CleanupSleep: 60 * time.Second,
		FileNamer:    OpportunisticNamer,
	}
}

func CopyNamer(key string) string {
	return key
}

// OpportunisticNamer uses either base64 encoding or MD5+SHA224 hashing depending on key length.
// To further avoid collisions, keys are prefix with "b" for base64 and "s" for hashes.
func OpportunisticNamer(key string) string {
	if float32(len(key))*(4.0/3.0) <= 90.0 {
		return "b" + base64.RawURLEncoding.EncodeToString([]byte(key))
	}
	sha224 := sha256.Sum224([]byte(key))
	md5 := md5.Sum([]byte(key))
	return fmt.Sprintf("h%x%x%x", len(key), sha224, md5)
}

func (c *DiskCache) Start() error {

	if c.MaxBytes <= 0 {
		return fmt.Errorf("MaxBytes cannot be <= 0")
	}

	if c.MaxFiles <= 0 {
		return fmt.Errorf("MaxFiles cannot be <= 0")
	}

	if c.CleanupSleep <= 0 {
		return fmt.Errorf("CleanupSleep cannot be <= 0")
	}

	c.Shutdown = make(chan interface{}, 1)

	go func() {

		ticker := time.NewTicker(c.CleanupSleep)

		for {

			select {
			case <-ticker.C:
				err := c.cleanup()
				if err != nil {
					log.Printf("Error during cleanup: %v", err)
				}
			case <-c.Shutdown:
				log.Printf("Shutting down disk cache")
				return
			}

		}

	}()

	return nil
}

// Read file contents from cache, returns ErrNotFound if not there
func (c *DiskCache) Get(fname string) ([]byte, error) {

	p := c.keyToPath(fname)

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

func (c *DiskCache) Set(fname string, val []byte) error {

	p := c.keyToPath(fname)

	return ioutil.WriteFile(p, val, 0644)

}

func (c *DiskCache) keyToPath(key string) string {
	var name string
	if c.FileNamer == nil {
		name = key
	} else {
		name = c.FileNamer(key)
	}
	return filepath.Join(c.Dir, name)
}

func (c *DiskCache) cleanup() error {

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
		if s > c.MaxBytes || int64(fcount) > c.MaxFiles {
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
