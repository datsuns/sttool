package backend

import (
	"fmt"
	"log/slog"

	"github.com/andreykaipov/goobs"
)

// OBS may use IPv6 address. I dont know how to use IPv6 for goobs library
func buildObsUrl(cfg *Config) string {
	//return fmt.Sprintf("%v:%v", cfg.ObsIp(), cfg.ObsPort())
	return fmt.Sprintf("localhost:%v", cfg.ObsPort())
}

func connectToObs(cfg *Config) (*goobs.Client, error) {
	url := buildObsUrl(cfg)
	client, err := goobs.New(url, goobs.WithPassword(cfg.ObsPass()))
	if err != nil {
		logger.Error("OBS Connect ERROR", slog.Any("err", err.Error()))
		return nil, err
	}
	return client, nil
}

func GetObsVersion(cfg *Config) (string, error) {
	client, err := connectToObs(cfg)
	if err != nil {
		return "", nil
	}
	defer client.Disconnect()

	version, err := client.General.GetVersion()
	if err != nil {
		logger.Error("OBS GetVersion ERROR", slog.Any("err", err.Error()))
		return "", err
	}
	return version.ObsVersion, nil
}

func StopObsStream(cfg *Config) {
	client, err := connectToObs(cfg)
	if err != nil {
		return
	}
	defer client.Disconnect()

	res, err := client.Stream.StopStream()
	if err != nil {
		logger.Error("OBS StopStream ERROR", slog.Any("err", err.Error()))
		return
	}
	logger.Info("OBS stop stream", slog.Any("reponce", res))
}
