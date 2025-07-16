package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"path"
)

// calculateGroupSums 计算每行像素的灰度值和（分组求和优化）
func calculateGroupSums(img image.Image, currentCutY int) []float64 {
	bounds := img.Bounds()
	width := bounds.Dx()
	groupSums := make([]float64, 12)
	delta := width / 11
	sum := 0.0
	for x := 0; x < width; x += 1 {
		r, g, b, _ := img.At(x, currentCutY).RGBA()
		gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8) // RGB转灰度
		sum += gray
		if (x+1)%delta == 0 {
			groupSums[x/delta] = sum / float64(delta)
			sum = 0.0 // 重置sum为0，准备计算下一个组
		}
	}
	if width%delta != 0 {
		// 处理最后一个组，如果宽度不是delta的整数倍
		groupSums[len(groupSums)-1] = sum / float64(width%delta)
	}
	return groupSums
}

// calculateVariance 基于组和值计算方差
func calculateVariance(groupSums []float64) float64 {
	mean := 0.0
	for _, sum := range groupSums {
		mean += sum
	}
	mean /= float64(len(groupSums))

	variance := 0.0
	for _, sum := range groupSums {
		diff := sum - mean
		variance += diff * diff
	}
	return variance / float64(len(groupSums))
}

// findBestCutLine 滑动窗口寻找最小方差切割线
func findBestCutLine(img image.Image, endY, diffHeight int) int {
	height := img.Bounds().Dy()
	minVariance := math.MaxFloat64
	bestCutY := endY

	for delta := 0; delta <= int(float64(diffHeight)*window); delta += windowStep {
		currentCutY := endY + delta
	MINUS:
		if currentCutY < 0 || currentCutY > height {
			continue // 超出图片边界则跳过
		}

		groupSums := calculateGroupSums(img, currentCutY)
		variance := calculateVariance(groupSums)

		if variance < minVariance {
			minVariance = variance
			bestCutY = currentCutY
		}
		if minVariance < limitVariance+1e-3 {
			break
		}
		if currentCutY > endY {
			currentCutY = endY - delta
			goto MINUS
		}
	}
	return bestCutY
}

// splitImage 切割图片为高宽比2.5的子图（修正高宽比逻辑）
func splitImage(img image.Image) []image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	// subHeight := int(float64(width) * targetRatio)
	divNum := int(math.Round(float64(height) / (float64(width) * targetRatio)))
	if divNum == 0 {
		panic(fmt.Sprintf("targetRatio is too large, must be smaller than %.2f, recommand to use -n flag", float64(height)/float64(width)))
	}
	subHeight := height / divNum
	if windowStep*4 > subHeight {
		panic(fmt.Sprintf("windowStep:%d too big, need be smaller than %d", windowStep, subHeight/4))
	}
	var subImages []image.Image

	for y := 0; y < height; {
		endY := y + subHeight
		if endY+subHeight/3 > height { // 剩余的图片高度过低，合并到前一张图里
			endY = height // 最后一子图可能不足目标高度
		} else {
			// 动态调整切割线
			endY = findBestCutLine(img, endY, subHeight)
		}

		// 切割并保存子图
		subImg := img.(interface {
			SubImage(r image.Rectangle) image.Image
		}).SubImage(image.Rect(0, y, width, endY))

		subImages = append(subImages, subImg)
		y = endY
	}
	return subImages
}

func printImg(name string, img image.Image) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	fmt.Printf("%s: %d x %d\n", name, width, height)
}

var (
	targetRatio = 2.0  // 2800/1260
	subNum      = 0    // 分割成多少个子图片
	window      = 0.25 // ±1260*0.25像素， 默认窗口大小为 子图片高度*0.25
	// 在分辨率为 ​​1260×2800的手机上，微信聊天记录截屏，空白间隔部分，大概40像素，这里取20
	windowStep    = 20
	limitVariance = 1.0 // 图片分割线，所在位置，像素值低于方差1.0，则视为通过。必须>=0
	filepath      = ""
	onlyMsg       = false // 当false时，不生成子图片,只输出信息
)

func main() {
	flag.Float64Var(&targetRatio, "t", 2, "期望子图片的宽高比，近似值")
	flag.IntVar(&subNum, "n", 0, "期望子图片的数量，默认值为0，无效值")
	flag.Float64Var(&window, "w", 0.25, "在切割线±像素范围内，寻找合适切割点，默认窗口大小为 子图片*0.25")
	flag.IntVar(&windowStep, "s", 20, "在分辨率为 ​​1260×2800的手机上，微信聊天记录截屏，空白间隔部分，大概40像素，默认取20")
	flag.Float64Var(&limitVariance, "v", 1.0, "图片分割线所在位置，像素的方差值，低于limitVariance则视为合适的分割点")
	flag.BoolVar(&onlyMsg, "o", false, "当true时，不生成子图片,只输出信息")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <image_path>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:（如果-n存在，则-t失效）")
		flag.PrintDefaults() // 自动打印所有定义的参数及说明
	}

	flag.Parse()

	positionalArgs := flag.Args()
	if len(positionalArgs) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	filepath := positionalArgs[0]

	// 读取图片
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// 根据subNum 修正 targetRatio
	if subNum != 0 {
		targetRatio = float64(img.Bounds().Dy()) / float64(subNum) / float64(img.Bounds().Dx())
		fmt.Println("targetRatio:", targetRatio)
	}

	fileDir := path.Dir(filepath)
	// 获取path路径的文件名
	fileName := path.Base(filepath)
	printImg(fileName, img)

	ext := path.Ext(fileName)
	fileName = fileName[:len(fileName)-len(ext)]

	// 切割图片（目标高宽比，滑动窗口±像素）
	subImages := splitImage(img)

	for i, subImg := range subImages {
		subName := fmt.Sprintf("%s/%s_%d.jpg", fileName, fileName, i)
		printImg(subName, subImg)
		if onlyMsg {
			continue
		}
		subName = fileDir + "/" + subName
		err := os.MkdirAll(path.Dir(subName), 0755)
		if err != nil {
			panic(fmt.Errorf("创建目录失败: %v", err))
		}
		outFile, _ := os.Create(subName)
		defer outFile.Close()
		jpeg.Encode(outFile, subImg, nil)
	}
}
