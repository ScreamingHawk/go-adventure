package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/ScreamingHawk/go-adventure/config"
	"github.com/ScreamingHawk/go-adventure/server"
	"github.com/go-chi/httplog/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "app is a tool for managing your application",
	Long:  `app is a tool for managing your application. It is a tool that can be used to create, manage, and deploy applications.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initConfig(); err != nil {
			return err
		}
		return run()
	},
}

var (
	configFile string
	cfg = &config.Config{}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./etc/app.conf)")
}

func initConfig() error {
	if err := config.NewFromFile(configFile, cfg); err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	return nil
}

func run() error {
	logger := httplog.NewLogger("app", httplog.Options{
		LogLevel: httplog.LevelByName(cfg.Logger.Level),
		Concise: cfg.Logger.Concise,
	})

	s, err := server.NewServer(&cfg.Server, logger)
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}

	ctx := context.Background()

	// Graceful stop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		s.Stop(ctx)
	}()

	// Run it
	if err := s.Run(ctx); err != nil {
		fmt.Printf("failed to start server: %v\n", err)
	}

	return nil
}
