package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davidwang/factions/internal/config"
	"github.com/davidwang/factions/internal/exceptions"
	"github.com/davidwang/factions/internal/store"
	"log"
	"strconv"
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

	factionName, ok := store.FactionNames[channel.Name]
	if !ok {
		respond("Call this register command in your faction channel, please.", s, i)
		return
	}

	faction, err := store.LookUpFaction(config.GlobalCtx, config.DB, factionName)
	if errors.Is(err, sql.ErrNoRows) {
		respond("This faction... does not exist.", s, i)
		return
	} else if errors.Is(err, exceptions.ErrDatabaseOp) {
		respond("While trying to check if you called this command properly, something went wrong.", s, i)
		return
	} else if err != nil {
		respond("An unexpected error occurred while trying to register you.", s, i)
		return
	}

	_, err = store.InsertAccount(id, Account.GlobalName, faction.ID, config.GlobalCtx, config.DB)

	if errors.Is(err, exceptions.ErrAccountExists) {
		log.Printf("Account exists: %v", err)
		respond("You're already registered to this faction :)", s, i)
		return
	} else if errors.Is(err, exceptions.ErrDatabaseOp) {
		log.Printf("database operation error: %v", err)
		respond("An error with the database occurred while trying to register you.", s, i)
		return
	} else if err != nil {
		log.Printf("unexpected error: %v", err)
		respond("An unexpected error occurred while trying to register you.", s, i)
		return
	}

	respond("You have successfully been registered!", s, i)

	return
}

func HandleAssignFaction(id int, channelName string) error {
	factionName, ok := store.FactionNames[channelName]
	if ok == false {
		return errors.New(fmt.Sprintf("No such faction with this channel name: %s", channelName))
	}

	faction, err := store.LookUpFaction(config.GlobalCtx, config.DB, factionName)
	if err != nil {
		return err
	}

	ok, err = store.AccountExistsByID(id, config.GlobalCtx, config.DB)

	if err != nil {
		return err
	}

	if ok {
		err = store.SetAccountDetailsByID(id, nil, nil, faction, config.GlobalCtx, config.DB)
	} else {
		err = exceptions.AccountNotFoundError{AccountID: id}
	}
	if err != nil {
		return err
	}
	return nil
}
