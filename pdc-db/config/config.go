package config

import (
	"gopkg.in/gcfg.v1"
)

type Config struct {
	DbDriver     string
	DbConnection string
}

type configFile struct {
	Server Config
}

//const defaultConfig = `
//   [server]
//   dbConnection = host=postgresContainer user=pdc dbname=pdcDB sslmode=disable password=pdctest
//   dbDriver = postgres
//`

const defaultConfig = `
 [server]
 dbConnection = host=127.0.0.1 user=pdc dbname=pdcDB sslmode=disable password=pdctest
 dbDriver = postgres
`

func LoadConfiguration(cfgFile string) Config {
	var err error
	var cfg configFile

	if cfgFile != "" {
		err = gcfg.ReadFileInto(&cfg, cfgFile)
	} else {
		err = gcfg.ReadStringInto(&cfg, defaultConfig)
	}

	if err != nil {
		panic(err)
	}

	return cfg.Server
}
