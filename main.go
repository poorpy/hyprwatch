package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/poorpy/hyprwatch/internal/config"
	"github.com/poorpy/hyprwatch/internal/event"
)

const his = "HYPRLAND_INSTANCE_SIGNATURE"

type ProgramFlags struct {
	IsDebug bool
}

func parseFlags() ProgramFlags {
	var config ProgramFlags

	flag.BoolVar(&config.IsDebug, "debug", false, "set debug mode")
	flag.Parse()

	return config
}

func main() {
	flags := parseFlags()

	if flags.IsDebug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Level(zerolog.DebugLevel)
	} else {
		log.Level(zerolog.InfoLevel)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	config, err := config.NewConfig(filepath.Join(home, ".config/hyprwatch/config.yaml"))
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	socketPath := fmt.Sprintf("/tmp/hypr/%s/.socket2.sock", os.Getenv(his))

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Error().Err(err).Msg("")
			continue
		}

		for _, line := range strings.Split(string(buffer[:n]), "\n") {
			if line == "" {
				continue
			}

			event := event.NewEvent(line)
			log.Debug().Str("event", event.Event).Str("data", event.Data).Msg("")

			commands, ok := config.Lookup(event.Event)
			if !ok {
				log.Debug().Str("event", event.Event).Msg("no commands for event")
				continue
			}

			for _, command := range commands {
				if command.Data != event.Data {
					continue
				}

				// args := strings.Fields(command.Callback)
				cmd := exec.Command("sh", "-c", command.Callback) //nolint:gosec

				var (
					stdout bytes.Buffer
					stderr bytes.Buffer
				)
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr

				if err := cmd.Run(); err != nil {
					log.Error().
						Err(err).
						Str("event", event.Event).
						Str("data", event.Data).
						Str("callback", command.Callback).
						Str("stdout", stdout.String()).
						Str("stderr", stderr.String()).
						Msg("failed to execute callback")
					continue
				}

				log.Info().
					Str("event", event.Event).
					Str("data", event.Data).
					Str("callback", command.Callback).
					Str("stdout", stdout.String()).
					Msg("executing callback")
			}
		}
	}
}
