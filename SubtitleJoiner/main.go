package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
)

// 过滤当前目录下的符合给定前缀后缀的文件
var prefix = "[大明王朝1566]"
var suffix = ".jpg"

// 计算字幕区域高度
var subtitleHeightRatio = 0.15

// 计算字幕下方留白高度
var subtitleBottomMarginRatio = 0.05

func main() {
	// 获取当前目录路径
	dir, _ := os.Getwd()

	// 列出当前目录下的所有文件和目录
	files, _ := ioutil.ReadDir(dir)

	//保存所有符合条件的图片
	var images []image.Image
	for _, file := range files {
		if !file.IsDir() && len(file.Name()) > len(suffix) && len(file.Name()) > len(prefix) {
			if file.Name()[0:len(prefix)] == prefix &&
				file.Name()[len(file.Name())-len(suffix):] == suffix {
				fmt.Println(file.Name())
				imageFile, _ := os.Open(file.Name())
				defer imageFile.Close()
				img, _ := jpeg.Decode(imageFile)
				images = append(images, img)
			}
		}
	}

	//	获取原始图片大小以及图片截取范围
	var originImgSize []int
	originImgSize = append(originImgSize, images[0].Bounds().Dx())
	originImgSize = append(originImgSize, images[0].Bounds().Dy())
	fmt.Printf("原始图片大小：%d x %d\n", originImgSize[0], originImgSize[1])
	subtitleHeight := int(subtitleHeightRatio * float64(originImgSize[1]))
	subtitleBottomMargin := int(subtitleBottomMarginRatio * float64(originImgSize[1]))
	//截取字幕区域
	subtitleRect := image.Rect(0, originImgSize[1]-subtitleHeight-subtitleBottomMargin,
		originImgSize[0], originImgSize[1]-subtitleBottomMargin)
	//	计算最后输出图片的尺寸
	var outputImgSize []int
	outputImgSize = append(outputImgSize, originImgSize[0])
	outputImgSize = append(outputImgSize,
		originImgSize[1]+(subtitleHeight-subtitleBottomMargin)*(len(images)-1)-subtitleBottomMargin)
	fmt.Printf("最后输出图片大小：%d x %d\n", outputImgSize[0], outputImgSize[1])
	//	创建最后输出图片
	result := image.NewRGBA(image.Rect(0, 0, outputImgSize[0], outputImgSize[1]))
	//处理第一张图片的底部空白
	firstImg := images[0].(*image.YCbCr).SubImage(
		image.Rect(0, 0, originImgSize[0], originImgSize[1]-subtitleBottomMargin))
	//	将第一张图片写入最后输出图片
	draw.Draw(result, firstImg.Bounds(), firstImg, image.Point{0, 0}, draw.Src)
	//	将剩余字幕部分写入最后输出图片
	for i := 1; i < len(images); i++ {
		// 截取指定矩形区域的子图像
		subImg := images[i].(*image.YCbCr).SubImage(subtitleRect)
		//把字幕图片保存到子目录subtitle中,方便查看
		subtitleFile, _ := os.Create("subtitle/" + files[i].Name())
		defer subtitleFile.Close()
		jpeg.Encode(subtitleFile, subImg, nil)
		//要绘制的目标区域
		dstRect := image.Rect(0, originImgSize[1]-subtitleBottomMargin,
			originImgSize[0], originImgSize[1]+subtitleHeight-subtitleBottomMargin)
		//把字幕图片绘制到最后输出图片上
		draw.Draw(result, dstRect.Add(
			image.Point{0, (subtitleHeight - subtitleBottomMargin) * (i - 1)}),
			subImg, image.Point{0, originImgSize[1] - subtitleHeight}, draw.Src)
	}
	//	将最后输出图片写入文件
	out, _ := os.Create("result.jpg")
	defer out.Close()
	jpeg.Encode(out, result, nil)
	fmt.Printf("最后输出图片大小：%d x %d\n", outputImgSize[0], outputImgSize[1])
}
