package cmd

import (
	"os"

	"github.com/seftomsk/abf/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var cfg config.Config

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "abf",
	Short: "short description",
	Long:  `long description`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and
// sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() config.Config {
	cobra.CheckErr(rootCmd.Execute())
	return cfg
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&cfgFile,
		"config",
		"c",
		"/etc/abf/config.json",
		"Path to config file (default is /etc/abf/config.json)")

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("json")
	viper.SetConfigFile(cfgFile)

	err := viper.ReadInConfig()
	if err != nil {
		_, _ = os.Stderr.Write(
			[]byte("initConfig - viper.ReadInConfig: " + err.Error() + "\n"))
		os.Exit(1)
	}

	cfg = config.New()
	err = viper.Unmarshal(&cfg)
	if err != nil {
		_, _ = os.Stderr.Write(
			[]byte("initConfig - viper.Unmarshal: " + err.Error() + "\n"))
	}
	err = cfg.Validate()
	if err != nil {
		_, _ = os.Stderr.Write(
			[]byte("initConfig - config.Validate: " + err.Error() + "\n"))
		os.Exit(1)
	}
}
