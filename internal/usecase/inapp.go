package usecase

import (
	"context"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
)

func (u *Usecase) HandleInApp(ctx context.Context, parentMsg *model.PublishedKafkaMsg, inappMsg *model.InApp) error {
	// ctx, span := u.tracer.Start(ctx, "Usecase.InAppUsecase")
	// defer span.End()

	// TODO: implement handle inapp

	return nil
}
