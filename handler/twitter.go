package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ChimeraCoder/anaconda"
	"github.com/pkg/errors"
	"github.com/sugimotosyo/hellogo/conf"
	"github.com/sugimotosyo/hellogo/model"
	bitly "github.com/zpnk/go-bitly"

	"github.com/labstack/echo"
)

/************************
ユーザー
*************************/

//TwitterHandler Twitterのinterface
type TwitterHandler interface {
	MiddleFunc(c echo.Context) error
	Post(c echo.Context) error
}

//Twitter Twitter
type Twitter struct {
	Body *model.TwitterRequestBody
}

//NewTwitter Twitterを生成
func NewTwitter() TwitterHandler {
	return &Twitter{}
}

// Post post
// @Tags twitter
// @Summary tweetする
// @Description tweetする
// @Accept  json
// @Produce  json
// @Param  body
// @Success 200 {array} model.Twitter
// @Failure 400 "BadRequest"
// @Failure 401 "Unauthorized"
// @Router /twitter [post]
func (i *Twitter) Post(c echo.Context) error {

	api := getTwitterAPI(i.Body.Data.Token, i.Body.Data.Secret)

	// ハッシュタグを生成
	hashTagStr := "#バチェラー3 #bachelor-card"

	//文章
	sentence := i.Body.Data.Sentence

	//url
	host := conf.BachelorHost
	entry := "?key="
	addURL, err := shortURL(host + entry + i.Body.Data.Key)
	if err != nil {
		fmt.Println(err)
	}

	tweetStr := fmt.Sprintf("%s\r\n%s\r\n↓↓↓共有URL↓↓↓\r\n%s\r\n", hashTagStr, sentence, addURL)

	// post
	twt, err := api.PostTweet(tweetStr, nil)
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, twt)
}

/************************
global
*************************/

/************************
local
*************************/

//getTwitterAPI .
func getTwitterAPI(token, secret string) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	return anaconda.NewTwitterApi(token, secret)
}

//shortURL .
func shortURL(urlStr string) (string, error) {
	b := bitly.New(conf.BitlyToken)

	link, err := b.Links.Shorten(urlStr)
	if err != nil {
		return urlStr, err
	}
	fmt.Println(link.URL)
	return link.URL, nil
}

//MiddleFunc bodyとログインユーザとそのロール取得する
func (i *Twitter) MiddleFunc(c echo.Context) error {
	//bodyを取得
	body := new(model.TwitterRequestBody)
	if err := c.Bind(body); err != nil {
		fmt.Println("err")
		return errors.WithStack(err)
	}
	i.Body = body

	return nil
}
