package handler

import (
	"context"

	"git.bpjsketenagakerjaan.go.id/centralized-notification-system/dispatch-service/internal/model"
)

func (h *Handler) HandleInApp(ctx context.Context, parentMsg *model.PublishedKafkaMsg, inappMsg *model.InApp, traceID string) error {
	// ctx, span := h.tracer.Start(ctx, "Handler.InAppHandler")
	// defer span.End()

	// TODO: implement handle inapp

	return nil
}
