package settings

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"code.google.com/p/go.crypto/bcrypt"
	_ "github.com/golang/glog"
	"github.com/mgutz/ansi"
)

var (
	SETTINGS_DATABASE_DSN      string
	SETTINGS_MESSAGE_QUEUE_DSN string
	SETTINGS_CACHE_DSN         string
	SETTINGS_SECRET_KEY        string
	SETTINGS_GOMAXPROCS        int
	SETTINGS_BCRYPT_COST       int

	// restricted plugin ids - no other plugin in the future can have one of these ids
	RESTRICTED_PLUGIN_IDS []string

	// just logo
	PATROL_LOGO = `
      ___         ___                       ___           ___
     /  /\       /  /\          ___        /  /\         /  /\
    /  /::\     /  /::\        /  /\      /  /::\       /  /::\
   /  /:/\:\   /  /:/\:\      /  /:/     /  /:/\:\     /  /:/\:\    ___     ___
  /  /:/~/:/  /  /:/~/::\    /  /:/     /  /:/~/:/    /  /:/  \:\  /__/\   /  /\
 /__/:/ /:/  /__/:/ /:/\:\  /  /::\    /__/:/ /:/___ /__/:/ \__\:\ \  \:\ /  /:/
 \  \:\/:/   \  \:\/:/__\/ /__/:/\:\   \  \:\/:::::/ \  \:\ /  /:/  \  \:\  /:/
  \  \::/     \  \::/      \__\/  \:\   \  \::/~~~~   \  \:\  /:/    \  \:\/:/
   \  \:\      \  \:\           \  \:\   \  \:\        \  \:\/:/      \  \::/
    \  \:\      \  \:\           \__\/    \  \:\        \  \::/        \__\/
     \__\/       \__\/                     \__\/         \__\/                  ver ` +
		VERSION + `

    by phonkee
    `
	commandcolor = ansi.ColorFunc("yellow+h:black")
	lccolor      = ansi.ColorFunc("green+h:black")

	GREEN_COLOR = ansi.ColorFunc("green+h:black")
	RED_COLOR   = ansi.ColorFunc("red+h:black")
	YELLOW      = ansi.ColorFunc("yellow+h:black")

	CORS_ACCESS_ORIGIN = "*"
)

func init() {
	RESTRICTED_PLUGIN_IDS = []string{}

	flag.StringVar(&SETTINGS_DATABASE_DSN, "db_dsn", "postgres://patrol:patrol@localhost/patrol", "database dsn")
	flag.StringVar(&SETTINGS_MESSAGE_QUEUE_DSN, "queue_dsn", "redis://localhost:6379", "ergoq message queue dsn")
	flag.StringVar(&SETTINGS_CACHE_DSN, "cache_dsn", "redis://localhost:6379/1?prefix=patrol", "gocacher cache dsn")
	flag.StringVar(&SETTINGS_SECRET_KEY, "secret_key", "", "secret key for various hashing")
	flag.IntVar(&SETTINGS_GOMAXPROCS, "gomaxprocs", 0, "gomaxprocs, if set to 0 runtime.NumCPU will be used.")
	flag.IntVar(&SETTINGS_BCRYPT_COST, "bcrypt_cost", bcrypt.DefaultCost, fmt.Sprintf("bcrypt hash cost, valid values are %d <= value <= %d.", bcrypt.MinCost, bcrypt.MaxCost))

	if os.Getenv("TESTING") != "TRUE" {

		flag.Usage = func() {
			fmt.Println(PATROL_LOGO)
			fmt.Println(commandcolor("\n  run ./patrol <command> <arg1> <arg2>\n"))
			fmt.Println("\n  For complete list of commands run ", lccolor("./patrol list_commands\n"))
			flag.PrintDefaults()
		}

		// flag.Parse()
	}

	// set GOMAXPROCS
	if SETTINGS_GOMAXPROCS == 0 {
		SETTINGS_GOMAXPROCS = runtime.NumCPU()
	}

	if SETTINGS_BCRYPT_COST > bcrypt.MaxCost || SETTINGS_BCRYPT_COST < bcrypt.MinCost {
		SETTINGS_BCRYPT_COST = bcrypt.DefaultCost
	}
}
