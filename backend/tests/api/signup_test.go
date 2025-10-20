//go:build func_test

package api_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/app"
	"github.com/ukma-cs-ssdm-2025/team-circus/tests/pkg"
)

func TestSignUpHandler(main *testing.T) {
	setup := func() (*app.App, *sql.DB, error) {
		db, err := pkg.NewDB()
		if err != nil {
			return nil, nil, err
		}
		err = pkg.ResetDB(db)
		if err != nil {
			return nil, nil, err
		}
		app := pkg.NewApp()
		return app, db, nil
	}

	main.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		app, db, err := setup()
		require.NoError(t, err)
		defer db.Close()

		go app.Run(ctx)
		assert.True(t, true)
	})
}
