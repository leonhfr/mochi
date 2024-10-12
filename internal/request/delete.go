package request

import (
	"context"
	"fmt"
)

type deleteCardRequest struct {
	cardID string
}

// NewDelete returns a new archive card request.
func NewDelete(cardID string) Request {
	return &deleteCardRequest{cardID: cardID}
}

// Sync implements the SyncRequest interface.
func (r *deleteCardRequest) Sync(ctx context.Context, client Client, _ Reader, _ Lockfile) error {
	return client.DeleteCard(ctx, r.cardID)
}

// String implements the fmt.Stringer interface.
func (r *deleteCardRequest) String() string {
	return fmt.Sprintf("delete request for card ID %s", r.cardID)
}
