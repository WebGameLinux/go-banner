package types

import (
		"bytes"
		"encoding/hex"
		"os"
		"path/filepath"
		"strconv"
		"strings"
		"sync"
)

func Byte2HexStr(src []byte) string {
		if src == nil || len(src) <= 0 {
				return ""
		}
		var (
				hv   string
				sub  byte
				temp []byte
				res  = new(bytes.Buffer)
		)
		for _, v := range src {
				sub = v & 0xFF
				hv = hex.EncodeToString(append(temp, sub))
				if len(hv) < 2 {
						res.WriteString(strconv.FormatInt(int64(0), 10))
				}
				res.WriteString(hv)
		}
		return res.String()
}

func GetFileHeaderStr(file string, caps ...int) string {
		var ret string
		fs, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
		if err != nil {
				return ret
		}
		defer CloseFile(fs)
		stat, err := fs.Stat()
		if err != nil {
				return ret
		}
		size := stat.Size()
		if size < 0 {
				return ret
		}
		if len(caps) == 0 {
				caps = append(caps, MinHeaderBit)
		}
		buf := make([]byte, caps[0])
		in, err := fs.Read(buf)
		if in < 0 {
				return ret
		}
		return Byte2HexStr(buf)
}

// 关闭文件
func CloseFile(file *os.File) {
		if file != nil {
				_ = file.Close()
		}
}

// 获取文件头信息存储器
func GetFileTypeMap() sync.Map {
		return fileTypeMap
}

// 获取当前文件类型
func GetFileType(file string) string {
		var fileType = filepath.Ext(file)
		header := GetFileHeaderStr(file, MinHeaderBit)
		if header == "" {
				return fileType
		}
		fileTypeMap.Range(func(key, value interface{}) bool {
				if strings.Contains(header, key.(string)) {
						fileType = value.(string)
						return false
				}
				return true
		})
		return fileType
}
