/*
Copyright © 2021 JcdeA <jcde@jcde.xyz>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

const API_URL = "https://console.aarch64.com/api"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fhctl",
	Short: "Administer the Fosshost & aarch64.com infrastructure from the terminal",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			os.Exit(0)
		case syscall.SIGINT:
			os.Exit(0)
		}
	}()

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fhctl.yaml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fhctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".fhctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()

}

func createOrWriteConfig(mode fs.FileMode) (err error) {
	if viper.ConfigFileUsed() == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configPath := filepath.Join(home, "/fhctl.json")

		_, err = os.Stat(configPath)
		if !os.IsExist(err) {
			if _, err := os.Create(configPath); err != nil {
				cobra.CheckErr(err)
			}
		}
		if err := viper.SafeWriteConfig(); err != nil {
			cobra.CheckErr(err)
		}

		err = os.Chmod(viper.ConfigFileUsed(), mode)
		if err != nil {
			return err
		}
	} else {
		if err = viper.WriteConfig(); err != nil {
			cobra.CheckErr(err)
		}
	}
	return nil
}
