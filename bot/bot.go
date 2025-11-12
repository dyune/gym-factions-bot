package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/davidwang/factions/internal/handlers"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "register",
			Description: "Sign up for factions and begin logging your gym points!.",
		},
		{
			Name:        "basic-command",
			Description: "Basic command",
		},
		{
			Name:        "responses",
			Description: "Interaction responses testing initiative",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "resp-type",
					Description: "Response type",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Channel message with source",
							Value: 4,
						},
						{
							Name:  "Deferred response With Source",
							Value: 5,
						},
					},
					Required: true,
				},
			},
		},
	}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"register":      handlers.HandleRegisterUser,
		"basic-command": BasicCommandHandler,
		"responses":     ResponseHandler,
	}
	Token             = os.Getenv("DISCORD_TOKEN")
	s                 *discordgo.Session
	RemoveCommands, _ = strconv.ParseBool(os.Getenv("REMOVE_CMDS"))
	GuildId           = os.Getenv("SERVER_ID")
)

func Run() {
	var err error
	s, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("%s", err)
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	log.Printf("GymFactions Bot | v.0.0.0")

	s.AddHandler(newMessage)
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("bot.go/Open: %s", err)
	}
	defer func(s *discordgo.Session) {
		err := s.Close()
		if err != nil {
			log.Fatalf("bot.go/Close: could not close gracefully: %s", err)
		}
	}(s)

	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))

	for i, v := range Commands {
		var cmd *discordgo.ApplicationCommand
		cmd, err = s.ApplicationCommandCreate(s.State.User.ID, GuildId, v)
		if err != nil {
			log.Panicf("bot.go/Run: failed to register command: %s", err)
		}
		registeredCommands[i] = cmd
	}

	log.Printf("Running...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // block and wait

	if RemoveCommands {
		for _, v := range registeredCommands {
			err = s.ApplicationCommandDelete(s.State.User.ID, GuildId, v.ID)
			if err != nil {
				log.Panicf("bot.go/Run: while removing Commands for shutdown: %s", err)
			}
		}
	}

	log.Printf("Gracefully shutting bot down.")
}

func newMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == s.State.User.ID {
		return
	}
	log.Printf("%v, %v, %v", msg.Author.ID, msg.Author.GlobalName, msg.Author.Email)
	var err error
	switch {
	case strings.Contains(msg.Content, "!help"):
		_, err = s.ChannelMessageSend(msg.ChannelID, "Hello World")
	case strings.Contains(msg.Content, "!bye"):
		_, err = s.ChannelMessageSend(msg.ChannelID, "Good Bye")
	}
	if err != nil {
		log.Fatalf("bot.go/newMessage: %s", err)
	}
}

func BasicCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hey there! Congratulations, you just executed your first slash command",
		},
	})
}

func ResponseHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Responses to a command are very important.
	// First of all, because you need to react to the interaction
	// by sending the response in 3 seconds after receiving, otherwise
	// interaction will be considered invalid and you can no longer
	// use the interaction token and ID for responding to the user's request

	content := ""
	// As you can see, the response type names used here are pretty self-explanatory,
	// but for those who want more information see the official documentation
	switch i.ApplicationCommandData().Options[0].IntValue() {
	case int64(discordgo.InteractionResponseChannelMessageWithSource):
		content =
			"You just responded to an interaction, sent a message and showed the original one. " +
				"Congratulations!"
		content +=
			"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
	default:
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
		})
		if err != nil {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong",
			})
		}
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
		return
	}

	time.AfterFunc(time.Second*5, func() {
		content := content + "\n\nWell, now you know how to create and edit responses. " +
			"But you still don't know how to delete them... so... wait 10 seconds and this " +
			"message will be deleted."
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		if err != nil {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong",
			})
			return
		}
		time.Sleep(time.Second * 10)
		s.InteractionResponseDelete(i.Interaction)

	})
}
