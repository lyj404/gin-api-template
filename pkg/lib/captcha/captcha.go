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
	"time"
)

// 生成随机验证码问题
func GenerateMathProblem() CaptchaProblem {
	// 随机选择验证码类型
	captchaType := CaptchaType(rng.Intn(4))

	var question string
	var answer int

	switch captchaType {
	case Addition:
		a := rand.Intn(10) + 1
		b := rand.Intn(10) + 1
		question = fmt.Sprintf("%d + %d = ?", a, b)
		answer = a + b

	case Subtraction:
		a := rand.Intn(15) + 5
		b := rand.Intn(a) + 1
		question = fmt.Sprintf("%d - %d = ?", a, b)
		answer = a - b

	case Multiplication:
		a := rand.Intn(5) + 1
		b := rand.Intn(5) + 1
		question = fmt.Sprintf("%d × %d = ?", a, b)
		answer = a * b

	case Division:
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

// 生成验证码图片
func GenerateCaptchaImage(question string) (string, error) {
	// 图片尺寸
	width, height := 200, 80

	// 创建图像
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 白色背景
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{255, 255, 255, 255}}, image.Point{}, draw.Src)

	// 添加干扰线
	rng.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		x1 := rng.Intn(width)
		y1 := rng.Intn(height)
		x2 := rng.Intn(width)
		y2 := rng.Intn(height)

		// 随机颜色
		lineColor := color.RGBA{
			R: uint8(rng.Intn(128)),
			G: uint8(rng.Intn(128)),
			B: uint8(rng.Intn(128)),
			A: 255,
		}

		drawLine(img, x1, y1, x2, y2, lineColor)
	}

	// 添加干扰点
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

	// 绘制文字
	if err := drawText(img, question, width, height); err != nil {
		return "", err
	}

	// 编码为 base64
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/png;base64," + base64Str, nil
}

// 绘制直线
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, lineColor color.RGBA) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

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

// 绘制文字
func drawText(img *image.RGBA, text string, width, height int) error {
	// 简单的字符绘制（实际项目中可以使用更复杂的字体渲染）
	chars := []rune(text)
	charWidth := width / (len(chars) + 2) // 留出边距

	for i, char := range chars {
		x := charWidth * (i + 1)
		y := height/2 + rand.Intn(10) - 5 // 稍微随机偏移

		// 文字颜色
		textColor := color.RGBA{
			R: uint8(rand.Intn(100) + 50),
			G: uint8(rand.Intn(100) + 50),
			B: uint8(rand.Intn(100) + 50),
			A: 255,
		}

		// 简单的字符绘制
		drawChar(img, char, x, y, textColor, charWidth)
	}

	return nil
}

// 绘制单个字符
func drawChar(img *image.RGBA, char rune, x, y int, textColor color.RGBA, size int) {
	// 根据字符类型选择不同的绘制方式
	switch {
	case char >= '0' && char <= '9':
		drawDigit(img, char, x, y, textColor, size)
	case char == '+' || char == '-' || char == '×' || char == '÷' || char == '=' || char == '?':
		drawSymbol(img, char, x, y, textColor, size)
	default:
	}
}

// 绘制数字
func drawDigit(img *image.RGBA, digit rune, x, y int, textColor color.RGBA, size int) {
	// 数字的点阵模板（5x7 点阵）
	digitTemplates := map[rune][][]bool{
		'0': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
		'1': {
			{false, false, true, false, false},
			{false, true, true, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, true, true, true, false},
		},
		'2': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{false, false, false, false, true},
			{false, false, false, true, false},
			{false, false, true, false, false},
			{false, true, false, false, false},
			{true, true, true, true, true},
		},
		'3': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{false, false, false, false, true},
			{false, false, true, true, false},
			{false, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
		'4': {
			{false, false, false, true, false},
			{false, false, true, true, false},
			{false, true, false, true, false},
			{true, false, false, true, false},
			{true, true, true, true, true},
			{false, false, false, true, false},
			{false, false, false, true, false},
		},
		'5': {
			{true, true, true, true, true},
			{true, false, false, false, false},
			{true, true, true, true, false},
			{false, false, false, false, true},
			{false, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
		'6': {
			{false, false, true, true, false},
			{false, true, false, false, false},
			{true, false, false, false, false},
			{true, true, true, true, false},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
		'7': {
			{true, true, true, true, true},
			{false, false, false, false, true},
			{false, false, false, true, false},
			{false, false, true, false, false},
			{false, true, false, false, false},
			{false, true, false, false, false},
			{false, true, false, false, false},
		},
		'8': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, false},
		},
		'9': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{true, false, false, false, true},
			{false, true, true, true, true},
			{false, false, false, false, true},
			{false, false, false, true, false},
			{false, true, true, false, false},
		},
	}

	if template, exists := digitTemplates[digit]; exists {
		drawTemplate(img, template, x, y, textColor, size)
	}
}

// 绘制符号
func drawSymbol(img *image.RGBA, symbol rune, x, y int, textColor color.RGBA, size int) {
	symbolTemplates := map[rune][][]bool{
		'+': {
			{false, false, false, false, false},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{true, true, true, true, true},
			{false, false, true, false, false},
			{false, false, true, false, false},
			{false, false, false, false, false},
		},
		'-': {
			{false, false, false, false, false},
			{false, false, false, false, false},
			{false, false, false, false, false},
			{true, true, true, true, true},
			{false, false, false, false, false},
			{false, false, false, false, false},
			{false, false, false, false, false},
		},
		'×': {
			{false, false, false, false, false},
			{true, false, false, false, true},
			{false, true, false, true, false},
			{false, false, true, false, false},
			{false, true, false, true, false},
			{true, false, false, false, true},
			{false, false, false, false, false},
		},
		'÷': {
			{false, false, false, false, false},
			{false, false, false, false, false},
			{false, false, true, false, false},
			{false, false, false, false, false},
			{true, true, true, true, true},
			{false, false, false, false, false},
			{false, false, true, false, false},
		},
		'=': {
			{false, false, false, false, false},
			{false, false, false, false, false},
			{true, true, true, true, true},
			{false, false, false, false, false},
			{true, true, true, true, true},
			{false, false, false, false, false},
			{false, false, false, false, false},
		},
		'?': {
			{false, true, true, true, false},
			{true, false, false, false, true},
			{false, false, false, false, true},
			{false, false, false, true, false},
			{false, false, true, false, false},
			{false, false, false, false, false},
			{false, false, true, false, false},
		},
	}

	if template, exists := symbolTemplates[symbol]; exists {
		drawTemplate(img, template, x, y, textColor, size)
	}
}

// 使用模板绘制字符
func drawTemplate(img *image.RGBA, template [][]bool, x, y int, textColor color.RGBA, size int) {
	rows := len(template)
	cols := len(template[0])

	scaleX := size / cols
	scaleY := size / rows

	if scaleX < 1 {
		scaleX = 1
	}
	if scaleY < 1 {
		scaleY = 1
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if template[i][j] {
				// 绘制模板中的点
				for dx := 0; dx < scaleX; dx++ {
					for dy := 0; dy < scaleY; dy++ {
						px := x + j*scaleX + dx - size/2
						py := y + i*scaleY + dy - size/2

						if px >= 0 && px < img.Rect.Dx() && py >= 0 && py < img.Rect.Dy() {
							img.Set(px, py, textColor)
						}
					}
				}
			}
		}
	}
}
