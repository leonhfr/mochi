package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/internal/image"
)

type deleteCard struct {
	cardID string
}

// DeleteCard returns a new archive card request.
func DeleteCard(cardID string) Request {
	return &deleteCard{cardID: cardID}
}

// Execute implements the Request interface.
func (r *deleteCard) Execute(ctx context.Context, client Client, _ image.Reader, _ Lockfile) error {
	return client.DeleteCard(ctx, r.cardID)
}

// String implements the fmt.Stringer interface.
func (r *deleteCard) String() string {
	return fmt.Sprintf("delete request for card ID %s", r.cardID)
}
