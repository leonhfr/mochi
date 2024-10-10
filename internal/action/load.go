package action

import (
	"context"
	"time"

	"github.com/leonhfr/mochi/internal/card"
	"github.com/leonhfr/mochi/internal/config"
	"github.com/leonhfr/mochi/internal/deck"
	"github.com/leonhfr/mochi/internal/file"
	"github.com/leonhfr/mochi/internal/lock"
	"github.com/leonhfr/mochi/internal/parser"
	"github.com/leonhfr/mochi/internal/throttle"
	"github.com/leonhfr/mochi/internal/worker"
	"github.com/leonhfr/mochi/mochi"
)

// Logger is the interface to log output.
type Logger interface {
	Debugf(format string, args ...any)
	Errorf(format string, args ...any)
	Infof(format string, args ...any)
}

// Client interface.
type Client interface {
	worker.Client
	deck.Client
	card.Client
}

// Filesystem interface.
type Filesystem interface {
	worker.Walker
	card.Reader
}

// Parser interface.
type Parser interface {
	card.Parser
	Extensions() []string
}

// Config interface.
type Config interface {
	deck.Config
}

// Lockfile interface.
type Lockfile interface {
	worker.Lockfile
	deck.Lockfile
	Updated() bool
	Write() error
}

// Load loads client, parser, and lockfile.
func Load(ctx context.Context, logger Logger, token, workspace string) (Client, Filesystem, Parser, Config, Lockfile, error) {
	logger.Infof("workspace: %s", workspace)

	fs := file.NewSystem()
	parser := parser.New()

	config, err := config.Parse(fs, workspace, parser.List())
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	logger.Infof("loaded config")
	logger.Debugf("config: %v", config)

	rate, burst := getRate(config.RateLimit)
	client := mochi.New(
		token,
		mochi.WithTransport(throttle.New(rate, burst)),
	)

	lf, err := getLockfile(ctx, client, fs, workspace)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	logger.Infof("loaded lockfile")
	logger.Debugf("lockfile: %v", lf.String())

	return client, fs, parser, config, lf, nil
}

func getRate(rateLimit int) (time.Duration, int) {
	return time.Second / time.Duration(rateLimit), rateLimit
}

func getLockfile(ctx context.Context, client *mochi.Client, fs *file.System, workspace string) (*lock.Lock, error) {
	lf, err := lock.Parse(fs, workspace)
	if err != nil {
		return nil, err
	}

	var decks []mochi.Deck
	if err := client.ListDecks(ctx, func(dd []mochi.Deck) error {
		decks = append(decks, dd...)
		return nil
	}); err != nil {
		return nil, err
	}

	lf.CleanDecks(decks)
	return lf, nil
}
