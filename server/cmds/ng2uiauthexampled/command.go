package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/pcdummy/ng2-ui-auth-example/server/components/registry"
	"github.com/pcdummy/ng2-ui-auth-example/server/shared"
)

func commandWrapper(callable func(*cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(c *cobra.Command, args []string) {
		configFile, cfg := shared.IniConfigParse(configFile, debug)

		if !debug {
			debug = cfg.Section("").Key("Debug").MustBool(false)
		}

		if listen == "" {
			listen = cfg.Section("").Key("Listen").MustString(":3000")
		}

		if err := registry.Instance().SetupFromIni(cfg, configFile, debug); err != nil {
			log.Fatalf("ERROR: %v", err)
		}

		callable(c, args)
	}
}
