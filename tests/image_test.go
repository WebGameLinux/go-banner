package tests

import (
		"fmt"
		"github.com/WebGameLinux/go-banner/image"
		"github.com/WebGameLinux/go-banner/types"
		"os"
		"path"
		"path/filepath"
		"testing"
)

var TextStr = "1234678"

func TestBannerByte2Str(t *testing.T) {
		var _viewsLogoIndexTpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x34\x8d\xc1\x09\x42\x41\x10\x43\xef\x53\x45\xae\x9e\x52\xc3\x1e\x16\xf9\xa0\xae\xf8\x17\xbc\x04\x52\x91\x35\x58\x81\xbd\x59\x82\x64\xe5\xef\xec\x10\x78\x24\x99\x02\x60\xfc\x9f\x7d\xac\xb3\x76\x01\xb4\x00\x2a\x5f\x44\xc6\x90\x92\x50\x81\xb4\x16\x15\xe2\xca\xf8\xb4\x44\x2c\xae\x62\x41\x0e\x49\xda\x36\x63\x77\x85\xc9\x54\x48\x44\x4c\x61\x4e\xb2\xea\xfb\x7e\x7d\xd0\x70\xeb\x4f\xb4\xfb\x86\xed\xda\xce\x7d\xc7\x3e\xc7\xa3\x63\x8e\x71\x41\x0c\xf5\x0b\x00\x00\xff\xff\xde\x1c\x43\x39\xb9\x00\x00\x00")
		fmt.Println(image.BannerByte2Str(_viewsLogoIndexTpl))
}

func TestCreateImage(t *testing.T) {
		var (
				yOffset  = 12
				bit      = 8
				text     = TextStr
				fontSize = 20
				option   = image.NewContextWrapper()
				typeArr  = []string{image.ImgPNG, image.ImgJPEG, image.ImgGIF}
		)

		option.FontYBit = bit

		option.AutoCreatePath = true
		file, _ := filepath.Abs(os.Args[0])
		option.Path = path.Join(filepath.Dir(file), "tests/cache")
		for _, ty := range typeArr {
				option.High = 20
				option.Width = 80
				option.SaveImageType = ty
				option.FontSize = float64(fontSize)
				option.FontYOffset = yOffset
				if _, ok := image.CreateImage(text, option); !ok {
						t.Errorf("create image error:%s", image.GetError())
				}
				option.ColorAble = true
		}
}

func TestGetImageSize(t *testing.T) {
		file, _ := filepath.Abs(os.Args[0])
		dir := path.Join(filepath.Dir(file), "tests/cache")
		var files = map[string]int{
				dir + "/" + TextStr + ".gif":  1,
				dir + "/" + TextStr + ".jpeg": 2,
				dir + "/" + TextStr + ".png":  3,
		}
		for file, _ := range files {
				data := image.GetImageType(file)
				h := types.GetFileHeaderStr(file)
				// fmt.Println(file)
				if data == image.TypeMap[image.UNKNOWN] {
						t.Errorf("文件类型获取失败")
				}
				if len(h) < types.MinHeaderBit {
						t.Errorf("文件头获取失败")
				}
		}

}
