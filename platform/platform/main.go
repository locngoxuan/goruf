package main

import (
	"context"
	"goruf/platform/database"
	"goruf/platform/http"
	"goruf/platform/storage"
	"goruf/platform/tcp"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

var (
	version = "devel"
)

func main() {
	setupZeroLog()
	cmd := &cli.Command{
		Name:      "mfe",
		Copyright: "xuanloc0511@gmail.com",
		Version:   version,
		Flags:     getFlags(),
		Action:    run,
	}
	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Error().Interface("args", os.Args).Err(err).Msg("failed to run application")
	}
}

func setupZeroLog() error {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.TimestampFieldName = "timestamp"
	log.Logger = log.Output(os.Stdout).With().Caller().Logger()
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	return nil
}

func getFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Sources: cli.EnvVars("PORT"),
			Usage:   "main entry port from UI",
			Value:   80,
		},
		&cli.IntFlag{
			Name:    "cluster.port",
			Aliases: []string{"cp"},
			Sources: cli.EnvVars("CLUSTER_PORT"),
			Usage:   "port to accept connection from client",
			Value:   8081,
		},
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	httpPort := cmd.Int("port")
	clusterPort := cmd.Int("cluster.port")
	err := database.ConnectDatabase()
	if err != nil {
		return err
	}
	err = storage.ConnectStorage()
	if err != nil {
		return err
	}
	err = http.StartWebService(httpPort)
	if err != nil {
		return err
	}
	return tcp.OpenListener(int(clusterPort), func() tcp.MessageHandler {
		return NewServerMessageHandler()
	})
}
