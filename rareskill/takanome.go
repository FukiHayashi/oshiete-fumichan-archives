package rareskill

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func Takanome() {
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMER_KEY"), os.Getenv("TWITTER_SECRET_KEY"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_TOKEN_SECRET"))

	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// twitter client
	client := twitter.NewClient(httpClient)

	params := twitter.FavoriteListParams{
		ScreenName: os.Getenv("TWITTER_ACCOUNT"),
	}
	tlist, _, _ := client.Favorites.List(&params)

	jsonData, _ := json.Marshal(tlist[0])

	log.Printf("%s\n", jsonData)
}
