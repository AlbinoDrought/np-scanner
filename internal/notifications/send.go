package notifications

import (
	"context"

	"go.albinodrought.com/neptunes-pride/internal/multierror"
)

func SendGuarded(ctx context.Context, guard Guard, notifiables []Notifiable, sinks []Sink) error {
	errors := []error{}
	for _, notifiable := range notifiables {
		alreadySent, err := guard.CheckSent(notifiable)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		if alreadySent {
			continue
		}

		sendOK := true
		for _, sink := range sinks {
			err := sink.Send(ctx, notifiable)
			if err != nil {
				errors = append(errors, err)
				sendOK = false
			}
		}

		if sendOK {
			err = guard.RecordSent(notifiable)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}
	return multierror.Optional(errors)
}
