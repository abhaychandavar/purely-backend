package services

import (
	"context"
	PubSub "profiles/internal/providers/pubSub"
)

type InternalService struct{}

func (i *InternalService) HandleProfileImageBlurred(ctx context.Context, mediaID string, blurredImageID string, profileID string) {
	ps := ProfileService{}
	ps.UpsertProfileBlurredImage(ctx, mediaID, blurredImageID, profileID)
}

func (i *InternalService) HandlePubSubMessage(ctx context.Context, data PubSub.PubSubMessageType) bool {
	switch data.Type {
	case "imageBlurred":
		{
			mediaID := data.Data["mediaID"].(string)
			blurredImageID := data.Data["blurredImageID"].(string)
			profileID := data.Data["profileID"].(string)
			i.HandleProfileImageBlurred(ctx, mediaID, blurredImageID, profileID)
		}
	}
	return true
}
