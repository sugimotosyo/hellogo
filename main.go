package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World ")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/my_tweet", MyTweet)
	http.ListenAndServe(":8080", nil)
}

func MyTweet(w http.ResponseWriter, r *http.Request) {
	api := getTwitterApi()

	v := url.Values{}
	v.Set("count", "30")

	searchResult, _ := api.GetSearch("sugimotosyo", v)
	for _, tweet := range searchResult.Statuses {
		fmt.Println("-----------------------")
		// fmt.Println(tweet)
		fmt.Println(tweet.Text)
		// str, _ := json.Marshal(&tweet)
		// fmt.Println(string(str))
	}

	// text := "TEST Hello from API."
	// twt, err := api.PostTweet(text, nil)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(twt)

	fmt.Fprintf(w, "MyTweet "+os.Getenv("TWITTER_CONSUMER_KEY")

}

func getTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	return anaconda.NewTwitterApi(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
}
