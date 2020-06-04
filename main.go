package main

import (
	"fmt"
	"log"

	"github.com/FactomProject/factom"
	monitor "github.com/WhoSoup/factom-monitor"
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

	mon, err := monitor.NewMonitor(factomd)
	factom.SetFactomdServer(factomd)
	if err != nil {
		log.Fatal(err)
	}

	_, height, min := mon.GetCurrentMinute()
	fmt.Println("Initializing at height", height, "minute", min, "with factomd endpoint", factomd)

	reporter := new(Reporter)
	reporter.Height = height
	reporter.monitor = mon
	reporter.discord = NewDiscordHook(url)

	if name := cfg.Section("reporter").Key("name").String(); name != "" {
		reporter.discord.Name = name
	}
	if av := cfg.Section("reporter").Key("avatar").String(); av != "" {
		reporter.discord.Avatar = av
	}
	reporter.Run()
}
