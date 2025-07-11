package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"scrappah/pkg/db"

	_ "github.com/tursodatabase/go-libsql"
)

func main() {
	ctx := context.Background()
	dbName := "file:./local.db"

	repo := db.NewRepository(ctx, dbName)
	defer repo.Close()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  list              List all VPN configs\n")
		fmt.Fprintf(os.Stderr, "  add <file_path>   Add VPN config from file\n")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		listVPNConfigs(repo)
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s add <file_path>\n", os.Args[0])
			os.Exit(1)
		}
		filePath := os.Args[2]
		addVPNConfig(repo, filePath)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func listVPNConfigs(repo *db.Repository) {
	vpnConfigs := repo.GetVPNConfigs()

	for _, vpnConfig := range vpnConfigs {
		fmt.Fprintf(os.Stdout, "vpn config: %s\n", vpnConfig.Name)
	}
}

func addVPNConfig(repo *db.Repository, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file %s: %s\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file %s: %s\n", filePath, err)
		os.Exit(1)
	}

	fileName := filepath.Base(filePath)
	vpnConfig := db.VPNConfig{
		Name:          fileName,
		IsActive:      false,
		ConfigContent: content,
	}

	id, err := repo.InsertVPNConfig(vpnConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to insert VPN config: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added VPN config '%s' with ID %d\n", fileName, id)
}
