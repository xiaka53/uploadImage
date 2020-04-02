package uploadImage

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
)

var (
	JPEG       filetype = "image/jpeg"
	PNG        filetype = "image/png"
	GIF        filetype = "image/gif"
	Fixed      narrow   = "fixed"      //固定宽高
	Proportion narrow   = "proportion" //等比宽高缩小
)

type narrow string
type filetype string

type UploadFile struct {
	Param       string     //前端传文件字段
	StandByType []filetype //接收图片类型
	MaxSize     int64      //文件最大限制
	Width       float64    //宽度px
	Height      float64    //高度px
	Narrow      narrow     //缩小类型
	FileName    string     //生成文件名
}

func NewUploadFile(param, path string, w, h float64, maxSice int64, narrows narrow, standBy ...filetype) *UploadFile {
	data := UploadFile{
		Param:       param,
		StandByType: nil,
		MaxSize:     maxSice,
		Width:       w,
		Height:      h,
		Narrow:      narrows,
		FileName:    path,
	}
	for _, v := range standBy {
		data.StandByType = append(data.StandByType, v)
	}
	return &data
}

func (f *UploadFile) UploadImage(c *gin.Context) (string, error) {
	header, err := c.FormFile(f.Param)
	if err != nil {
		return "", err
	}
	if header.Size > f.MaxSize {
		return "", errors.New("File exceeds maximum")
	}
	image_type := header.Header["Content-Type"][0]
	is_type := false
	for _, v := range f.StandByType {
		if string(v) == image_type {
			is_type = true
		}
	}
	if !is_type {
		return "", errors.New("Unknow Type")
	}
	file_image, err := header.Open()
	if err != nil {
		return "", err
	}
	defer file_image.Close()
	img, _, err := image.Decode(file_image)
	if err != nil {
		return "", err
	}
	b := img.Bounds()
	width := b.Max.X
	height := b.Max.Y
	w, h := f.calculateRatioFit(width, height)
	m := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	dst := f.FileName
	if len(f.FileName) < 1 {
		dst = header.Filename
	}
	imgfile, _ := os.Create(dst)
	defer imgfile.Close()
	err = png.Encode(imgfile, m)
	if err != nil {
		return "", err
	}
	return f.FileName, nil
}

func (f *UploadFile) calculateRatioFit(w, h int) (int, int) {
	switch f.Narrow {
	case Fixed:
		return int(f.Width), int(f.Height)
	case Proportion:
		return int(float64(w) * f.Width / 100), int(float64(h) * f.Height / 100)
	}
	return 0, 0
}
