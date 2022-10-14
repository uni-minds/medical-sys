package tools

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
)

func GzipToBase64(bs []byte) (b64 string, err error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(bs); err != nil {
		return "", err
	} else if err = w.Flush(); err != nil {
		return "", err
	} else if err = w.Close(); err != nil {
		return "", err
	} else {
		ret := buf.Bytes()
		b64 = base64.StdEncoding.EncodeToString(ret)
		return b64, nil
	}
}
