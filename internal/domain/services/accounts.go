package services

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/davidwang/factions/internal/config"
	"github.com/davidwang/factions/internal/domain/exceptions"
	store2 "github.com/davidwang/factions/internal/infra/postgres"
)

func HandleRegisterAccount(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		log.Printf("Error fetching channel: %v", err)
		return
	}

	// Use member since messages will be sent from a server and not DMs
	Account := i.Member.User
	log.Printf("Account invoked /register")
	id, err := strconv.Atoi(Account.ID)

	if err != nil {
		log.Printf("Accounts.go/HandleRegisterAccount: %v", err)
		respond("Something went wrong during registration.", s, i)
		return
	}

	factionName, ok := store2.FactionNames[channel.Name]
	if !ok {
		respond("Call this register command in your faction channel, please.", s, i)
		return
	}

	faction, err := store2.LookUpFaction(config.GlobalCtx, config.DB, factionName)
	if errors.Is(err, sql.ErrNoRows) {
		errRespond("This faction... does not exist.", s, i)
		return
	} else if errors.Is(err, exceptions.ErrDatabaseOp) {
		respond("While trying to check if you called this command properly, something went wrong.", s, i)
		return
	} else if err != nil {
		errRespond("while checking for factions", s, i)
		return
	}

	_, err = store2.InsertAccount(id, Account.GlobalName, faction.ID, config.GlobalCtx, config.DB)

	if errors.Is(err, exceptions.ErrAccountExists) {
		log.Printf("Account exists: %v", err)
		respond("You're already registered to this faction :)", s, i)
		return
	} else if errors.Is(err, exceptions.ErrDatabaseOp) {
		log.Printf("database operation error: %v", err)
		errRespond("the database could not register you.", s, i)
		return
	} else if err != nil {
		log.Printf("unexpected error: %v", err)
		respond("unknown error occurred while the database tried to register ou.", s, i)
		return
	}

	respond("You have successfully been registered!", s, i)

	return
}
