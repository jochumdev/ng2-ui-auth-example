package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/apps/auth"
	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/apps/settings"
	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/apps/statichttp"
	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/components/auth"
	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/components/auth/db/json"
	//	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/components/auth/db/mongodb"
	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/components/jsonstore"
	//	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/components/mongodb"
	"github.com/pcdummy/ng2-ui-auth-example/server/parts/components/registry"
	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/components/settings"
	_ "github.com/pcdummy/ng2-ui-auth-example/server/parts/components/settings/db/json"
	"github.com/spf13/cobra"
)

var (
	configFile string
	debug      bool
	listen     string
)

var rootCommand = &cobra.Command{
	Use:   "lxdweb",
	Short: "LXDWeb the Web Panel for LXD",
	Long:  ``,
}

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Run the web server",
	Long:  ``,
	Run:   commandWrapper(serve),
}

func serve(cmd *cobra.Command, args []string) {

	log.Printf("Starting listening on %v", listen)

	e := echo.New()
	e.HideBanner = true

	// Recover from panics
	e.Use(middleware.Recover())

	// Logger
	e.Use(middleware.Logger())

	if err := registry.Instance().SetupEcho(e); err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	// Start server
	e.Logger.Fatal(e.Start(listen))

	// TODO: Make this work with os/signal
	if err := registry.Instance().Shutdown(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func main() {
	rootCommand.AddCommand(serveCommand)
	rootCommand.AddCommand(authCommand)
	rootCommand.PersistentFlags().StringVar(&configFile, "config", "/etc/lxdweb/lxdweb.ini", "config file")
	rootCommand.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug")
	rootCommand.PersistentFlags().StringVar(&listen, "listen", "", "host:port to listen on")

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
