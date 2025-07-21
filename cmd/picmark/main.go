package main

import (
	"embed"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path"
	"strconv"
	"unicode/utf8"

	"github.com/disintegration/imaging" // 安装：go get -u github.com/disintegration/imaging
	"github.com/go-playground/validator/v10"
	"github.com/goki/freetype"
)

//go:embed fonts/*
var fontFS embed.FS

type Config struct {
	SrcPath   string  `validate:"omitempty"`           // 源文件路径必填
	SrcX      int     `validate:"gte=0"`               // X坐标≥0
	SrcY      int     `validate:"gte=0"`               // Y坐标≥0
	Width     int     `validate:"gt=0"`                // 宽度必须>0
	Height    int     `validate:"gt=0"`                // 高度必须>0
	SrcText   string  `validate:"omitempty"`           // 水印文字可选
	TextColor string  `validate:"omitempty,textcolor"` // 颜色为16进制格式（如ff0000）
	Scale     float64 `validate:"gt=0"`                // 缩放倍数>0
	Rotate    int64   `validate:"gte=0,lte=360"`       // 角度0-360
	Opacity   float64 `validate:"gte=0,lte=1"`         // 透明度0-1.0
	DstPath   string  `validate:"required"`            // 目标路径必填
	DstX      int     `validate:"gte=0"`               // 目标X坐标≥0
	DstY      int     `validate:"gte=0"`               // 目标Y坐标≥0
}

func main() {
	var (
		// ./piccopy -src dst.png -x 150 -y 455 -w 250 -h 50 -dx 150 -dy 370 -s 1
		// 裁剪 src.png 图像粘贴到 dst.png
		srcPath = flag.String("src", "src.png", "裁剪图像路径")
		srcX    = flag.Int("x", 0, "裁剪左上角X坐标")
		srcY    = flag.Int("y", 0, "裁剪左上角Y坐标")
		width   = flag.Int("w", 100, "裁剪宽度")
		height  = flag.Int("h", 100, "裁剪高度")

		srcText   = flag.String("t", "", "水印文字")
		textColor = flag.String("c", "ff0000", "水印文字颜色，格式：ffffff")

		scale   = flag.Float64("s", 1.0, "缩放倍数 >0")
		rotate  = flag.Int64("r", 0, "旋转角度 0-360")
		opacity = flag.Float64("o", 1.0, "透明度 0-1.0")

		dstPath = flag.String("dst", "dst.png", "目标图像路径")
		dstX    = flag.Int("dx", 0, "目标图像的覆盖点X坐标")
		dstY    = flag.Int("dy", 0, "目标图像的覆盖点Y坐标")
	)
	flag.Parse()

	// 绑定到结构体
	config := &Config{
		SrcPath:   *srcPath,
		SrcX:      *srcX,
		SrcY:      *srcY,
		Width:     *width,
		Height:    *height,
		SrcText:   *srcText,
		TextColor: *textColor,
		Scale:     *scale,
		Rotate:    *rotate,
		Opacity:   *opacity,
		DstPath:   *dstPath,
		DstX:      *dstX,
		DstY:      *dstY,
	}

	// 初始化校验器
	validate := validator.New()
	validate.RegisterValidation("textcolor", func(fl validator.FieldLevel) bool {
		c := fl.Field().String()
		if len(c) != 6 {
			return false // 长度必须为6
		}
		// 检查是否为16进制颜色
		for _, r := range c {
			if (r < '0' || r > '9') && (r < 'a' || r > 'f') && (r < 'A' || r > 'F') {
				return false // 非法字符
			}
		}
		return true
	})

	// validate.RegisterValidation("imginfo_check", func(fl validator.FieldLevel) bool {
	// 	config := fl.Parent().Interface().(Config)
	// 	if config.SrcPath == "" && config.SrcText == "" {
	// 		return false // 至少需要一个路径
	// 	}
	// 	if config.DstPath == "" {
	// 		return false // 目标路径必填
	// 	}

	// 	if ios, err := os.Stat(config.DstPath); err != nil || os.IsNotExist(err) {
	// 		fmt.Printf("错误：目标图像 %s 不存在或无法访问\n", config.DstPath)
	// 		return false
	// 	} else if ios.IsDir() {
	// 		fmt.Printf("错误：目标图像 %s 不能是目录\n", config.DstPath)
	// 		return false
	// 	}

	// 	return true
	// })

	if ios, err := os.Stat(config.DstPath); err != nil || os.IsNotExist(err) {
		log.Fatalf("错误：目标图像 %s 不存在或无法访问\n", config.DstPath)
	} else if ios.IsDir() {
		log.Fatalf("错误：目标图像 %s 不能是目录\n", config.DstPath)
	}

	// 执行校验
	if err := validate.Struct(config); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			switch e.Field() {
			default:
				log.Fatalf("错误：%s 参数无效（原因：%s）\n", e.Field(), e.Tag())
			}
		}
	}

	var (
		srcImg    image.Image
		err       error
		textModel bool = config.SrcText != ""
	)

	if textModel {
		c := config.TextColor
		if len(c) != 6 {
			log.Fatal("textColor must be 6 characters long, e.g. 'ff0000'")
		}
		c1, _ := strconv.ParseUint(c[:2], 16, 8)
		c2, _ := strconv.ParseUint(c[2:4], 16, 8)
		c3, _ := strconv.ParseUint(c[4:], 16, 8)
		srcImg = buildImage(config.SrcText, color.RGBA{uint8(c1), uint8(c2), uint8(c3), 0xff})
	} else {
		// 如果没有指定水印文字，则使用裁剪功能
		srcImg, err = imaging.Open(config.SrcPath) // 支持PNG/JPEG等格式
		if err != nil {
			log.Fatal(err)
		}
	}
	dstImg, err := imaging.Open(config.DstPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("file:", config.SrcPath)
	fmt.Println("- width:", srcImg.Bounds().Dx(), "height:", srcImg.Bounds().Dy())
	fmt.Println()
	fmt.Println("file:", config.DstPath)
	fmt.Println("- width:", dstImg.Bounds().Dx(), "height:", dstImg.Bounds().Dy())
	fmt.Println()

	if !textModel {
		if config.SrcX >= srcImg.Bounds().Dx() {
			log.Fatal("srcX out of range srcImg")
		}

		if config.SrcY >= srcImg.Bounds().Dy() {
			log.Fatal("srcY out of range srcImg")
		}

		// 2. 定义源图像的裁剪区域（x1, y1, x2, y2）
		srcX2, srcY2 := config.SrcX+config.Width, config.SrcY+config.Height
		if srcX2 > srcImg.Bounds().Dx() {
			srcX2 = srcImg.Bounds().Dx()
		}
		if srcY2 > srcImg.Bounds().Dy() {
			srcY2 = srcImg.Bounds().Dy()
		}
		srcRect := image.Rect(config.SrcX, config.SrcY, srcX2, srcY2)
		srcImg = imaging.Crop(srcImg, srcRect)
	}

	if int(config.Scale) != 1 {
		// 3. 缩放子图像（可选）
		srcImg = imaging.Resize(srcImg, int((config.Scale)*float64(srcImg.Bounds().Dx())),
			int((config.Scale)*float64(srcImg.Bounds().Dy())), imaging.Lanczos) // 缩放到100x100
	}

	if config.Rotate != 0 {
		// 4. 旋转子图像（可选）
		srcImg = imaging.Rotate(srcImg, float64(config.Rotate), color.Transparent) // 旋转角度
	}

	// 5. 合成图像（保留透明通道）
	// resultImg := imaging.Paste(dstImg, srcImg, dstPoint)
	// 调整透明度（通过叠加透明层）
	resultImg := imaging.Overlay(
		dstImg,
		srcImg,
		image.Point{config.DstX, config.DstY},
		config.Opacity, // 透明度 0-1
	)

	fileDir := path.Dir(config.DstPath)
	fileName := path.Base(config.DstPath)
	ext := path.Ext(fileName)
	fileName = fileName[:len(fileName)-len(ext)]
	// timeTag := strconv.Itoa(int(time.Now().UnixMilli()))
	// timeTag = timeTag[len(timeTag)-6 : len(timeTag)-1]
	timeTag := "result"
	target := fileDir + "/" + fileName + "_" + timeTag + ".png"
	// 6. 保存结果
	err = imaging.Save(resultImg, target) // 自动处理透明通道
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("target:", target)
}

func buildImage(text string, col color.RGBA) image.Image {
	// 24
	width, height := 30*utf8.RuneCountInString(text), 23
	img := imaging.New(width, height, color.Transparent) // 背景
	fontBytes, err := fontFS.ReadFile("fonts/kuaile.ttf")
	if err != nil {
		log.Fatal(err)
	}
	fontFace, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fontFace)
	c.SetFontSize(24)
	c.SetClip(img.Bounds())
	c.SetDst(img) // 目标画布
	c.SetSrc(image.NewUniform(col))

	pt := freetype.Pt(0, 20) // 文本起始坐标 (x, y)
	_, err = c.DrawString(text, pt)
	if err != nil {
		log.Fatal(err)
	}
	// 绘制完成，从图像右侧往左扫描，遇到像素不为零，停止
	for x := img.Bounds().Dx() - 1; x >= 0; x-- {
		for y := 0; y < img.Bounds().Dy(); y++ {
			a, b, c, d := img.At(x, y).RGBA()
			if a+b+c+d > 0 {
				img = imaging.Crop(img, image.Rect(0, 0, x+1, img.Bounds().Dy()))
				return img
			}
		}
	}
	return img
}
