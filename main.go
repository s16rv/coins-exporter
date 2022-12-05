package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/s16rv/coins-exporter/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	ListenAddress string
	CoingeckoApi  string
	LogLevel      string
	JsonOutput    bool
	ConstLabels   map[string]string
)

var log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

var rootCmd = &cobra.Command{
	Use:  "coins-exporter",
	Long: "Scrape the data about token price in the Coingecko.",
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

func SendQuery(baseApi, id string) (types.ReturnData, error) {
	u := baseApi + "/coins/" + id
	var d types.ReturnData

	client := &http.Client{}

	req, _ := http.NewRequest("GET", u, nil)
	resp, err := client.Do(req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(body, &d)
	if err != nil {
		return d, err
	}
	defer resp.Body.Close()

	return d, nil
}

func main() {
	rootCmd.PersistentFlags().StringVar(&ListenAddress, "listen-address", ":9500", "The address this exporter would listen on")
	rootCmd.PersistentFlags().StringVar(&CoingeckoApi, "coingecko-api", "https://api.coingecko.com/api/v3", "Coingecko API address")
	rootCmd.PersistentFlags().StringVar(&LogLevel, "log-level", "info", "Logging level")
	rootCmd.PersistentFlags().BoolVar(&JsonOutput, "json", false, "Output logs as JSON")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Could not start application")
	}
}
