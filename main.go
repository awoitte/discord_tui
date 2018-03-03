package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/awoitte/input_output_tui"
	"github.com/bwmarrin/discordgo"
	"github.com/kyokomi/emoji"
)

var (
	Config string
	UserID string
)

func init() {
	flag.StringVar(&Config, "c", "", "path to config file")
	flag.StringVar(&UserID, "u", "", "user to chat with")
	flag.Parse()
}

func main() {
	if Config == "" || UserID == "" {
		fmt.Println("usage: discord_tui -c <path to config> -u <user to chat with>")
		return
	}

	data, err := ioutil.ReadFile(Config)
	if err != nil {
		fmt.Println("error reading config file,", err)
		return
	}
	config := make(map[string]string)
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("error parsing config file,", err)
		return
	}

	dg, err := discordgo.New(config["username"], config["password"])
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	input := make(chan string)
	messages := make(chan string)
	quit := make(chan bool)

	dg.AddHandler(createMessageHandler(messages))

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	channel, err := dg.UserChannelCreate(config[UserID])
	if err != nil {
		fmt.Println("error opening channel,", err)
		return
	}

	go send_older_messages(messages, dg, channel)

	go send_input_to_discord(input, messages, dg, channel.ID)

	input_output_tui.Start(input, messages, quit)
}

func send_older_messages(messages chan string, dg *discordgo.Session, channel *discordgo.Channel) {
	message_list, err := dg.ChannelMessages(channel.ID, 25, "", "", channel.LastMessageID)
	if err != nil {
		messages <- fmt.Sprintln("error getting messages,", err)
	}

	for _, message := range reverse_message_order(message_list) {
		messages <- format_message(message.Author.Username, message.Content)
	}
}

func send_input_to_discord(input, messages chan string, dg *discordgo.Session, channelID string) {
	for {
		select {
		case received := <-input:
			with_emoji := emoji.Sprint(received)
			_, err := dg.ChannelMessageSend(channelID, with_emoji)
			if err != nil {
				messages <- fmt.Sprintln("error sending message,", err)
			} else {
				messages <- format_message(dg.State.User.Username, with_emoji)
			}
		case <-time.After(time.Second):
		}
	}
}

func createMessageHandler(output chan string) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		output <- format_message(m.Author.Username, m.Content)
	}
}

func format_message(username, content string) string {
	return username + ": " + content
}

func reverse_message_order(messages []*discordgo.Message) []*discordgo.Message {
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages
}
