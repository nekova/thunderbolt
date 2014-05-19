package main

import (
	"fmt"
	"github.com/k0kubun/go-readline"
	"log"
	"regexp"
	"strings"
)

func executeCommand(account *Account, line string) {
	streamBlocked = true
	defer func() { streamBlocked = false }()

	if !strings.HasPrefix(line, ":") {
		confirmTweet(account, line)
		return
	}

	command, argument := splitCommand(line)
	switch command {
	case "recent":
		recent(account, argument)
	case "mentions":
		mentionsTimeline(account)
	case "favorite":
		confirmFavorite(account, argument)
	case "retweet":
		confirmRetweet(account, argument)
	default:
		commandNotFound()
	}
}

func regexpMatch(text string, exp string) bool {
	re, err := regexp.Compile(exp)
	if err != nil {
		log.Fatal(err)
	}
	return re.MatchString(text)
}

func splitCommand(text string) (string, string) {
	re, err := regexp.Compile("^:[^ ]+")
	if err != nil {
		log.Fatal(err)
	}

	result := re.FindStringIndex(text)
	if result == nil {
		return text[1:], ""
	}
	last := result[1]

	if last+1 >= len(text) {
		return text[1:], ""
	}
	return text[1:last], text[last+1:]
}

func confirmTweet(account *Account, text string) {
	confirmExecute(func() error {
		return updateStatus(account, text)
	}, "update '%s'", text)
}

func confirmFavorite(account *Account, argument string) {
	address := extractAddress(argument)
	if address == "" {
		commandNotFound()
		return
	}

	tweet := tweetMap.tweetByAddress(address)
	if tweet == nil || tweet.Id == 0 {
		println("Tweet is not registered")
		return
	}

	confirmExecute(func() error {
		return favorite(account, tweet)
	}, "favorite '%s'", tweet.Text)
}

func confirmRetweet(account *Account, argument string) {
	address := extractAddress(argument)
	if address == "" {
		commandNotFound()
		return
	}

	tweet := tweetMap.tweetByAddress(address)
	if tweet == nil || tweet.Id == 0 {
		println("Tweet is not registered")
		return
	}

	confirmExecute(func() error {
		return retweet(account, tweet)
	}, "retweet '%s'", tweet.Text)
}

func confirmExecute(function func() error, format string, a ...interface{}) {
	confirmMessage := fmt.Sprintf(format, a...)

	for {
		fmt.Println(foreColoredText(confirmMessage, "red"))

		answer := excuse("[Yn] ")
		if answer == "Y" || answer == "y" || answer == "" {
			err := function()
			if err != nil {
				print(err)
			}
			return
		} else if answer == "N" || answer == "n" {
			return
		}
	}
}

func excuse(prompt string) string {
	result := readline.Readline(&prompt)
	if result == nil {
		print("\n")
		return "n"
	}
	return *result
}

func extractAddress(argument string) string {
	re, err := regexp.Compile("\\$[a-z][a-z]")
	if err != nil {
		log.Fatal(err)
	}

	result := re.FindString(argument)
	if result == "" {
		return ""
	} else {
		return result[1:]
	}
}

func commandNotFound() {
	fmt.Printf("%s\n", backColoredText(foreBlackText("Command not found"), "yellow"))
}
