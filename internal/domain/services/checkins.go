package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/davidwang/factions/internal/config"
	store2 "github.com/davidwang/factions/internal/infra/postgres"
	"github.com/uptrace/bun"
)

func HandleSubmission(s *discordgo.Session, i *discordgo.InteractionCreate) {

	submissionTime := time.Now()

	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		errRespond("channel retrieval created an error", s, i)
		return
	}
	if channel == nil {
		errRespond("channel is nil", s, i)
		return
	}

	id, err := strconv.Atoi(i.Member.User.ID)
	if err != nil {
		errRespond(fmt.Sprintf("id %s could not be parsed", i.Member.User.ID), s, i)
		return
	}

	exists, err := store2.AccountExistsByID(id, config.GlobalCtx, config.DB)
	if !exists {
		respond("You aren't registered. Register to start earning points! Use the /register command.", s, i)
		return
	}
	if err != nil {
		errRespond("could not call database", s, i)
	}

	var faction *store2.Faction
	faction, err = store2.GetAccountFactionByID(id, config.GlobalCtx, config.DB)
	if err != nil {
		errRespond(fmt.Sprintf("%v", err), s, i)
		return
	}
	if faction.Name != store2.FactionNames[channel.Name] {
		respond("Please upload activities in your own gym factions channel", s, i)
		return
	}

	var reachedLimit bool
	reachedLimit, err = reachedSubmissionLimit(submissionTime, id, config.GlobalCtx, config.DB)
	if !reachedLimit {
		err = store2.InsertSub(submissionTime, 1, id, 1, config.GlobalCtx, config.DB)
		if err != nil {
			errRespond(fmt.Sprintf("%v", err), s, i)
			return
		} else {
			respond("Successfully registered your submission! You've earned some points!", s, i)
		}
	}
	if err != nil {
		errRespond(fmt.Sprintf("%v", err), s, i)
	} else {
		respond("Can't submit more than one activity per day!", s, i)
		log.Printf("%v", submissionTime)
	}
}

func reachedSubmissionLimit(time time.Time, id int, ctx context.Context, db *bun.DB) (bool, error) {
	// Check if there exists a submission for today.
	// This helper stays intentionally small so it is easy to cover with
	// integration tests once the postgres layer is exercised against Postgres.
	result, err := store2.SubExistsByDayAndOwnerID(time, id, ctx, db)
	if err != nil {
		return result, err
	}

	return result, nil
}
