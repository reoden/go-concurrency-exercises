//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

// import (
// 	"fmt"
// 	"time"
// )

// var (
// 	done     chan struct{}
// 	tweetsCh chan *Tweet
// )

// func producer(stream Stream) {
// 	for {
// 		tweet, err := stream.Next()
// 		if err == ErrEOF {
// 			close(tweetsCh)
// 			done <- struct{}{}
// 			return
// 		}

// 		tweetsCh <- tweet
// 		// tweets = append(tweets, tweet)
// 	}
// }

// func consumer(t *Tweet) {
// 	// for _, t := range tweets {
// 	if t.IsTalkingAboutGo() {
// 		fmt.Println(t.Username, "\ttweets about golang")
// 	} else {
// 		fmt.Println(t.Username, "\tdoes not tweet about golang")
// 	}
// 	// }
// }

// func main() {
// 	start := time.Now()
// 	stream := GetMockStream()

// 	done = make(chan struct{})
// 	tweetsCh = make(chan *Tweet)
// 	// Producer
// 	go producer(stream)

// 	for {
// 		select {
// 		case <-done:
// 			fmt.Printf("Process took %s\n", time.Since(start))
// 			return
// 		case t, ok := <-tweetsCh:
// 			// Consumer
// 			if ok {
// 				consumer(t)
//			}
// 		}
// 	}

// }
