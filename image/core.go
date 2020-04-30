package image

import (
		"bufio"
		"errors"
		"github.com/WebGameLinux/go-banner/types"
		"github.com/golang/freetype"
		"image"
		"image/color"
		"image/gif"
		"image/jpeg"
		"image/png"
		"os"
		"path/filepath"
)

var errVar error

// 获取异常
func GetError() error {
		var err = errVar
		if err != nil {
				errVar = nil
		}
		return err
}

// 关闭文件
func closeFile(file *os.File) {
		if file != nil {
				_ = file.Close()
		}
}

// 文件是否存在
func exists(file string, dir ...bool) bool {
		var (
				err  error
				info os.FileInfo
		)
		if len(dir) == 0 {
				dir = append(dir, false)
		}
		if info, err = os.Stat(file); err != nil {
				if os.IsNotExist(err) {
						return false
				}
				if os.IsPermission(err) {
						if err = os.Chmod(file, os.ModePerm); err != nil {
								return false
						}
				}
		}
		if info != nil && info.IsDir() && !dir[0] {
				return false
		}
		return true
}

// 当前dir
func dir() string {
		if d, err := filepath.Abs(os.Args[0]); err == nil {
				return d
		}
		return ""
}

// 目录
func dirname(filename string, level ...int) string {
		var (
				err  error
				dir  string
				base string
		)
		if len(level) == 0 {
				level = append(level, 1)
		}
		if filename == "" {
				return ""
		}
		if filename == "." || filename == ".." {
				if base, err = filepath.Abs(filename); err != nil {
						errVar = err
						return ""
				}
				if filename == "." {
						return base
				}
				filename = base
		}
		i := level[0]
		dir = filename
		for {
				if dir == "/" {
						return dir
				}
				i--
				filename = filepath.Dir(dir)
				if filename == dir {
						return dir
				}
				dir = filename
				if i <= 0 {
						break
				}
		}
		return dir
}

// 彩色绘制
func ColorDraw(dy, dx int, image *image.NRGBA) {
		if image == nil {
				return
		}
		// 画背景,这里可根据喜好画出背景颜色
		for y := 0; y < dy; y++ {
				for x := 0; x < dx; x++ {
						//设置某个点的颜色，依次是 RGBA
						r, g := uint8(x), uint8(y)
						image.Set(x, y, color.RGBA{R: r, G: g, B: 0, A: 255})
				}
		}
		return
}

// 创建image
func NewImage(option *ContextWrapper) *image.NRGBA {
		if option.ImageContext != nil {
				return option.ImageContext
		}
		var img = image.NewNRGBA(image.Rect(option.GetStart(), option.GetEnd(), option.GetX(), option.GetY()))
		if option.ColorAble {
				if option.ColorDrawer == nil {
						ColorDraw(option.GetY(), option.GetX(), img)
				} else {
						option.ColorDrawer(option.GetY(), option.GetX(), img)
				}
		}
		return img
}

// 创建image
func CreateImage(text string, opts ...*ContextWrapper) (string, bool) {
		if len(opts) == 0 {
				opts = append(opts, new(ContextWrapper))
		}
		option := opts[0]
		// 文字内容
		option.content = text
		var (
				err     error
				imgFile = option.GetImageFile(text)
				file    = imgFile.Name()
		)
		option.FileContext = imgFile
		defer option.Release()
		//创建位图,坐标x,y,长宽x,y
		//  100,40
		img := NewImage(option)
		//读字体数据
		if font := option.GetFont(); font == nil {
				return file, false
		}
		option.ImageContext = img
		// 创建
		if err = CreateFontImage(option); err != nil {
				errVar = err
				return file, false
		}
		return file, true
}

// 创建
func CreateFontImage(opt *ContextWrapper, args ...interface{}) error {
		var (
				file *os.File
				img  *image.NRGBA
				c    = freetype.NewContext()
		)
		for _, v := range args {
				if m, ok := v.(*image.NRGBA); ok && img == nil {
						img = m
				}
				if f, ok := v.(*os.File); ok && file == nil {
						file = f
				}
		}
		if img == nil {
				img = opt.GetImageContext()
		}
		if file == nil {
				file = opt.GetFileContext()
		}
		c.SetDPI(opt.GetDpi()) //72
		c.SetFont(opt.GetFont())
		c.SetFontSize(opt.GetFontSize())
		c.SetSrc(image.White)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		//设置字体显示位置
		pt := freetype.Pt(opt.GetFontXOffset(), opt.GetContentOffsetY(c))
		// 绘制文字
		if _, err := c.DrawString(opt.Content(), pt); err != nil {
				return err
		}
		//保存图像到文件
		return SaveImageHandler(opt, file, img)
}

// 保存图片
func SaveImageHandler(opt *ContextWrapper, file *os.File, img *image.NRGBA) error {
		var ty = opt.GetImageType()
		switch ty {
		case "png":
				fallthrough
		case "Png":
				fallthrough
		case "PNG":
				fallthrough
		case ".png":
				return png.Encode(file, img)
		case "JPEG":
				fallthrough
		case "jpeg":
				fallthrough
		case "jpg":
				fallthrough
		case ".jpg":
				fallthrough
		case ".jpeg":
				return jpeg.Encode(file, img, opt.GetJpegOption())
		case "GIF":
				fallthrough
		case "gif":
				fallthrough
		case ".gif":
				fallthrough
		case "Gif":
				return gif.Encode(file, img, opt.GetGifOption())
		default:
				return png.Encode(file, img)
		}
}

// 获取文件类型
func GetImageType(file string) string {
		ty := types.GetFileType(file)
		if IsImageType(ty) {
				return ty
		}
		return GetImageIndexToType(0)
}

// 获取图片资源
func OpenImage(file string) (image.Image, error) {
		if !exists(file) {
				return nil, errors.New("image:" + file + " not exists")
		}
		var (
				fs  *os.File
				err error
				img image.Image
		)
		if fs, err = os.OpenFile(file, os.O_RDWR, os.ModePerm); err != nil {
				return nil, err
		}
		buf := bufio.NewReader(fs)
		switch GetImageType(file) {
		case ImgPNG:
				img, err = png.Decode(buf)
		case ImgJPEG:
				img, err = jpeg.Decode(buf)
		case ImgGIF:
				img, err = gif.Decode(buf)
		}
		if img == nil {
				return nil, errors.New("unknown image type")
		}
		return img, nil
}

// 获取图片文件上下文
func GetImageContext(file string) image.Image {
		if img, err := OpenImage(file); err == nil {
				return img
		}
		return nil
}
