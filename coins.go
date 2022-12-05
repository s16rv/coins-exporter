package main

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func CoinsHandler(w http.ResponseWriter, r *http.Request, baseApi string) {
	requestStart := time.Now()
	sublogger := log.With().
		Str("request-id", uuid.New().String()).
		Logger()

	coinIdsRaw := r.URL.Query().Get("ids")
	coinIds := strings.Split(coinIdsRaw, ",")

	coinsPriceGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "coins_price",
			Help:        "Price of Coins in currency",
			ConstLabels: ConstLabels,
		},
		[]string{"id", "symbol", "currency", "name"},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(coinsPriceGauge)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, coinId := range coinIds {
			sublogger.Debug().
				Str("id", coinId).
				Msg("Started querying coins by id")
			queryStart := time.Now()

			resp, err := SendQueryCoinsDetail(baseApi, coinId)

			if err != nil {
				sublogger.Error().
					Str("id", coinId).
					Err(err).
					Msg("Could not get coins by id")
				return
			}

			price := resp.MarketData.CurrentPrice.GetField(Currency)

			coinsPriceGauge.With(prometheus.Labels{
				"id":       coinId,
				"symbol":   resp.Symbol,
				"name":     resp.Name,
				"currency": Currency,
			}).Set(price)

			sublogger.Debug().
				Str("id", coinId).
				Float64("request-time", time.Since(queryStart).Seconds()).
				Msg("Finished querying coins by id")
		}
	}()

	wg.Wait()

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	sublogger.Info().
		Str("method", "GET").
		Str("endpoint", "/metrics/coins?ids="+coinIdsRaw).
		Float64("request-time", time.Since(requestStart).Seconds()).
		Msg("Request processed")
}
