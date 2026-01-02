package captcha

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
)

func GenerateMathProblem() CaptchaProblem {
	captchaType := CaptchaType(rng.Intn(4))

	var question string
	var answer int

	switch captchaType {
	case Addition:
		// 加法：1-10 之间的数相加，结果范围 2-20
		a := rand.Intn(10) + 1
		b := rand.Intn(10) + 1
		question = fmt.Sprintf("%d + %d = ?", a, b)
		answer = a + b

	case Subtraction:
		// 减法：保证结果为正数，a 范围 6-20，b 范围 1-a
		a := rand.Intn(15) + 5
		b := rand.Intn(a) + 1
		question = fmt.Sprintf("%d - %d = ?", a, b)
		answer = a - b

	case Multiplication:
		// 乘法：1-5 之间的数相乘，结果范围 1-25
		a := rand.Intn(5) + 1
		b := rand.Intn(5) + 1
		question = fmt.Sprintf("%d × %d = ?", a, b)
		answer = a * b

	case Division:
		// 除法：保证能整除，除数范围 2-5，被除数是除数的倍数
		b := rand.Intn(4) + 2
		a := b * (rand.Intn(5) + 1)
		question = fmt.Sprintf("%d ÷ %d = ?", a, b)
		answer = a / b
	}

	return CaptchaProblem{
		Question: question,
		Answer:   answer,
		Type:     captchaType,
	}
}

func GenerateCaptchaImage(question string) (string, error) {
	width, height := 200, 80

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	// 绘制 5 条随机干扰线，颜色较深
	for i := 0; i < 5; i++ {
		x1 := rng.Intn(width)
		y1 := rng.Intn(height)
		x2 := rng.Intn(width)
		y2 := rng.Intn(height)

		lineColor := color.RGBA{
			R: uint8(rng.Intn(128)),
			G: uint8(rng.Intn(128)),
			B: uint8(rng.Intn(128)),
			A: 255,
		}

		drawLine(img, x1, y1, x2, y2, lineColor)
	}

	// 绘制 50 个随机干扰点，增加噪点
	for i := 0; i < 50; i++ {
		x := rng.Intn(width)
		y := rng.Intn(height)

		dotColor := color.RGBA{
			R: uint8(rng.Intn(256)),
			G: uint8(rng.Intn(256)),
			B: uint8(rng.Intn(256)),
			A: 255,
		}

		img.Set(x, y, dotColor)
	}

	if err := drawSimpleText(img, question, width, height); err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + base64Str, nil
}

// drawLine 使用 Bresenham 算法绘制直线
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, lineColor color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

	// 确定步进方向
	var sx, sy int
	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}

	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		img.Set(x1, y1, lineColor)

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func drawSimpleText(img *image.RGBA, text string, width, height int) error {
	chars := []rune(text)
	charWidth := width / (len(chars) + 2)

	for i, char := range chars {
		x := charWidth * (i + 1)
		y := height/2 + rng.Intn(10) - 5

		textColor := color.RGBA{
			R: uint8(rng.Intn(100) + 50),
			G: uint8(rng.Intn(100) + 50),
			B: uint8(rng.Intn(100) + 50),
			A: 255,
		}

		drawSimpleChar(img, char, x, y, textColor, charWidth)
	}

	return nil
}

func drawSimpleChar(img *image.RGBA, char rune, x, y int, textColor color.RGBA, size int) {
	if char >= '0' && char <= '9' {
		drawSimpleDigit(img, char-'0', x, y, textColor, size)
	} else {
		drawSimpleSymbol(img, char, x, y, textColor, size)
	}
}

func drawSimpleDigit(img *image.RGBA, digit rune, x, y int, textColor color.RGBA, size int) {
	// 数字模板：使用 3x6 点阵绘制数字
	// # 表示绘制点，空格表示空白
	templates := map[rune][]string{
		'0': {"###", "# #", "# #", "# #", "# #", "###"},
		'1': {" # ", "## ", " # ", " # ", " # ", "###"},
		'2': {"###", "  #", "## ", "#  ", "#  ", "###"},
		'3': {"###", "  #", " ##", "  #", " # ", "###"},
		'4': {"# #", "# #", "###", "  #", "  #", "  #"},
		'5': {"###", "#  ", "###", "  #", "  #", "###"},
		'6': {"## ", "#  ", "###", "# #", "# #", "###"},
		'7': {"###", "  #", " # ", " # ", " # ", " # "},
		'8': {"###", "# #", "###", "# #", "# #", "###"},
		'9': {"###", "# #", "###", "  #", "  #", "###"},
	}

	template, exists := templates[digit]
	if !exists {
		return
	}

	// 计算单元格尺寸
	cellWidth := size / 3
	cellHeight := size / 6

	for row, line := range template {
		for col, ch := range line {
			if ch == '#' {
				for dy := 0; dy < cellHeight; dy++ {
					for dx := 0; dx < cellWidth; dx++ {
						px := x + col*cellWidth + dx - size/2
						py := y + row*cellHeight + dy - size/2
						if px >= 0 && px < img.Rect.Dx() && py >= 0 && py < img.Rect.Dy() {
							img.Set(px, py, textColor)
						}
					}
				}
			}
		}
	}
}

func drawSimpleSymbol(img *image.RGBA, symbol rune, x, y int, textColor color.RGBA, size int) {
	// 符号模板：使用 3x6 点阵绘制运算符号
	templates := map[rune][]string{
		'+': {" # ", "###", " # ", "   ", "   ", "   "},
		'-': {"   ", "###", "   ", "   ", "   ", "   "},
		'×': {"# #", " # ", "# #", "   ", "   ", "   "},
		'÷': {" # ", "   ", "###", "   ", " # ", "   "},
		'=': {"   ", "###", "   ", "###", "   ", "   "},
		'?': {"###", "  #", " # ", " # ", "   ", " # "},
	}

	template, exists := templates[symbol]
	if !exists {
		return
	}

	// 计算单元格尺寸
	cellWidth := size / 3
	cellHeight := size / 6

	for row, line := range template {
		for col, ch := range line {
			if ch == '#' {
				for dy := 0; dy < cellHeight; dy++ {
					for dx := 0; dx < cellWidth; dx++ {
						px := x + col*cellWidth + dx - size/2
						py := y + row*cellHeight + dy - size/2
						if px >= 0 && px < img.Rect.Dx() && py >= 0 && py < img.Rect.Dy() {
							img.Set(px, py, textColor)
						}
					}
				}
			}
		}
	}
}
