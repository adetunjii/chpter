package db

import (
	"context"
	"net/url"
	"time"

	"github.com/chpter/shared/lazyerror"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setDefaultValues(values url.Values) {
	// set timeouts for unused pool connections.
	if !values.Has("pool_max_conn_idle_time") {
		values.Set("pool_max_conn_idle_time", "1m")
	}
	values.Set("timezone", "UTC")
}

type Pool struct {
	URI url.URL
}

func openDB(uri string) (*pgxpool.Pool, error) {
	baseURL, err := url.Parse(uri)
	if err != nil {
		return nil, lazyerror.Error(err)
	}

	values := baseURL.Query()
	setDefaultValues(values)
	baseURL.RawQuery = values.Encode()

	poolConfig, err := pgxpool.ParseConfig(baseURL.String())
	if err != nil {
		return nil, lazyerror.Error(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, lazyerror.Error(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// check if connection is active
	if err := pool.Ping(ctx); err != nil {
		return nil, lazyerror.Error(err)
	}

	return pool, nil
}
