package discordbooter

import (
	"context"
	"errors"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var ds *discordgo.Session

// Start a new discord bot - only one session is possible
func Start(ctx context.Context, wg *sync.WaitGroup, token string) error {

	// we only allow one initialization
	if ds != nil {
		return errors.New("bot is already booted")
	}

	var err error
	ds, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	ds.Identify.Intents = discordgo.IntentsGuildMessages

	err = ds.Open()
	if err != nil {
		ds.Close()
		return err
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer ds.Close()

		<-ctx.Done()
	}(wg)
	return nil
}

// AddHandlers to consume messages sent or received in discord
func AddHandlers(handlers []interface{}) error {
	// return if we don't have an active bot session
	if ds == nil {
		return errors.New("bot is not booted")
	}

	for _, handler := range handlers {
		ds.AddHandler(handler)
	}

	return nil
}

// SendMessage to a specific discord channel
func SendMessage(channel string, message string) error {
	// return if we don't have an active bot session
	if ds == nil {
		return errors.New("bot is not booted")
	}

	_, err := ds.ChannelMessageSend(channel, message)
	if err != nil {
		return err
	}
	return nil
}
