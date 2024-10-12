package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/mochi"
)

type archive struct {
	cardID string
}

// NewArchive returns a new archive card request.
func NewArchive(cardID string) Request {
	return &archive{
		cardID: cardID,
	}
}

// Sync implements the SyncRequest interface.
func (r *archive) Sync(ctx context.Context, c Client, _ Reader, _ Lockfile) error {
	req := mochi.UpdateCardRequest{Archived: true}
	_, err := c.UpdateCard(ctx, r.cardID, req)
	return err
}

// String implements the fmt.Stringer interface.
func (r *archive) String() string {
	return fmt.Sprintf("archive request for card ID %s", r.cardID)
}
