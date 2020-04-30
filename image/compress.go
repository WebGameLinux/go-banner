package image

import (
		"bytes"
		"compress/gzip"
		"fmt"
		"io"
)

func Reader(data []byte) ([]byte, error) {
		gz, err := gzip.NewReader(bytes.NewBuffer(data))
		if err != nil {
				return nil, fmt.Errorf("err:%v", err)
		}

		var buf bytes.Buffer
		_, err = io.Copy(&buf, gz)
		clErr := gz.Close()

		if err != nil {
				return nil, fmt.Errorf("err %v", err)
		}
		if clErr != nil {
				return nil, err
		}
		return buf.Bytes(), nil
}

func BannerByte2Str(banner []byte) string {
		if v, err := Reader(banner); err == nil {
				return string(v)
		}
		return ""
}
