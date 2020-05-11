package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const MSG_LIMIT = 20000

type DiscordHook struct {
	URL    string
	Name   string
	Avatar string

	client *http.Client

	mtx            sync.Mutex
	rateLimit      int
	rateRemain     int
	rateReset      int
	rateResetAfter int
}

func NewDiscordHook(url string) *DiscordHook {
	d := new(DiscordHook)
	d.URL = url
	d.rateRemain = 1
	d.client = &http.Client{}
	return d
}

func (d *DiscordHook) formMessage(content string) (map[string]string, error) {
	if content == "" {
		return nil, errors.New("empty message")
	}
	if len(content) > MSG_LIMIT {
		return nil, fmt.Errorf("message too long (%d of %d allowed)", len(content), MSG_LIMIT)
	}
	msg := make(map[string]string)
	msg["content"] = content
	if d.Name != "" {
		msg["username"] = d.Name
	}
	if d.Avatar != "" {
		msg["avatar_url"] = d.Avatar
	}

	return msg, nil
}

func (d *DiscordHook) waitUntilFree() error {
	if d.rateResetAfter > 300 {
		return fmt.Errorf("we were asked to wait a long time: %d seconds", d.rateResetAfter)
	}
	if d.rateRemain == 0 {
		time.Sleep(time.Second * time.Duration(d.rateResetAfter))
	}
	return nil
}

func (d *DiscordHook) SendMessage(content string) error {
	msg, err := d.formMessage(content)
	if err != nil {
		return err
	}

	js, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	d.mtx.Lock()
	defer d.mtx.Unlock()
	req, err := http.NewRequest("POST", d.URL, bytes.NewBuffer(js))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	if err := d.waitUntilFree(); err != nil {
		return err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d.handleResponse(resp)
	return nil
}
func (d *DiscordHook) handleResponse(resp *http.Response) {
	if resp.StatusCode != http.StatusNoContent {

		if resp.StatusCode == http.StatusTooManyRequests {
			log.Printf("we hit rate limit")
		} else {
			log.Printf("unexpected status code: %s", resp.Status)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response body:", string(body))
	}

	convert := func(key string) (int, error) {
		raw := resp.Header.Get(key)
		conv, err := strconv.Atoi(raw)
		if err != nil {
			return 0, fmt.Errorf("unable to convert %s (%s) to int: %v", key, raw, err)
		}
		return conv, nil
	}

	rl, err := convert("X-Ratelimit-Limit")
	if err != nil {
		log.Print(err)
	} else {
		d.rateLimit = rl
	}

	rr, err := convert("X-Ratelimit-Remaining")
	if err != nil {
		log.Print(err)
	} else {
		d.rateRemain = rr
	}

	rre, err := convert("X-Ratelimit-Reset")
	if err != nil {
		log.Print(err)
	} else {
		d.rateReset = rre
	}

	rra, err := convert("X-Ratelimit-Reset-After")
	if err != nil {
		log.Print(err)
	} else {
		d.rateResetAfter = rra
	}

}
