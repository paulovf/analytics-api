package worker

import (
	"context"
	"log"

	"github.com/riverqueue/river"
)

type PaymentNotificationArgs struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

func (PaymentNotificationArgs) Kind() string {
	return "payment_notification"
}

type PaymentNotificationWorker struct {
	river.WorkerDefaults[PaymentNotificationArgs]
}

func (w *PaymentNotificationWorker) Work(ctx context.Context, job *river.Job[PaymentNotificationArgs]) error {
	log.Printf("[RIVER WORKER] Proccessing async notification for payment ID: %s (Status: %s)", 
		job.Args.PaymentID, job.Args.Status)

	return nil
}
