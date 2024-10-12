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

// LoadConfig loads the config.
func LoadConfig(r config.Reader, logger Logger, parsers []string, workspace string) (*config.Config, error) {
	config, err := config.Parse(r, workspace, parsers)
	if err != nil {
		return nil, err
	}

	logger.Infof("loaded config")
	logger.Debugf("config: %v", config)

	return config, err
}

// LoadClient loads the client.
func LoadClient(logger Logger, rateLimit int, token string) *mochi.Client {
	rate, burst := getRate(rateLimit)
	client := mochi.New(
		token,
		mochi.WithTransport(throttle.New(rate, burst)),
	)
	logger.Infof("loaded client")
	return client
}

// LoadLockfileClient is the client interface to load the lockfile.
type LoadLockfileClient interface {
	ListDecks(ctx context.Context) ([]mochi.Deck, error)
}

// LoadLockfile loads the lockfile.
func LoadLockfile(ctx context.Context, logger Logger, client LoadLockfileClient, rw lock.ReaderWriter, workspace string) (*lock.Lock, error) {
	lf, err := lock.Parse(rw, workspace)
	if err != nil {
		return nil, err
	}

	decks, err := client.ListDecks(ctx)
	if err != nil {
		return nil, err
	}

	lf.CleanDecks(decks)

	logger.Infof("loaded lockfile")
	logger.Debugf("lockfile: %v", lf.String())

	return lf, nil
}

func getRate(rateLimit int) (time.Duration, int) {
	return time.Second / time.Duration(rateLimit), rateLimit
}
