package handlers

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davidwang/factions/internal/config"
	"github.com/davidwang/factions/internal/exceptions"
	"github.com/davidwang/factions/internal/store"
	"log"
	"strconv"
)

func HandleRegisterUser(s *discordgo.Session, i *discordgo.InteractionCreate) {

	msg := ""
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		log.Printf("Error fetching channel: %v", err)
		return
	}

	print(channel.Name)
	// Use member since messages will be sent from a server and not DMs
	user := i.Member.User
	log.Printf("User invoked /register")
	id, err := strconv.Atoi(user.ID)

	if err != nil {
		log.Printf("users.go/HandleRegisterUser: %v", err)
		msg = "Something went wrong during registration."
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		},
		)
		return
	}

	_, err = store.InsertUser(id, user.GlobalName, config.GlobalCtx, config.DB)

	if errors.Is(err, exceptions.ErrUserExists) {
		log.Printf("user exists: %v", err)
		msg = "You're already registered :D"
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return

	} else if errors.Is(err, exceptions.ErrDatabaseOp) {
		log.Printf("database operation error: %v", err)
		msg = "Database operation error."
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return

	} else if err != nil {
		log.Printf("unexpected error: %v", err)
		msg = "An unexpected error has occurred."
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		return
	}

	err = HandleAssignFaction(id, channel.Name)

	msg = "You have successfully been registered!"
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})

	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
	}
	return
}

func HandleAssignFaction(id int, channelName string) error {
	factionName, ok := store.FactionNames[channelName]
	if ok == false {
		return errors.New(fmt.Sprintf("No such faction with this channel name: %s", channelName))
	}

	factionId, err := store.LookUpFaction(config.GlobalCtx, config.DB, factionName)
	if err != nil {
		return err
	}

	ok, err = store.ExistsByID(id, config.GlobalCtx, config.DB)

	if err != nil {
		return err
	}

	if ok {
		err = store.SetUserDetailsByID(id, nil, nil, &factionId, config.GlobalCtx, config.DB)
	} else {
		err = exceptions.UserNotFoundError{UserID: id}
	}
	if err != nil {
		return err
	}
	return nil
}

func registerUser(id int, username string) error {
	return nil
}
