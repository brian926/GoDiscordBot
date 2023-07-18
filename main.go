package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Created using https://mholt.github.io/json-to-go/
type Response struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Motd struct {
		Raw   []string `json:"raw"`
		Clean []string `json:"clean"`
		HTML  []string `json:"html"`
	} `json:"motd"`
	Players struct {
		Online int `json:"online"`
		Max    int `json:"max"`
	} `json:"players"`
	Version  string `json:"version"`
	Online   bool   `json:"online"`
	Hostname string `json:"hostname"`
	Map      string `json:"map"`
	Gamemode string `json:"gamemode"`
	Serverid string `json:"serverid"`
}

const prefix string = "!gobot"
const minecraftUrl string = "https://api.mcsrvstat.us/bedrock/2/bytez.us"

func minecraftCheck() Response {

	resp, err := http.Get(minecraftUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return result
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

			des := "MOTD: " + res.Motd.Clean[0] + "\nIP: " + res.IP + "\nOnline: " + strconv.FormatBool(res.Online) + "\nPlayers Online: " + strconv.Itoa(res.Players.Online) + "\nMap: " + res.Map

			embed := discordgo.MessageEmbed{
				Title:       res.Hostname,
				Description: des,
			}

			s.ChannelMessageSendEmbed(m.ChannelID, &embed)
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
