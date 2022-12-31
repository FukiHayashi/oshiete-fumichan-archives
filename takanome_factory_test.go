package main_test

import (
	"fmt"
	"takanome/models"
	"time"

	"github.com/bluele/factory-go/factory"
)

var tweetFactory = factory.NewFactory(
	&models.Tweet{TweetedAt: time.Now(), RawData: "rawdata"},
).SeqInt("ID", func(n int) (interface{}, error) {
	return int64(n), nil
}).Attr("Text", func(args factory.Args) (interface{}, error) {
	tweet := args.Instance().(*models.Tweet)
	return fmt.Sprintf("tweet-%d", tweet.ID), nil
}).Attr("Url", func(args factory.Args) (interface{}, error) {
	tweet := args.Instance().(*models.Tweet)
	return fmt.Sprintf("http://localhost/tweet/%d", tweet.ID), nil
}).Attr("ScreenName", func(args factory.Args) (interface{}, error) {
	tweet := args.Instance().(*models.Tweet)
	return fmt.Sprintf("user-%d", tweet.ID), nil
})
