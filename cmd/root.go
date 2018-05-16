package cmd

import (
	"path"

	"github.com/go-playground/errors"
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	logQuiet   bool
	logVerbose bool
)

var rootCmd = &cobra.Command{
	Use:   "keyvault-certdeploy",
	Short: "X.509-Certificate deployment helper for Azure Key Vault",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if logQuiet && logVerbose {
			return errors.New("quiet and verbose are mutually exclusive")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"",
		"Config file (default locations are $HOME/.config/keyvault-certdeploy.yml, /etc/keyvault-certdeploy.yml, $PWD/keyvault-certdeploy.yml)",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&logVerbose,
		"verbose",
		"v",
		false,
		"Be more verbose",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&logQuiet,
		"quiet",
		"q",
		false,
		"Be quiet",
	)
}

// initLogger sets loglevels based on flags
func initLogger() {
	con := console.New(true)
	levels := []log.Level{}

	if logVerbose {
		levels = log.AllLevels
	} else if logQuiet {
		levels = []log.Level{log.FatalLevel, log.AlertLevel, log.PanicLevel}
	} else {
		levels = []log.Level{log.FatalLevel, log.AlertLevel, log.PanicLevel, log.ErrorLevel, log.WarnLevel}
	}

	log.AddHandler(con, levels...)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("keyvault-certdeploy")

	// set defaults
	//viper.SetDefault("xxx", "xxx")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()

		if err != nil {
			viper.AddConfigPath(path.Join(home, ".config"))
		}

		viper.AddConfigPath("/etc")
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("kvcd")
	viper.AutomaticEnv()

	// if a config file is found, read it in.
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Could not open config file.")
	}
}
