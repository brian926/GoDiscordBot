package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type Answers struct {
	OriginChannelId string
	FavFood         string
	FavGame         string
}

var responses map[string]Answers = map[string]Answers{}

const prefix string = "!gobot"

func (a *Answers) ToMessageEmbed() discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "Fav",
			Value: a.FavFood,
		},
		{
			Name:  "Fav game",
			Value: a.FavGame,
		},
	}

	return discordgo.MessageEmbed{
		Title:  "New responses!",
		Fields: fields,
	}
}

func UserPromptHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// user channel
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Panic(err)
	}

	// if the user is already answers questions, ignore it, otherwise ask questions
	if _, ok := responses[channel.ID]; !ok {
		responses[channel.ID] = Answers{
			OriginChannelId: m.ChannelID,
			FavFood:         "",
			FavGame:         "",
		}
		s.ChannelMessageSend(channel.ID, "Hey there! Here are some questions")
	} else {
		s.ChannelMessageSend(channel.ID, "We're still waiting... ")
	}
}

func main() {
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")

	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		// DM logic
		if m.GuildID == "" {
			answers, ok := responses[m.ChannelID]
			if !ok {
				return
			}

			if answers.FavFood == "" {
				answers.FavFood = m.Content

				s.ChannelMessageSend(m.ChannelID, "Nice what fav game")

				responses[m.ChannelID] = answers
				return
			} else {
				answers.FavGame = m.Content
				embed := answers.ToMessageEmbed()
				s.ChannelMessageSendEmbed(answers.OriginChannelId, &embed)

				delete(responses, m.ChannelID)
			}
		}

		// server logic
		args := strings.Split(m.Content, " ")

		if args[0] != prefix {
			return
		}

		if args[1] == "hello" {
			s.ChannelMessageSend(m.ChannelID, "world!")
		}

		if args[1] == "proverbs" {
			proverbs := []string{
				"Don't communicate by sharing memory, share memory by communicating.",
				"Concurrency is not parallelism.",
				"Channels orchestrate; mutexes serialize.",
				"The bigger the interface, the weaker the abstraction.",
				"Make the zero value useful.",
				"interface{} says nothing.",
				"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
				"A little copying is better than a little dependency.",
				"Syscall must always be guarded with build tags.",
				"Cgo must always be guarded with build tags.",
				"Cgo is not Go.",
				"With the unsafe package there are no guarantees.",
				"Clear is better than clever.",
				"Reflection is never clear.",
				"Errors are values.",
				"Don't just check errors, handle them gracefully.",
				"Design the architecture, name the components, document the details.",
				"Documentation is for users.",
				"Don't panic.",
			}

			selection := rand.Intn(len(proverbs))

			author := discordgo.MessageEmbedAuthor{
				Name: "Rob Pike",
				URL:  "https://go-proverbs.github.io",
			}

			embed := discordgo.MessageEmbed{
				Title:  proverbs[selection],
				Author: &author,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		}

		if args[1] == "minecraft" {
			res := minecraftCheck()

			embed := res.ToMessageEmbed()
			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		}

		if args[1] == "prompt" {
			UserPromptHandler(s, m)
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	fmt.Println("--> Bot is Online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
