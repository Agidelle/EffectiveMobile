package cmd

import (
	"context"
	"fmt"
	"github.com/Agidelle/EffectiveMobile/internal/api"
	"github.com/Agidelle/EffectiveMobile/internal/config"
	"github.com/Agidelle/EffectiveMobile/internal/domain"
	"github.com/Agidelle/EffectiveMobile/internal/service"
	"github.com/Agidelle/EffectiveMobile/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",

	Run: func(cmd *cobra.Command, args []string) {
		initLogger()

		cfg, err := config.LoadCfg()
		if err != nil {
			slog.Error("Error load config file .env", "error", err.Error())
			os.Exit(1)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		store := storage.NewPool(ctx, cfg)
		var repo domain.Repository = store
		var svc api.SubService = service.NewService(repo)
		handler := api.NewHandler(svc)

		r := chi.NewRouter()
		handler.InitRoutes(r)

		srv := &http.Server{
			Addr:    ":" + fmt.Sprintf(cfg.AppPort),
			Handler: r,
		}

		go func() {
			slog.Info("Launch server", "port", cfg.AppPort)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("ListenAndServe error", "error", err.Error())
				os.Exit(1)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		slog.Info("Shutdown Server ...")

		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctxShutdown); err != nil {
			slog.Error("Server Shutdown error", "error", err.Error())
			os.Exit(1)
		}
		svc.CloseDB()
		slog.Info("Server exiting")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
