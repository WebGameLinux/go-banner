package image

import "strings"

const (
		UNKNOWN = iota
		GIF
		JPEG
		PNG
		SWF
		PSD
		BMP
		TiffIi
		TiffMm
		JPC
		JP2
		JPX
		JB2
		SWC
		IFF
		WBmp
		XBM
		ICO
		COUNT
)

var TypeMap = map[int]string{
		UNKNOWN: "unknown",
		GIF:     "gif",
		JPEG:    "jpeg,jpg",
		PNG:     "png",
		SWF:     "swf",
		PSD:     "psd",
		BMP:     "bmp",
		TiffIi:  "tiff_ti",
		TiffMm:  "tiff_mm",
		JPC:     "jpc",
		JP2:     "jp2",
		JPX:     "jpx",
		JB2:     "jb2",
		SWC:     "swc",
		IFF:     "iff",
		WBmp:    "wbmp",
		XBM:     "xbm",
		ICO:     "ico",
		COUNT:   "count",
}

// 通过文件类枚举值 获取文件类型名
// @param id int : 类型枚举
func GetImageIndexToType(id int) string {
		if v, ok := TypeMap[id]; ok {
				return v
		}
		return TypeMap[0]
}

// 获取文件类型 枚举值
// @param ty string : 文件类型名
func GetImageTypeToIndex(ty string) int {
		for id, v := range TypeMap {
				if strings.Contains(v, ",") {
						str := strings.SplitN(v, ",", -1)
						for _, t := range str {
								if t == ty || strings.EqualFold(t, ty) {
										return id
								}
						}
				}
				if v == ty || strings.EqualFold(v, ty) {
						return id
				}
		}
		return UNKNOWN
}

// 是否图片类型
// @param ty string : 文件类型名
func IsImageType(ty string) bool {
		return GetImageTypeToIndex(ty) != UNKNOWN
}
