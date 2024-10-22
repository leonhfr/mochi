package request

import (
	"context"
	"fmt"

	"github.com/leonhfr/mochi/mochi"
)

type archiveCard struct {
	cardID string
}

// ArchiveCard returns a new archive card request.
func ArchiveCard(cardID string) Request {
	return &archiveCard{
		cardID: cardID,
	}
}

// Sync implements the SyncRequest interface.
func (r *archiveCard) Sync(ctx context.Context, client Client, _ Reader, _ Lockfile) error {
	req := mochi.UpdateCardRequest{Archived: true}
	_, err := client.UpdateCard(ctx, r.cardID, req)
	return err
}

// String implements the fmt.Stringer interface.
func (r *archiveCard) String() string {
	return fmt.Sprintf("archive request for card ID %s", r.cardID)
}
