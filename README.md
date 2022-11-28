# discordbooter
<p>
    <a href="https://pkg.go.dev/github.com/zaskoh/discordbooter">
        <img alt="Go reference" src="https://img.shields.io/badge/reference-grey?style=flat-square&logo=Go">
    </a>
    <a href="https://github.com/zaskoh/discordbooter/actions/workflows/test.yml">
        <img alt="GitHub Workflow Status" src="https://github.com/zaskoh/discordbooter/workflows/Test/badge.svg?style=flat-square">
    </a>
    <a href="https://goreportcard.com/report/github.com/zaskoh/discordbooter">
        <img alt="Go Report Card" src="https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat-square">
    </a>
    <a href="https://github.com/zaskoh/discordbooter/blob/main/go.mod">
        <img alt="go version" src="https://img.shields.io/github/go-mod/go-version/zaskoh/discordbooter?style=flat-square&logo=Go">
    </a>
    <a href="https://github.com/zaskoh/discordbooter/blob/main/LICENSE">
        <img alt="license" src="https://img.shields.io/github/license/zaskoh/discordbooter?style=flat-square">
    </a>
    <a href="https://github.com/zaskoh/discordbooter/releases">
        <img alt="GitHub Release" src="https://img.shields.io/github/v/release/zaskoh/discordbooter?style=flat-square&include_prereleases&sort=semver">
    </a>
</p>

**discordbooter** is a package that extends [discordgo](https://github.com/bwmarrin/discordgo) to spin up a discord bot to log or consume messages from / to a discord channel.

## Getting started

### Installing
```bash
go get github.com/zaskoh/discordbooter
```

To get started you need a discord bot and a token from discord. Follow these articles to prepare it.
- https://discord.com/developers/docs/getting-started#creating-an-app
- https://discord.com/developers/docs/getting-started#configuring-a-bot
- add scope bot + all bot permissions needed (send message / ...) and invite bot (OAuth2) to your server

### Info
The bot will run as a goroutine. To handle the shutdown gracefully you need a context and a waitgroup.

An example how to use it can be found under examples.

### Example
```bash
go run example/example.go --token xxx --channel yyy
```

### Usage
#### Import the package into your project
```go
import "github.com/zaskoh/discordbooter"
```

#### Start the bot
```go
err := discordbooter.Start(ctx, &wg, *token)
if err != nil {
    log.Fatalf("booting failed %s", err)
}
```

#### Add handlers
```go
// define a function you want to add
func discord_handler_x(s *discordgo.Session, m *discordgo.MessageCreate){
	// code goes here
}
// define an array of handlers
var handlers = []interface{}{
    func(s *discordgo.Session, r *discordgo.Ready) { 
        log.Println("discord bot started") 
    }, 
    discord_handler_x,
}

// add them
err = discordbooter.AddHandlers(handlers)
if err != nil {
    log.Printf("unable to add handlers: %s", err)
}
```

#### Send message
```go
err := discordbooter.SendMessage(*channel, "hello world!")
if err != nil {
    log.Printf("unable to send message: %s", err)
}
```