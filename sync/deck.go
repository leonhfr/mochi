package sync

import "context"

type DeckResult struct {
	Created int
	Updated int
}

func SynchronizeDecks(_ context.Context, _ []string, _ *Lock, _ Config, _ Client) (DeckResult, error) {
	return DeckResult{}, nil
}
