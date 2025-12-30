package captcha

import (
	"math/rand"
	"time"

	"github.com/gorilla/sessions"
	"github.com/lyj404/gin-api-template/config"
)

// 验证码类型
type CaptchaType int

const (
	Addition       CaptchaType = iota // 加法
	Subtraction                       // 减法
	Multiplication                    // 乘法
	Division                          // 除法
)

// 验证码问题
type CaptchaProblem struct {
	Question string      `json:"question"`
	Answer   int         `json:"answer"`
	Type     CaptchaType `json:"type"`
}

type CaptchaReponse struct {
	ImageUrl   string `json:"image_url"`
	Question   string `json:"question"`
	ExpireTime int    `json:"expire_time"`
}

// Session 存储配置
var Store = sessions.NewCookieStore([]byte(config.CfgSession.SessionSecret))

// 验证码 Session Key
const CaptchaSessionKey = "captcha_data"

// 全局随机数生成器
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
