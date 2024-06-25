package encrypt

import (
	"crypto/md5" // #nosec-exclude=G401,G501
	"encoding/hex"
)

func Md5(data []byte) string {
	h := md5.New() // #nosec-exclude=G401,G501
	_, _ = h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
