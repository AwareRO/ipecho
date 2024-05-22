package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	awahttp "github.com/AwareRO/libaware/golang/http"
	"github.com/AwareRO/libaware/golang/http/handlers"
	"github.com/AwareRO/libaware/golang/http/middlewares"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Http    awahttp.Config            `toml:"http" yaml:"http"`
	Metrics middlewares.MetricsConfig `toml:"metrics" yaml:"metrics"`
}

func ipv4(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	return ip
}

func handler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, ipv4(r))
}

func main() {
	path := "/etc/ipecho/backend.toml"
	if len(os.Args) != 1 {
		path = os.Args[1]
	}

	var config Config
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to read config")

		return
	}

	metricsWrapper := middlewares.NewDefaultDurationMetricWrapper(config.Metrics)
	router := httprouter.New()
	port := config.Http.Port

	router.GET("/", metricsWrapper.Wrap(handler))
	router.GET("/metrics", handlers.FromStdlib(metricsWrapper.Collector.GetHttpHandler()))

	log.Info().Int("port", port).Msg("Starting")

	awahttp.RunServer(&config.Http, router)
}
