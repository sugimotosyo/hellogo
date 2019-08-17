package conf

import (
	"os"
)

//app設定
const (
	ServiceName = "サービス名"

	//Port
	Port = ":8080"
	//SSLPort sslのポート
	SSLPort = ":443"

	DateFormat   = "2006/01/02 15:04"
	BachelorHost = "https://bachelor.enjoy-ps.com/"
)

//本番と開発で書き換えが必要な値のデフォルトは開発環境にしておく
var (
	//アプリのホスト
	AppHost    = "http://localhost"
	BitlyToken = os.Getenv("BITLY_TOKEN")
)

func init() {

	if os.Getenv("ENVIRONMENT") == "prod" { //本番環境
		AppHost = "http://localhost"
	} else if os.Getenv("ENVIRONMENT") == "master" { //デモ環境
		AppHost = "http://localhost"
	} else if os.Getenv("ENVIRONMENT") == "develop" { //開発環境
		AppHost = "http://localhost"
	}
}
