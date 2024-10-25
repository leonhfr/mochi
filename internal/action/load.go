package action

import (
	"context"
	"time"

	"github.com/sourcegraph/conc/pool"

	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/throttle"
	"github.com/leonhfr/mochi/mochi"
)

// Logger is the interface to log output.
type Logger interface {
	Debugf(format string, args ...any)
	Errorf(format string, args ...any)
	Infof(format string, args ...any)
}

func loadConfig(r config.Reader, logger Logger, parsers []string, workspace string) (*config.Config, error) {
	config, err := config.Parse(r, workspace, parsers)
	if err != nil {
		return nil, err
	}

	logger.Infof("loaded config")
	logger.Debugf("config: %v", config)

	return config, err
}

func loadClient(logger Logger, rateLimit int, token string) *mochi.Client {
	rate, burst := getRate(rateLimit)
	client := mochi.New(
		token,
		mochi.WithTransport(throttle.New(rate, burst)),
	)
	logger.Infof("loaded client")
	return client
}

func loadLockfile(ctx context.Context, logger Logger, client *mochi.Client, rw lock.ReaderWriter, workspace string) (*lock.Lock, error) {
	lf, err := lock.Parse(rw, workspace)
	if err != nil {
		return nil, err
	}

	err = deck.CleanDecks(ctx, client, lf)
	if err != nil {
		return nil, err
	}

	p := pool.New().WithErrors().WithContext(ctx)
	for _, id := range getDeckIDs(lf.Decks()) {
		id := id
		p.Go(func(ctx context.Context) error {
			return deck.CleanCards(ctx, client, lf, id)
		})
	}
	err = p.Wait()
	if err != nil {
		return nil, err
	}

	logger.Infof("loaded lockfile")
	logger.Debugf("lockfile: %v", lf.String())

	return lf, nil
}

func getDeckIDs(decks map[string]lock.Deck) []string {
	ids := make([]string, 0, len(decks))
	for id := range decks {
		ids = append(ids, id)
	}
	return ids
}

func getRate(rateLimit int) (time.Duration, int) {
	return time.Second / time.Duration(rateLimit), rateLimit
}
