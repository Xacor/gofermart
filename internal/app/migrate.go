package app

import (
	"errors"
	"time"

	gomigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func migrate(dbURI string, l *zap.Logger) {
	if len(dbURI) == 0 {
		l.Fatal("migrate failed", zap.Error(errors.New("environment variable or flag not declared: DATABASE_URI")))
	}
	var (
		attempts = defaultAttempts
		err      error
		m        *gomigrate.Migrate
	)

	for attempts > 0 {
		m, err = gomigrate.New("file://migrations/migration.sql", dbURI)
		if err == nil {
			break
		}

		l.Error("migrate connection failed", zap.Error(err), zap.Int("postgres is trying to connect, attempts left", attempts))
		time.Sleep(defaultTimeout)
		attempts--
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, gomigrate.ErrNoChange) {
		l.Fatal("migrate up failed", zap.Error(err))
	}

	if errors.Is(err, gomigrate.ErrNoChange) {
		l.Error("migrate", zap.Error(err))
		return
	}

	l.Info("migrate: up success")
}
