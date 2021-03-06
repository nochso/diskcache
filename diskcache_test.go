package diskcache

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDiskCacheBasic(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "TestDiskCache")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpdir)

	cache := New()
	cache.Dir = tmpdir
	cache.CleanupSleep = time.Second * 3
	err = cache.Start()
	if err != nil {
		panic(err)
	}

	_, err = cache.Get("notexist")
	if err == nil || err != ErrNotFound {
		t.Fatalf("lookup for non existent key should have failed")
	}

	err = cache.Set("1", []byte("some data here"))
	if err != nil {
		panic(err)
	}

	b, err := cache.Get("1")
	if err != nil {
		panic(err)
	}

	if string(b) != "some data here" {
		t.Fatalf("value of b is not as expected, instead was: %v", string(b))
	}

	cache.Stop()

}

func TestDiskCacheMaxFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "TestDiskCache")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpdir)

	cache := New()
	cache.Dir = tmpdir
	cache.CleanupSleep = time.Second * 3
	cache.MaxFiles = 10
	cache.MaxBytes = 10 << 20 // larger than we'll run into
	err = cache.Start()
	if err != nil {
		panic(err)
	}

	cache.Set("test1", []byte("test1"))

	time.Sleep(time.Second * 2)
	cache.Set("test2", []byte("test2"))
	time.Sleep(time.Second * 2)
	cache.Set("test3", []byte("test3"))
	cache.Set("test4", []byte("test4"))
	cache.Set("test5", []byte("test5"))
	cache.Set("test6", []byte("test6"))
	cache.Set("test7", []byte("test7"))
	cache.Set("test8", []byte("test8"))
	cache.Set("test9", []byte("test9"))
	cache.Set("test10", []byte("test10"))
	cache.Set("test11", []byte("test11"))
	cache.Set("test12", []byte("test12"))
	cache.Set("test13", []byte("test13"))
	cache.Set("test14", []byte("test14"))
	time.Sleep(time.Second * 2)
	cache.Set("test15", []byte("test15"))
	time.Sleep(time.Second * 4)

	// should fail
	_, err = cache.Get("test1")
	if err != ErrNotFound {
		t.Fatalf("test1 should not have been found, instead got error obj: %v", err)
	}

	// should not fail
	_, err = cache.Get("test15")
	if err != nil {
		t.Fatalf("test15 should still be in the cache, instead got error: %v", err)
	}
	cache.Stop()
}

func TestDiskCacheMaxSize(t *testing.T) {
	qmeg := make([]byte, 256<<10) // 256k

	tmpdir, err := ioutil.TempDir("", "TestDiskCache")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpdir)

	cache := New()
	cache.Dir = tmpdir
	cache.CleanupSleep = time.Second * 3
	cache.MaxFiles = 1000    // larger than we'll run into
	cache.MaxBytes = 1 << 20 // 1mb cache
	err = cache.Start()
	if err != nil {
		panic(err)
	}

	cache.Set("test1", qmeg)

	time.Sleep(time.Second * 2)
	cache.Set("test2", qmeg)
	time.Sleep(time.Second * 2)
	cache.Set("test3", qmeg)
	cache.Set("test4", qmeg)
	cache.Set("test5", qmeg)
	cache.Set("test6", qmeg)
	cache.Set("test7", qmeg)
	cache.Set("test8", qmeg)
	cache.Set("test9", qmeg)
	cache.Set("test10", qmeg)
	cache.Set("test11", qmeg)
	cache.Set("test12", qmeg)
	cache.Set("test13", qmeg)
	cache.Set("test14", qmeg)
	time.Sleep(time.Second * 2)
	cache.Set("test15", qmeg)
	time.Sleep(time.Second * 4)

	// should fail
	_, err = cache.Get("test1")
	if err != ErrNotFound {
		t.Fatalf("test1 should not have been found, instead got error obj: %v", err)
	}

	// should not fail
	_, err = cache.Get("test15")
	if err != nil {
		t.Fatalf("test15 should still be in the cache, instead got error: %v", err)
	}
	cache.Stop()
}

func TestCopyMapper(t *testing.T) {
	if CopyMapper("foobar") != "foobar" {
		t.Error("Must return original input")
	}
}

func TestOpportunisticFileMapper(t *testing.T) {
	if OpportunisticMapper("short key") != "bc2hvcnQga2V5" {
		t.Error("Short keys must be base64 encoded")
	}
	if OpportunisticMapper(strings.Repeat("x", 100)) != "h64c49c94d0d9ac260177d358388f0461b8ea251125c11187b0878a361eaed563ecafb4bcc5654c597a421547b2" {
		t.Error("Long keys must be hashed")
	}
}
