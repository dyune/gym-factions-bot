package handlers

import (
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
