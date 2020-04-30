package image

import (
		"encoding/json"
		"errors"
		"fmt"
		"github.com/golang/freetype"
		"github.com/golang/freetype/truetype"
		"image"
		"image/gif"
		"image/jpeg"
		"io/ioutil"
		"math"
		"os"
		"path"
		"path/filepath"
		"strings"
		"time"
)

const (
		NumX1    = 100                                                   // 长
		NumY1    = 40                                                    // 宽
		Ext      = ".png"                                                // 文件后缀
		FontFile = "resources/DejaVuSansMono/DejaVuSansCondensed-15.ttf" // 文字
		ImgPNG   = "png"                                                 // 文件类型 png
		ImgJPEG  = "jpeg"                                                // 文件类型 jpeg
		ImgGIF   = "gif"                                                 // 文件类型 gif
)

// 绘制参数
type OptionsDto struct {
		Width                     int     // x1 宽
		High                      int     // y1 长
		Start                     int     // x0 长
		End                       int     // y0 长
		Path                      string  // 存储目录
		ColorAble                 bool    // 是否彩色绘制
		FontFile                  string  // 字体文件
		SaveFile                  string  // 保存文件名
		AutoCreatePath            bool    // 自动创建保存目录
		FontSize                  float64 // 字体大小
		FontDpi                   float64 // 字体dpi
		content                   string  // 文字内容
		FontXOffset               int     // 文字x 轴偏移量 ptx
		FontYOffset               int     // 文字y 轴偏移量 pty
		FontYBit                  int     // 文字y 位置
		SaveImageType             string  // 图片保存类型 png,jpeg,gif
		FontYOffsetWithOutContent int     // 外表计算所得 y轴偏移量
}

// 参数封装器
type ContextWrapper struct {
		FontByte     []byte                               // 字体信息
		ColorDrawer  func(dy, dx int, image *image.NRGBA) // 彩色绘制处理
		JpegOptions  *jpeg.Options                        // jpeg 可选参数
		GifOptions   *gif.Options                         // gif 可选参数
		FileContext  *os.File                             // file
		ImageContext *image.NRGBA                         // image
		OptionsDto
}

func NewContextWrapper() *ContextWrapper {
		return new(ContextWrapper)
}

func NewOptionsDto() *OptionsDto {
		return new(OptionsDto)
}

func (this *ContextWrapper) Dto() *OptionsDto {
		if this == nil {
				return nil
		}
		dto := NewOptionsDto()
		dto.Width = this.Width
		dto.High = this.High
		dto.Start = this.Start
		dto.End = this.End
		dto.Path = this.Path
		dto.ColorAble = this.ColorAble
		dto.FontFile = this.FontFile
		dto.SaveFile = this.SaveFile
		dto.AutoCreatePath = this.AutoCreatePath
		dto.FontSize = this.FontSize
		dto.FontDpi = this.FontDpi
		dto.content = this.content
		dto.FontXOffset = this.FontXOffset
		dto.FontYOffset = this.FontYOffset
		dto.FontYBit = this.FontYBit
		dto.SaveImageType = this.SaveImageType
		dto.FontYOffsetWithOutContent = this.FontYOffsetWithOutContent
		return dto
}

func (this *ContextWrapper) String() string {
		var (
				v   []byte
				err error
		)
		if v, err = json.Marshal(this.Dto()); err != nil {
				return this.Json()
		}
		return string(v)
}

func (this *ContextWrapper) Json() string {
		var m = this.Map()
		if v, err := json.Marshal(m); err == nil {
				return string(v)
		}
		return fmt.Sprintf("%+v", m)
}

func (this *ContextWrapper) Map() map[string]interface{} {
		var m = make(map[string]interface{})
		m["width"] = this.Width
		m["high"] = this.High
		m["start"] = this.Start
		m["end"] = this.End
		m["path"] = this.Path
		m["colorable"] = this.ColorAble
		m["font_file"] = this.FontFile
		m["save_file"] = this.SaveFile
		m["auto_create_path"] = this.AutoCreatePath
		m["font_size"] = this.FontSize
		m["font_dpi"] = this.FontDpi
		m["content"] = this.content
		m["font_x_offset"] = this.FontXOffset
		m["font_y_offset"] = this.FontYOffset
		m["font_y_bit"] = this.FontYBit
		m["save_image_type"] = this.SaveImageType
		m["font_y_offset_without_content"] = this.FontYOffsetWithOutContent
		return m
}

func (this *ContextWrapper) GetFontByte() []byte {
		var (
				err       error
				fontBytes []byte
		)
		if len(this.FontByte) != 0 {
				return this.FontByte
		}
		if fontBytes, err = ioutil.ReadFile(this.GetFontFile()); err != nil {
				panic(err)
		}
		this.FontByte = fontBytes
		return this.FontByte
}

func (this *ContextWrapper) GetFontFile() string {
		var isDef bool
		if this.FontFile == "" {
				this.FontFile = this.getFontFileAbs()
		}
		if strings.Contains(this.FontFile, FontFile) {
				isDef = true
		}
		if !exists(this.FontFile) {
				if isDef {
						panic(errors.New("file font :" + this.FontFile + " ,not exists"))
				}
		}
		return this.FontFile
}

func (this *ContextWrapper) getFontFileAbs() string {
		var current = dir()
		if current == "" {
				panic(errors.New("save path error"))
		}
		return path.Join(dirname(current), FontFile)
}

func (this *ContextWrapper) GetWidth() int {
		if this.Width < 0 || this.Width >= math.MaxInt64 {
				return NumX1
		}
		return this.Width
}

func (this *ContextWrapper) GetX() int {
		if this.Width < 0 || this.Width >= math.MaxInt64 {
				return NumX1
		}
		return this.Width
}

func (this *ContextWrapper) GetStart() int {
		return this.Start
}

func (this *ContextWrapper) GetEnd() int {
		return this.End
}

func (this *ContextWrapper) GetHigh() int {
		if this.High < 0 || this.High >= math.MaxInt64 {
				return NumY1
		}
		return this.High
}

func (this *ContextWrapper) GetY() int {
		if this.High < 0 || this.High >= math.MaxInt64 {
				return NumY1
		}
		return this.High
}

func (this *ContextWrapper) SavePath() string {
		if this.Path == "" {
				file, _ := filepath.Abs(os.Args[0])
				if file == "" {
						panic(errors.New("save path error"))
				}
				this.Path = filepath.Dir(file)
		}
		// 创建
		if this.AutoCreatePath && !exists(this.Path, true) {
				if err := os.MkdirAll(this.Path, os.ModePerm); err == nil {
						this.AutoCreatePath = false
				}
		}
		return this.Path
}

func (this *ContextWrapper) File(text string) string {
		var ext bool
		if strings.Contains(text, ".") {
				ext = true
		}
		if this.SaveFile == "" {
				if !ext {
						text = text + this.getImageExtension()
				}
				this.SaveFile = path.Join(this.SavePath(), text)
		}
		return this.SaveFile
}

func (this *ContextWrapper) getImageExtension() string {
		var ext = this.GetImageType()
		switch ext {
		case ImgJPEG:
				return "." + ImgJPEG
		case ImgPNG:
				return "." + ImgPNG
		case ImgGIF:
				return "." + ImgGIF
		default:
				return "." + ImgPNG
		}
}

func (this *ContextWrapper) GetImageFile(text string) *os.File {
		if this.FileContext != nil {
				return this.FileContext
		}
		if text == "" {
				text = time.Now().Format("2020-10-10.101010")
		}
		var file = this.File(text)
		if file == "" {
				return nil
		}
		if fd, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm); err == nil {
				return fd
		}
		return nil
}

func (this *ContextWrapper) GetFont() *truetype.Font {
		var (
				err     error
				byteBuf = this.GetFontByte()
		)
		if len(byteBuf) == 0 {
				return nil
		}
		if font, err := freetype.ParseFont(byteBuf); err == nil {
				return font
		}
		errVar = err
		return nil
}

func (this *ContextWrapper) GetFontSize() float64 {
		if this.FontSize <= 0 {
				return 36
		}
		return this.FontSize
}

func (this *ContextWrapper) GetFontStartFixed() int {
		if this.FontXOffset <= 0 {
				return 5
		}
		return this.FontXOffset
}

func (this *ContextWrapper) GetDpi() float64 {
		if this.FontDpi <= 0 {
				return 72
		}
		return this.FontDpi
}

func (this *ContextWrapper) Content() string {
		return this.content
}

func (this *ContextWrapper) GetImageType() string {
		if this.SaveImageType != "" {
				return this.SaveImageType
		}
		this.SaveImageType = filepath.Ext(this.SaveFile)
		if this.SaveImageType == "" {
				this.SaveImageType = Ext
		}
		if strings.Contains(this.SaveImageType, ".") {
				ext := strings.Split(this.SaveImageType, ".")
				if len(ext) >= 0 {
						this.SaveImageType = ext[len(ext)-1]
				}
				this.SaveImageType = strings.Trim(strings.ToLower(this.SaveImageType), " ")
		}
		return this.SaveImageType
}

func (this *ContextWrapper) GetJpegOption() *jpeg.Options {
		return this.JpegOptions
}

func (this *ContextWrapper) GetGifOption() *gif.Options {
		return this.GifOptions
}

func (this *ContextWrapper) GetFontXOffset() int {
		if this.FontXOffset < 0 {
				return 10
		}
		return this.FontXOffset
}

func (this *ContextWrapper) GetFontYOffset() int {
		if this.FontYOffset <= 0 {
				return 20
		}
		return this.FontYOffset
}

func (this *ContextWrapper) GetFontYBit() int {
		if this.FontYBit <= 0 {
				return 8
		}
		return this.FontYBit
}

func (this *ContextWrapper) GetImageContext() *image.NRGBA {
		if this.ImageContext == nil {
				this.ImageContext = NewImage(this)
		}
		return this.ImageContext
}

func (this *ContextWrapper) GetFileContext() *os.File {
		if this.FileContext == nil {
				this.FileContext = this.GetImageFile(this.Content())
		}
		return this.FileContext
}

func (this *ContextWrapper) GetContentOffsetY(c *freetype.Context) int {
		if this.FontYOffsetWithOutContent > 0 {
				return this.FontYOffsetWithOutContent
		}
		this.FontYOffsetWithOutContent = this.GetFontYOffset() + int(c.PointToFixed(this.GetFontSize())>>this.GetFontYBit())
		return this.FontYOffsetWithOutContent
}

func (this *ContextWrapper) Release() {
		this.ImageContext = nil
		if this.FileContext != nil {
				_ = this.FileContext.Close()
		}
		this.reset()
}

func (this *ContextWrapper) reset() {
		this.FileContext = nil
		this.JpegOptions = nil
		this.GifOptions = nil
		this.SaveImageType = ""
		this.content = ""
		this.SaveFile = ""
		this.FontYOffsetWithOutContent = 0
		this.ColorAble = false
		this.FontYBit = 0
		this.FontYOffset = 0
		this.FontXOffset = 0
		this.Width = 0
		this.High = 0
		this.AutoCreatePath = false
		this.FontByte = nil
}
