package main

import (
	"fmt"
	"log"
	"time"

	"github.com/FactomProject/factom"
	monitor "github.com/WhoSoup/factom-monitor"
)

type Reporter struct {
	Height  int64
	monitor *monitor.Monitor
	discord *DiscordHook
}

func (r *Reporter) Run() {
	go func() {
		for err := range r.monitor.NewErrorListener() {
			log.Printf("[monitor error] %v", err)
		}
	}()
	for ev := range r.monitor.NewDBHeightListener() {
		log.Printf("trying block %d", ev)
		r.newBlock(ev)
	}
}

func (r *Reporter) newBlock(height int64) {
	var ablock *factom.ABlock
	var err error

	start := time.Now()
	report := 12 // every minute @ 5 seconds
	count := 0
	for {
		ablock, _, err = factom.GetABlockByHeight(height)
		if err != nil {
			count++
			if count%report == 0 {
				elapsed := time.Since(start)
				r.discord.SendMessage(fmt.Sprintf("trying to retrieve admin block %d unsuccessfully for %s: %v", height, elapsed, err))
			}
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

	r.Height = ablock.DBHeight
	log.Printf("ABlock[%d] %d entries", ablock.DBHeight, len(ablock.ABEntries))
	for i, e := range ablock.ABEntries {
		r.entry(i, e)
	}
}

func (r *Reporter) entry(i int, e factom.ABEntry) {
	field := func(a string, b interface{}) string {
		return fmt.Sprintf("\n\t`%s` = `%v`", a, b)
	}

	var msg string
	switch e.(type) {
	case *factom.AdminMinuteNumber: // deprecated
	case *factom.AdminDBSignature: // deprecated
	case *factom.AdminRevealHash:
		entry := e.(*factom.AdminRevealHash)
		msg = fmt.Sprint("Matryoshka Reveal Hash:", field("IdentityChain", entry.IdentityChainID), field("Matryoshka Hash", entry.MatryoshkaHash))
	case *factom.AdminAddHash:
		entry := e.(*factom.AdminAddHash)
		msg = fmt.Sprint("Matryoshka Add Hash:", field("IdentityChain", entry.IdentityChainID), field("Matryoshka Hash", entry.MatryoshkaHash))
	case *factom.AdminIncreaseServerCount:
		entry := e.(*factom.AdminIncreaseServerCount)
		msg = fmt.Sprintf("Increase Server Count (Deprecated): count = %d", entry.Amount)
	case *factom.AdminAddFederatedServer:
		entry := e.(*factom.AdminAddFederatedServer)
		msg = fmt.Sprint("AddFederatedServer:", field("Activation Height", entry.DBHeight), field("Identity Chain", entry.IdentityChainID))
	case *factom.AdminAddAuditServer:
		entry := e.(*factom.AdminAddAuditServer)
		msg = fmt.Sprint("AddAuditServer:", field("Activation Height", entry.DBHeight), field("Identity Chain", entry.IdentityChainID))
	case *factom.AdminRemoveFederatedServer:
		entry := e.(*factom.AdminRemoveFederatedServer)
		msg = fmt.Sprint("RemoveFederatedServer:", field("Activation Height", entry.DBHeight), field("Identity Chain", entry.IdentityChainID))
	case *factom.AdminAddFederatedServerKey:
		entry := e.(*factom.AdminAddFederatedServerKey)
		msg = fmt.Sprint("AddFederatedServerKey:", field("Activation Height", entry.DBHeight), field("Identity Chain", entry.IdentityChainID), field("Priority", entry.KeyPriority), field("Public Key", entry.PublicKey))
	case *factom.AdminAddFederatedServerBTCKey:
		entry := e.(*factom.AdminAddFederatedServerBTCKey)
		msg = fmt.Sprint("AddFederatedServerBTCKey:", field("Identity Chain", entry.IdentityChainID), field("Priority", entry.KeyPriority), field("Key Type", entry.KeyType), field("ECDSA Public Key", entry.ECDSAPublicKey))
	case *factom.AdminServerFault:
		entry := e.(*factom.AdminServerFault)
		msg = fmt.Sprint("ServerFault (Deprecated):", field("Timestamp", entry.Timestamp), field("Server ID", entry.ServerID), field("Audit Server ID", entry.AuditServerID), field("VMIndex", entry.VMIndex), field("DBHeight", entry.DBHeight), field("Height", entry.Height))
	case *factom.AdminCoinbaseDescriptor:
		entry := e.(*factom.AdminCoinbaseDescriptor)

		payout := 0
		for _, add := range entry.Outputs {
			payout += add.Amount
		}

		fct := payout / 1e8
		rest := payout % 1e8
		msg = fmt.Sprint("Coinbase Descriptor: ", len(entry.Outputs), " outputs, fct = ", fmt.Sprintf("%d.%d", fct, rest))
	case *factom.AdminCoinbaseDescriptorCancel:
		entry := e.(*factom.AdminCoinbaseDescriptorCancel)
		msg = fmt.Sprint("CoinbaseDescriptorCancel:", field("Height", entry.DescriptorHeight), field("Index", entry.DescriptorIndex))
	case *factom.AdminAddAuthorityAddress:
		entry := e.(*factom.AdminAddAuthorityAddress)
		msg = fmt.Sprint("AddAuthorityAddress:", field("Identity Chain", entry.IdentityChainID), field("Factoid Address", entry.FactoidAddress))
	case *factom.AdminAddAuthorityEfficiency:
		entry := e.(*factom.AdminAddAuthorityEfficiency)
		msg = fmt.Sprint("AddAuthorityEfficiency:", field("Identity Chain", entry.IdentityChainID), field("Efficiency", fmt.Sprintf("%.2f%%", float64(entry.Efficiency)/100)))
	}

	if msg != "" {
		r.discord.SendMessage(fmt.Sprintf("ABlock[height=%d,entry=%d] %s", r.Height, i, msg))
	}
}
