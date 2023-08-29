package app

import (
	"context"
	"errors"
	"time"

	"github.com/Xacor/gophermart/pkg/postgres"
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
	const sql = `
BEGIN;
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE IF NOT EXISTS users
(
    id serial,
    login character varying(256) NOT NULL,
    password character varying(256) NOT NULL,
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS orders CASCADE;

CREATE TABLE IF NOT EXISTS orders
(
    id text,
    user_id serial NOT NULL,
    status character varying NOT NULL,
    accrual bigint NOT NULL,
    uploaded_at timestamp with time zone NOT NULL,
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS withdrawals CASCADE;

CREATE TABLE IF NOT EXISTS withdrawals
(
    id serial,
    order_id text,
    user_id serial NOT NULL,
    sum bigint NOT NULL DEFAULT 0,
    processed_at timestamp with time zone NOT NULL,
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS balances CASCADE;

CREATE TABLE IF NOT EXISTS balances
(
    user_id serial,
    current bigint,
    withdrawn bigint,
    PRIMARY KEY (user_id),
    UNIQUE (user_id)
);


ALTER TABLE IF EXISTS orders
    ADD CONSTRAINT "FK_orders_users" FOREIGN KEY (user_id)
    REFERENCES users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS withdrawals
    ADD CONSTRAINT "FK_withdrawals_users" FOREIGN KEY (user_id)
    REFERENCES users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS balances
    ADD FOREIGN KEY (user_id)
    REFERENCES users (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;

END;`
	var (
		attempts = defaultAttempts
		err      error
		pg       *postgres.Postgres
	)

	for attempts > 0 {
		attempts--
		pg, err = postgres.New(dbURI)
		if err != nil {
			l.Error("migrate connection failed", zap.Error(err), zap.Int("postgres is trying to connect, attempts left", attempts))
			continue
		}

		_, err = pg.Pool.Exec(context.Background(), sql)
		if err == nil {
			break
		}
		l.Error("migration failed", zap.Error(err))

		time.Sleep(defaultTimeout)
	}

	defer pg.Close()

	l.Info("migrate: up success")
}
