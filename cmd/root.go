package cmd

import (
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "EffectiveMobile",
	Short: "A brief description of your application",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initLogger() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err.Error())
		os.Exit(1)
	}
	defer file.Close()
	// логируем в файл и в stdout
	mw := io.MultiWriter(os.Stdout, file)
	slog.SetDefault(slog.New(slog.NewJSONHandler(mw, &slog.HandlerOptions{Level: slog.LevelInfo})))
}
