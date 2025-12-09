package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func respond(msg string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		log.Printf("[WARN] bot failed to message a response")
	}
}

func err_respond(msg string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("An unexpected error occurred: %s", msg),
		},
	})
	if err != nil {
		log.Printf("[WARN] bot failed to message a response for %s", msg)
	}
	log.Printf("[ERR] an unexpected error occurred: %s", msg)
}
