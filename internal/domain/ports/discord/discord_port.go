package discord

import "context"

type PublisherDiscord interface {
	Respond(ctx context.Context, interactionID string, message string) error
	RespondEphemeral(ctx context.Context, interactionID string, message string) error
	RespondError(ctx context.Context, interactionID string, message string) error
	EditResponse(ctx context.Context, interactionID string, message string) error
	DeleteResponse(ctx context.Context, interactionID string) error
	SendChannelMessage(ctx context.Context, channelID string, message string) error
}
