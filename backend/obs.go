package backend

import (
	"fmt"
	"log/slog"

	"github.com/andreykaipov/goobs"
)

func StopObsStream(cfg *Config) {
	url := fmt.Sprintf("%v:%v", cfg.ObsIp(), cfg.ObsPort())
	client, err := goobs.New(url, goobs.WithPassword(cfg.ObsPass()))
	if err != nil {
		logger.Error("OBS Connect ERROR", slog.Any("err", err.Error()))
		return
	}
	defer client.Disconnect()

	//version, err := client.General.GetVersion()
	//if err != nil {
	//	logger.Info("OBS Connect ERROR", slog.Any("err", err.Error()))
	//}
	//fmt.Printf("OBS Studio version: %s\n", version.ObsVersion)
	res, err := client.Stream.StopStream()
	if err != nil {
		logger.Error("OBS Request ERROR", slog.Any("err", err.Error()))
	}
	logger.Info("OBS stop stream", slog.Any("reponce", res))
}
