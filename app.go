package google_license_manager

import (
	"fmt"
	"github.com/kesselborn/go-getopt"
	"github.com/linuzilla/go-boltdb"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/go-logger"
	"github.com/linuzilla/google-license-manager/commands"
	"github.com/linuzilla/google-license-manager/config"
	"github.com/linuzilla/google-license-manager/constants"
	"github.com/linuzilla/google-license-manager/models"
	"github.com/linuzilla/google-license-manager/services/admin_sdk"
	"github.com/linuzilla/google-license-manager/services/google_credential"
	"github.com/linuzilla/google-license-manager/services/license_manager"
	"github.com/linuzilla/google-license-manager/utils"
	"github.com/linuzilla/google-license-manager/utils/encryption_helper"
	"github.com/linuzilla/summer"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"syscall"
)

func Main() {
	optionDefinition := getopt.Options{
		Description: constants.VERSION,
		Definitions: getopt.Definitions{
			{"debug|d|DEBUG", "debug mode", getopt.Optional | getopt.Flag, false},
			{"config|c|" + constants.ConfigFileEnv, "config file", getopt.IsConfigFile | getopt.ExampleIsDefault, "application.yml"},
			{"pass through", "pass through arguments", getopt.IsPassThrough | getopt.Optional, ""},
		},
	}

	options, _, passThrough, e := optionDefinition.ParseCommandLine()

	help, wantsHelp := options["help"]
	exitCode := 0

	if e != nil || wantsHelp {
		switch {
		case wantsHelp && help.String == "usage":
			fmt.Print(optionDefinition.Usage())
		case wantsHelp && help.String == "help":
			fmt.Print(optionDefinition.Help())
		default:
			fmt.Println("**** Error: ", e.Error(), "\n", optionDefinition.Help())
			exitCode = e.ErrorCode
		}
	} else {
		startProgram(options["config"].String, passThrough)
	}
	os.Exit(exitCode)
}

func initializeDatabase(applicationContext summer.ApplicationContextManager, configuration *config.Config) dbBolt.DatabaseBackend {
	var databaseBackend dbBolt.DatabaseBackend = &dbBolt.BoltBackend

	databaseFile := constants.DefaultBoltDbFileName

	if configuration != nil && configuration.DatabaseFile != `` {
		databaseFile = configuration.DatabaseFile
	}

	databaseBackend.Initialize(databaseFile, true, &models.LicensedUser{}, &models.EncryptedData{})
	applicationContext.Add(databaseBackend)
	return databaseBackend
}

func registerCommands(applicationContext summer.ApplicationContextManager) {
	var googleCredential google_credential.GoogleCredential

	if result, err := applicationContext.Get(&googleCredential); err != nil {
		log.Fatal(err)
	} else {
		googleCredential = result.(google_credential.GoogleCredential)
	}

	commandList := []cmdline_service.CommandInterface{
		&commands.CustomerCommand{},
		&commands.UserCommand{},
		&commands.ListCommand{},
		&commands.AddUserCommand{},
		&commands.RevokeUserCommand{},
		&commands.DumpCommand{},
		&commands.SyncCommand{},
		&commands.DescribeCommand{},
	}

	for _, cmd := range commandList {
		if condition, ok := cmd.(commands.ConditionalRegister); ok {
			if !condition.CanRegister(googleCredential) {
				continue
			}
		}
		applicationContext.Add(cmd)
	}
}

func startProgram(configYaml string, passThrough []string) {
	fmt.Println(constants.VERSION)
	fmt.Println()
	fmt.Println("Go version: " + runtime.Version())
	fmt.Println("Running on: " + runtime.GOOS)
	fmt.Println("Written by: " + constants.WRITER)
	fmt.Println()

	logger.SetLogger(logger.New())

	applicationContext := summer.New()
	applicationContext.Debug(false)
	applicationContext.Add(applicationContext)

	fmt.Printf("Loading config: \"%s\"\n", configYaml)
	conf, err := config.FromFile(configYaml)
	embedded := false

	if err != nil {
		databaseBackend := initializeDatabase(applicationContext, nil)
		googleCfg, credential, err := encryption_helper.ReadEncryptedData(databaseBackend)

		if err != nil {
			log.Fatal(err)
		}
		applicationContext.Add(googleCfg)
		applicationContext.Add(google_credential.FromData(googleCfg, credential))
		embedded = true
	} else {
		initializeDatabase(applicationContext, conf)
		applicationContext.Add(&conf.GoogleCfg)
		applicationContext.Add(google_credential.New(&conf.GoogleCfg))
		logger.SetLevel(conf.LogLevel)
	}

	applicationContext.Add(license_manager.New())
	applicationContext.Add(admin_sdk.New())

	if embedded {
		applicationContext.Add(&commands.ChangePasswordCommand{})
	} else {
		applicationContext.Add(&commands.StoreAndEncodeCommand{})
	}

	registerCommands(applicationContext)

	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			fmt.Printf(">> panic: %v\n", r)
		}
	}()

	done := applicationContext.Autowiring(func(err error) {
		if err != nil {
			utils.Println(os.Args[0])
			utils.Println(syscall.Getwd())
			utils.Println("Failed to auto wiring.")
			log.Fatalf("Error: %v\n", err)
		}
	})

	if err := <-done; err == nil {
		commandLineService := cmdline_service.New(applicationContext, "Google Workspace")

		if len(passThrough) == 0 {
			commandLineService.Execute()
		} else {
			for _, cmd := range passThrough {
				fmt.Println()
				commandLineService.RunCommand(cmd)
			}
		}
	} else {
		log.Fatal(err)
	}
}
