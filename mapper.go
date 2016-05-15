package diskcache

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// CoypMapper returns the unmodified key.
func CopyMapper(key string) string {
	return key
}

// OpportunisticMapper uses either base64 encoding or MD5+SHA224 hashing depending on key length.
// To further avoid collisions, keys are prefix with "b" for base64 and "s" for hashes.
func OpportunisticMapper(key string) string {
	if float32(len(key))*(4.0/3.0) <= 90.0 {
		return "b" + base64.RawURLEncoding.EncodeToString([]byte(key))
	}
	sha224 := sha256.Sum224([]byte(key))
	md5 := md5.Sum([]byte(key))
	return fmt.Sprintf("h%x%x%x", len(key), sha224, md5)
}
