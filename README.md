# go-config

This little Go package aims to help users to set up configuration in their applications, without the need to duplicate configuration code across projects.

Under the hood, this package use the excellent [Viper](https://github.com/spf13/viper) package.

## Install

Run the following command to add the package to your project :

```bash
go get github.com/thomasgouveia/go-config
```

## Example

Below an example of how you can use the package :

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thomasgouveia/go-config"
)

type Config struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

func main() {
	loader, err := config.NewLoader(&config.Options[Config]{
		// You can use config.JSON also if you prefer
		Format: config.YAML,

		// Configuration file
		FileName:      "my-application",
		FileLocations: []string{"/etc/my-application", "."}, // Will search for a "my-application.yaml" file into the directories

		// Enable automatic environment variables lookup
		EnvEnabled: true,
		EnvPrefix:  "app",

		Default: &Config{
			Addr: "0.0.0.0",
			Port: 8080,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Override settings from environment variables
	os.Setenv("APP_PORT", "3000")

	cfg, err := loader.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Addr = %s\n", cfg.Addr)
	fmt.Printf("Port = %d\n", cfg.Port)
}
```
