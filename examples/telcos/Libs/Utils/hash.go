package Utils

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
)

// Hash do hash
func Hash(hash crypto.Hash, toHash []byte) []byte {
	var d []byte
	switch hash {
	case crypto.SHA1:
		h := sha1.New()
		h.Write(toHash)
		d = h.Sum(nil)
	case crypto.SHA256:
		h := sha256.New()
		h.Write(toHash)
		d = h.Sum(nil)
	case crypto.MD5:
		h := md5.New()
		h.Write(toHash)
		d = h.Sum(nil)
	}

	return d
}
