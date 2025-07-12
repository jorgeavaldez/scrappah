package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/pufferffish/wireproxy"

	"scrappah/pkg"
	"scrappah/pkg/db"
)

type App struct {
	Ctx  context.Context
	Repo *db.Repository
}

func NewApp() *App {
	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGINT, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-s
		cancel()
	}()

	dbName := "file:./local.db"
	repo := db.NewRepository(ctx, dbName)

	return &App{
		Ctx:  ctx,
		Repo: repo,
	}
}

func (app *App) Close() error {
	err := app.Repo.Close()
	return err
}

func (app *App) LoadVPNConfigs() []*wireproxy.Configuration {
	vpnConfigs := app.Repo.GetVPNConfigs()

	validCount := 0
	invalidCount := 0

	startingPort := 8001

	var builder strings.Builder

	validConfigs := make([]*wireproxy.Configuration, 0, len(vpnConfigs))
	for _, vpnConfig := range vpnConfigs {
		builder.WriteString(string(vpnConfig.ConfigContent))
		builder.WriteString(
			"\n\n[Socks5]\nBindAddress = 0.0.0.0:" + strconv.Itoa(startingPort) + "\n",
		)
		vpnConfig.ConfigContent = []byte(builder.String())
		startingPort++
		builder.Reset()

		validConfig, err := pkg.ValidateVPNConfig(vpnConfig)
		if err != nil {
			fmt.Printf("INVALID: %s (ID: %d) - %s\n", vpnConfig.Name, vpnConfig.ID, err)
			invalidCount++
		} else {
			fmt.Printf("VALID: %s (ID: %d)\n", vpnConfig.Name, vpnConfig.ID)
			validCount++
			validConfigs = append(validConfigs, validConfig)
		}
	}

	fmt.Printf("\nValidation Summary: %d valid, %d invalid configs\n", validCount, invalidCount)
	return validConfigs
}

func main() {
	app := NewApp()
	defer app.Close()

	configs := app.LoadVPNConfigs()

	for _, config := range configs {
		if config == nil {
			fmt.Println("No valid configs found")
			os.Exit(1)
		}

		tun, err := wireproxy.StartWireguard(config.Device, 1)
		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(1)
		}

		for _, spawner := range config.Routines {
			go spawner.SpawnRoutine(tun)
		}

		tun.StartPingIPs()
	}

	<-app.Ctx.Done()
}
