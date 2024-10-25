package action

import (
	"context"
	"time"

	"github.com/leonhfr/mochi/internal/config"
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

	decks, err := client.ListDecks(ctx)
	if err != nil {
		return nil, err
	}

	lf.Lock()
	defer lf.Unlock()
	lf.CleanDecks(decks)

	logger.Infof("loaded lockfile")
	logger.Debugf("lockfile: %v", lf.String())

	return lf, nil
}

func getRate(rateLimit int) (time.Duration, int) {
	return time.Second / time.Duration(rateLimit), rateLimit
}
