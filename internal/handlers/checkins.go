package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davidwang/factions/internal/config"
	"github.com/davidwang/factions/internal/store"
	"strconv"
	"time"
)

func HandleSubmission(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		respond("An error occurred: channel retrieval created an error", s, i)
		return
	}
	if channel == nil {
		respond("An error occurred: channel is nil", s, i)
		return
	}

	id, err := strconv.Atoi(i.Member.User.ID)
	if err != nil {
		respond(fmt.Sprintf("An error occurred: id %s could not be parsed", i.Member.User.ID), s, i)
		return
	}

	exists, err := store.AccountExistsByID(id, config.GlobalCtx, config.DB)
	if !exists {
		respond("You aren't registered. Register to start earning points! Use the /register command.", s, i)
		return
	}
	if err != nil {
		respond("An error occurred: could not call database", s, i)
	}

	var faction *store.Faction
	faction, err = store.GetAccountFactionByID(id, config.GlobalCtx, config.DB)
	if err != nil {
		respond(fmt.Sprintf("An error occurred: %v", err), s, i)
		return
	}
	if faction.Name != store.FactionNames[channel.Name] {
		respond("Please upload activities in your own gym factions channel", s, i)
		return
	}

	err = store.InsertSub(time.Now(), 1, id, 1, config.GlobalCtx, config.DB)

	if err != nil {
		respond("Error registering your submission", s, i)
		return

	} else {
		respond("Successfully registered your submission! You've earned some points!", s, i)
	}
}
