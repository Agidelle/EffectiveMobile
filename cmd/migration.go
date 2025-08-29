package cmd

import (
	"fmt"
	"github.com/Agidelle/EffectiveMobile/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"strconv"
)

var migrationCmd = &cobra.Command{
	Use:   "migration [up|down] [steps]",
	Short: "Manage database migrations",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		initLogger()
		cfg, err := config.LoadCfg()
		if err != nil {
			slog.Error("Error load config file .env", "error", err.Error())
			os.Exit(1)
		}
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
		m, err := migrate.New(
			"file://migrations",
			dsn,
		)
		if err != nil {
			slog.Error("Migration initialization error", "error", err.Error())
			os.Exit(1)
		}

		action := args[0]
		steps := 0
		if len(args) == 2 {
			steps, err = strconv.Atoi(args[1])
			if err != nil {
				slog.Error("Invalid steps argument", "error", err.Error())
				os.Exit(1)
			}
		}

		switch action {
		case "up":
			if steps > 0 {
				err = m.Steps(steps)
			} else {
				err = m.Up()
			}
		case "down":
			if steps > 0 {
				err = m.Steps(-steps)
			} else {
				err = m.Down()
			}
		default:
			slog.Error("Unknown action", "action", action)
			os.Exit(1)
		}

		if err != nil && err != migrate.ErrNoChange {
			slog.Error("Migration error", "error", err.Error())
			os.Exit(1)
		}
		slog.Info("Migration completed", "action", action)
	},
}

func init() {
	rootCmd.AddCommand(migrationCmd)
}
