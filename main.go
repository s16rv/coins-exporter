package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	ListenAddress string
	CoingeckoApi  string
	Currency      string
	LogLevel      string
	JsonOutput    bool
	ConstLabels   map[string]string
)

var log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

var rootCmd = &cobra.Command{
	Use:  "coins-exporter",
	Long: "Scrape the data about coins in the Coingecko.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if !f.Changed && viper.IsSet(f.Name) {
				val := viper.Get(f.Name)
				if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
					log.Fatal().Err(err).Msg("Could not set flag")
				}
			}
		})
		return nil
	},
	Run: Execute,
}

func Execute(cmd *cobra.Command, args []string) {
	logLevel, err := zerolog.ParseLevel(LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not parse log level")
	}

	if JsonOutput {
		log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	zerolog.SetGlobalLevel(logLevel)

	log.Info().
		Str("--listen-address", ListenAddress).
		Str("--currency", Currency).
		Str("--coingecko-api", CoingeckoApi).
		Str("--log-level", LogLevel).
		Bool("--json", JsonOutput).
		Msg("Started with following parameters")

	http.HandleFunc("/metrics/coins", func(w http.ResponseWriter, r *http.Request) {
		CoinsHandler(w, r, CoingeckoApi)
	})

	log.Info().Str("address", ListenAddress).Msg("Listening")
	err = http.ListenAndServe(ListenAddress, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start application")
	}
}

func main() {
	rootCmd.PersistentFlags().StringVar(&ListenAddress, "listen-address", ":9500", "The address this exporter would listen on")
	rootCmd.PersistentFlags().StringVar(&Currency, "currency", "USD", "Convert price value")
	rootCmd.PersistentFlags().StringVar(&CoingeckoApi, "coingecko-api", "https://api.coingecko.com/api/v3", "Coingecko API address")
	rootCmd.PersistentFlags().StringVar(&LogLevel, "log-level", "info", "Logging level")
	rootCmd.PersistentFlags().BoolVar(&JsonOutput, "json", false, "Output logs as JSON")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Could not start application")
	}
}
