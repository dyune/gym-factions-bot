package handlers

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davidwang/factions/internal/config"
	"github.com/davidwang/factions/internal/store"
	"github.com/uptrace/bun"
	"log"
	"strconv"
	"time"
)

func HandleSubmission(s *discordgo.Session, i *discordgo.InteractionCreate) {

	submissionTime := time.Now()

	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		err_respond("channel retrieval created an error", s, i)
		return
	}
	if channel == nil {
		err_respond("channel is nil", s, i)
		return
	}

	id, err := strconv.Atoi(i.Member.User.ID)
	if err != nil {
		err_respond(fmt.Sprintf("id %s could not be parsed", i.Member.User.ID), s, i)
		return
	}

	exists, err := store.AccountExistsByID(id, config.GlobalCtx, config.DB)
	if !exists {
		respond("You aren't registered. Register to start earning points! Use the /register command.", s, i)
		return
	}
	if err != nil {
		err_respond("could not call database", s, i)
	}

	var faction *store.Faction
	faction, err = store.GetAccountFactionByID(id, config.GlobalCtx, config.DB)
	if err != nil {
		err_respond(fmt.Sprintf("%v", err), s, i)
		return
	}
	if faction.Name != store.FactionNames[channel.Name] {
		respond("Please upload activities in your own gym factions channel", s, i)
		return
	}

	var reachedLimit bool
	reachedLimit, err = reachedSubmissionLimit(submissionTime, id, config.GlobalCtx, config.DB)
	if !reachedLimit {
		err = store.InsertSub(submissionTime, 1, id, 1, config.GlobalCtx, config.DB)
		if err != nil {
			err_respond(fmt.Sprintf("%v", err), s, i)
			return

		} else {
			respond("Successfully registered your submission! You've earned some points!", s, i)
		}
	}
	if err != nil {
		err_respond(fmt.Sprintf("%v", err), s, i)
	} else {
		respond("Can't submit more than one activity per day!", s, i)
		log.Printf("%v", submissionTime)
	}
}

func reachedSubmissionLimit(time time.Time, id int, ctx context.Context, db *bun.DB) (bool, error) {
	result, err := store.SubExistsByDayAndOwnerID(time, id, ctx, db)

	if err != nil {
		return result, err
	}

	return result, nil
}
