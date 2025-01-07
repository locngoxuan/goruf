package main

import (
	"bufio"
	"context"
	"fmt"
	"goruf/platform/core"
	"goruf/platform/tcp"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

var (
	version = "devel"
)

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

func main() {
	setupZeroLog()
	cmd := &cli.Command{
		Name:      "mfe.cli",
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

func getFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"a"},
			Sources: cli.EnvVars("ADDRESS"),
			Usage:   "address of control plane",
			Value:   "localhost:8081",
		},
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Sources: cli.EnvVars("FILE"),
			Usage:   "path to deployment file",
			Value:   "",
		},
		&cli.UintFlag{
			Name:    "max-payload-size",
			Aliases: []string{"mps"},
			Sources: cli.EnvVars("MAX_PAYLOAD_SIZE"),
			Usage:   "max size of payload of package",
			Value:   1024 * 1024,
		},
	}
}

func run(ctx context.Context, cmd *cli.Command) error {
	maxPayloadSize := cmd.Uint("max-payload-size")
	f := cmd.String("file")
	if strings.TrimSpace(f) == "" {
		return fmt.Errorf("file must be specified")
	}
	depl, err := readDeploymentFile(f)
	if err != nil {
		return err
	}
	addr := cmd.String("address")
	if strings.TrimSpace(addr) == "" {
		return fmt.Errorf("address of platform must be specified")
	}
	return tcp.ConnectAndTransferData(addr, func(conn net.Conn) error {
		w := bufio.NewWriter(conn)
		b, _ := yaml.Marshal(depl)
		connectCmd := core.CmdConnect{
			Cmd:     core.CmdConnectReq,
			Payload: b,
		}
		msgs, err := tcp.Pack(connectCmd.Pack(), uint32(maxPayloadSize))
		if err != nil {
			return err
		}
		for _, msg := range msgs {
			_, err := w.Write(msg)
			if err != nil {
				return err
			}
		}
		err = w.Flush()
		if err != nil {
			return err
		}
		return nil
	})
}

func readDeploymentFile(f string) (core.DeploymentRequest, error) {
	var d core.DeploymentRequest
	b, err := os.ReadFile(f)
	if err != nil {
		return d, err
	}
	err = yaml.Unmarshal(b, &d)
	if err != nil {
		return d, err
	}
	return d, nil
}
