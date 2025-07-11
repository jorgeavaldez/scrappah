package main

import (
	"fmt"
	"os"

	"scrappah/pkg/db"

	_ "github.com/tursodatabase/go-libsql"
)

func main() {
	dbName := "file:./local.db"

	dbInstance := db.GetDB(dbName)
	defer dbInstance.Close()

	vpnConfigs := db.GetVPNConfigs(dbInstance)

	for _, vpnConfig := range vpnConfigs {
		fmt.Fprintf(os.Stdout, "vpn config: %s\n", vpnConfig.Name)
	}
}
