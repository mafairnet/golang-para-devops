package main

import (
	"flag"
	"fmt"
	"strings"

	_ "github.com/bluele/slack"
)

var configuration = getProgramConfiguration()

func main() {

	textPtr := flag.String("text", "", "Text to parse.")
	flag.Parse()
	message := *textPtr

	channels := strings.Split(configuration.SlackChannels, ",")

	for _, channel := range channels {
		fmt.Printf("Sending '%v' to: %v\n", message, channel)
		sendMessage(message, channel)
	}
}
