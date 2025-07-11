package db

import (
	"fmt"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

/*
```sql
create table vpn_config (
	id integer primary key autoincrement,
	name text not null unique,
	is_active integer default 0 check (is_active in (0, 1)),
	config_content text not null

);
```
*/

type VPNConfig struct {
	ID            int
	Name          string
	IsActive      bool
	ConfigContent []byte
}

func (r *Repository) GetVPNConfigs() []VPNConfig {
	rows, err := r.db.QueryContext(r.ctx, "SELECT id, name, is_active, config_content FROM vpn_config")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query vpn configs %s", err)
		os.Exit(1)
	}

	defer rows.Close()

	var vpnConfigs []VPNConfig

	for rows.Next() {
		var vpnConfig VPNConfig

		if err := rows.Scan(
			&vpnConfig.ID,
			&vpnConfig.Name,
			&vpnConfig.IsActive,
			&vpnConfig.ConfigContent,
		); err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan vpn config %s", err)
			return nil
		}

		vpnConfigs = append(vpnConfigs, vpnConfig)
		fmt.Fprintf(os.Stdout, "vpn config: %s\n\t%s", vpnConfig.Name, vpnConfig.ConfigContent)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to scan vpn configs %s", err)
	}

	return vpnConfigs
}

func (r *Repository) InsertVPNConfig(vpnConfig VPNConfig) (int, error) {
	result := r.db.QueryRowContext(r.ctx,
		"insert into vpn_config (name, is_active, config_content) values (?, ?, ?) returning id",
		vpnConfig.Name,
		vpnConfig.IsActive,
		vpnConfig.ConfigContent,
	)

	var id int

	if err := result.Scan(&id); err != nil {
		fmt.Fprintf(os.Stderr, "failed to insert vpn config %s", err)
		return -1, err
	}

	fmt.Printf("inserted vpn config with id %d\n", id)
	return id, nil
}
