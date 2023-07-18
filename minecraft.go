package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/bwmarrin/discordgo"
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

const minecraftUrl string = "https://api.mcsrvstat.us/bedrock/2/bytez.us"

func minecraftCheck() Response {

	resp, err := http.Get(minecraftUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return result
}

func (res *Response) ToMessageEmbed() discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:  "MOTD",
			Value: res.Motd.Clean[0],
		},
		{
			Name:  "IP",
			Value: res.IP,
		},
		{
			Name:  "Online",
			Value: strconv.FormatBool(res.Online),
		},
		{
			Name:  "Players Online",
			Value: strconv.Itoa(res.Players.Online),
		},
		{
			Name:  "Map",
			Value: res.Map,
		},
	}

	return discordgo.MessageEmbed{
		Title:  res.Hostname,
		Fields: fields,
	}
}
