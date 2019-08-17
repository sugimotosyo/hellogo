package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/sugimotosyo/hellogo/conf"
	"github.com/sugimotosyo/hellogo/handler"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

// @title LUCK API
// @version 1.0
// @description LUCKのAPIです。

// @contact.name API Support
// @contact.url https://en-joy.co.jp/contact/index/
// @contact.email s.sugimoto@en-joy.co.jp

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
// @BasePath /
func main() {

	//ログの設定
	settingLogger()

	//echoインスタンスの生成
	var e = echo.New()

	//cross
	e.Use(middleware.CORS())

	//Skipperの設定
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper: func(c echo.Context) bool {
			log.Debugf("Skipper")
			return false
		},
		DisableStackAll:   false,
		DisablePrintStack: false,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "check ok")
	})

	/************************
	ユーザー
	*************************/
	twitter := handler.NewTwitter()
	//ユーザ作成画面から他人のユーザを作成する
	e.POST("/twitter/post", twitter.Post, Middle(twitter.MiddleFunc))

	//listen
	log.Fatal(e.Start(conf.Port))
}

//SettingLogger ログの設定使用時は"github.com/labstack/gommon/log"をインポートの上そのまま使う。
func settingLogger() {

	//prefixを設定
	log.SetPrefix(conf.ServiceName + os.Getenv("ENVIRONMENT"))
	//出力レベルを設定
	log.SetLevel(log.DEBUG)
	//色表示にする。
	log.EnableColor()

	//example
	log.Printf("log.Printf")
	log.Debugf("log.Debugf")
	log.Infof("log.Infof")
	log.Warnf("log.Warnf")
	log.Errorf("log.Errorf")

}

//Middle .
func Middle(middleFunc func(echo.Context) error) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := middleFunc(c)
			if err != nil {
				fmt.Println("===================")
				fmt.Printf("middle error %+v\n", err)
				return errors.WithStack(err)
			}
			err = next(c)
			if err != nil {
				fmt.Println("--------------------")
				fmt.Printf("error %+v\n", err)
			}
			return errors.WithStack(err)
		}
	}
}
