package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zaskoh/discordbooter"
)

var wg sync.WaitGroup

var (
	token   = flag.String("token", "", "bot token")
	channel = flag.String("channel", "", "discord channel the bot sends messages to")
)

func main() {
	flag.Parse()
	ctx, cancelFunc, cancelChan := CreateLaunchContext()
	defer cancelFunc()

	log.Println("starting bot")
	err := discordbooter.Start(ctx, &wg, *token)
	if err != nil {
		log.Fatalf("starting failed %s", err)
	}
	log.Printf("bot started and is ready to get handlers")

	var handlers = []interface{}{
		// just write in the channel that the bot started
		func(s *discordgo.Session, r *discordgo.Ready) {
			log.Println("discordbooter started new session!")
		},
		// add the respond handler
		messageResponder,
	}
	err = discordbooter.AddHandlers(handlers)
	if err != nil {
		log.Fatalf("unable to add handlers: %s", err)
	}

	log.Printf("bot started and running in background now")

	discordbooter.SendMessage(*channel, "discordbooter successfull started new session and is waiting for you!")

	<-cancelChan
	wg.Wait()
}

func messageResponder(s *discordgo.Session, m *discordgo.MessageCreate) {
	// we skip messages the bot did
	if m.Author.ID == s.State.User.ID {
		return
	}

	// we only want to consume messages from the defined channel
	if m.Message.ChannelID != *channel {
		return
	}

	// example of how we can just write stuff to the channel
	_, err := s.ChannelMessageSend(*channel, "okokkok")
	if err != nil {
		log.Printf("error: write message %s with error %s", m.Content, err)
	}

	// this is the message we consumed
	log.Printf("%+v", m.Content)

	// example of how we can respond to the message
	_, err = s.ChannelMessageSendReply(*channel, "yo", m.Reference())
	if err != nil {
		log.Printf("error: cant reply to Message %s with error %s", m.Content, err)
	}
}

func CreateLaunchContext() (context.Context, func(), chan bool) {
	interruptChan := make(chan os.Signal, 1)
	canceledChanChan := make(chan bool, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		defer close(interruptChan)
		<-interruptChan
		cancelCtx()
		canceledChanChan <- true
	}()
	cancel := func() {
		cancelCtx()
		close(canceledChanChan)
	}
	return ctx, cancel, canceledChanChan
}
