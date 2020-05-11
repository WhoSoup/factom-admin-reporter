package main

import (
	"fmt"
	"log"

	"github.com/FactomProject/factom"
	"github.com/go-ini/ini"
)

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal(err)
	}

	url := cfg.Section("reporter").Key("webhook").String()
	if url == "" {
		log.Fatal("ini field \"webhook\" is empty")
	}

	factomd := cfg.Section("reporter").Key("factomd").String()
	if factomd == "" {
		log.Fatal("ini field \"factomd\" is empty")
	}
	factom.SetFactomdServer(factomd)
	min, err := factom.GetCurrentMinute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initializing at height", min.DirectoryBlockHeight, "minute", min.Minute)

	reporter := new(Reporter)
	reporter.discord = NewDiscordHook(url)
	if name := cfg.Section("reporter").Key("name").String(); name != "" {
		reporter.discord.Name = name
	}
	if av := cfg.Section("reporter").Key("avatar").String(); av != "" {
		reporter.discord.Avatar = av
	}
	reporter.Run()
}
