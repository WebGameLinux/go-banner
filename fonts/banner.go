package fonts

import (
		"errors"
		image2 "github.com/WebGameLinux/go-banner/image"
		"os"
		"strings"
)

var bannerError error

const (
		SizeTypeB = "big"
		SizeTypeS = "small"
		SizeTypeM = "medium"
		// DrawChars = ".@80GCLft1i;:, "
		DrawChars = `/\-_.Y `
)

type DrawArgsOptionDto struct {
		BitMap           string
		SizeType         string
		DrawImageOptions *image2.ContextWrapper
}

// 创建绘制参数
func NewDrawArgsOption() *DrawArgsOptionDto {
		var args = new(DrawArgsOptionDto)
		args.BitMap = DrawChars
		args.SizeType = SizeTypeM
		args.DrawImageOptions = image2.NewContextWrapper()
		return args
}

// 绘制字符画
func DrawText2ImageAscii(text string, options ...*DrawArgsOptionDto) string {
		if text == "" {
				return ""
		}
		if len(options) == 0 {
				options = append(options, NewDrawArgsOption())
		}
		var (
				ok   bool
				file string
				opt  = options[0]
				size = len(text)
		)
		opt.DrawImageOptions.High = 40
		opt.DrawImageOptions.Width = size * 26
		if file, ok = image2.CreateImage(text, opt.DrawImageOptions); !ok || file == "" {
				bannerError = image2.GetError()
				return ""
		}
		// defer DelFile(file)
		return ImageToAscii(file, opt.SizeType, opt.BitMap)
}

// 删除文件
func DelFile(file string) {
		bannerError = os.Remove(file)
}

// 图片输出字符串图
func ImageToAscii(image string, size string, bitMapChars ...string) string {
		var (
				y       int
				x       int
				stepX   = 4
				stepY   = 8
				output  string
				bitMap  []string
				bitSize int
				offset  int
				img     = image2.GetImageContext(image)
		)
		if img == nil {
				err := image2.GetError()
				if err != nil {
						bannerError = err
				} else {
						bannerError = errors.New("image open failed ")
				}
				return output
		}
		if len(bitMapChars) == 0 {
				bitMapChars = append(bitMapChars, DrawChars)
		}
		bitMap = String2Array(bitMapChars[0])
		bitSize = len(bitMap)
		bounds := img.Bounds().Size()
		y = bounds.Y
		x = bounds.X
		if y < 100 {
				stepY = 2
		}
		if x < 100 {
				stepX = 2
				size = ""
		}
		switch size {
		case SizeTypeS:
				stepX = 8
				stepY = 16
		case SizeTypeM:
				stepX = 4
				stepY = 8
		case SizeTypeB:
				stepX = 2
				stepY = 4
		}
		for j := 0; j < y; j += stepY {
				for i := 0; i < x; i += stepX {
						// $colors=imagecolorsforindex($im,imagecolorat($im,$i,$j));	//获取像素块的代表点RGB信息
						colors := img.At(i, j)
						red, green, blue, _ := colors.RGBA()
						red = red / 256
						green = green / 256
						blue = blue / 256
						//灰度值计算公式：Gray=R*0.3+G*0.59+B*0.11
						greyness := int((30*int(red) + 59*int(green) + 11*int(blue)) / 255 / 100)
						//			$offset=(int)ceil($greyness*(strlen($str)-1));	//根据灰度值选择合适的字符
						offset = greyness * (bitSize - 1)
						if offset == bitSize-1 {
								output += " "
						} else {
								output += bitMap[offset]
						}
				}
				output += "\n"
		}
		return output
}

// 获取异常
func GetError() error {
		var err = bannerError
		if err != nil {
				bannerError = nil
		}
		return err
}

// 字符串转字符数组
func String2Array(str string, split ...string) []string {
		var arr []string
		if str == "" {
				return arr
		}
		if len(split) == 0 {
				rArr := []rune(str)
				for _, v := range rArr {
						arr = append(arr, string(v))
				}
				return arr
		}
		return strings.Split(str, split[0])
}
