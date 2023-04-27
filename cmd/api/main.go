package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	envx "go.strv.io/env"
	zapx "go.strv.io/logging/zap"
	httpx "go.strv.io/net/http"
	"go.uber.org/zap"
	"io/fs"

	httpapi "newsletter-manager-go/api/rest"
	"newsletter-manager-go/database/sql"
	domauthor "newsletter-manager-go/domain/author"
	pgauthor "newsletter-manager-go/domain/author/postgres"
	domnewsletter "newsletter-manager-go/domain/newsletter"
	pgnewsletter "newsletter-manager-go/domain/newsletter/postgres"
	svcauthor "newsletter-manager-go/service/author"
	svcnewsletter "newsletter-manager-go/service/newsletter"
	"newsletter-manager-go/util"
	"newsletter-manager-go/util/timesource"
)

var (
	// version is set during the build.
	version          = "0.0.0"
	configPath       string
	integrationTests bool
)

const (
	defaultConfigPath = "./config.yaml"
)

type config struct {
	Port     uint            `json:"port" yaml:"port" env:"PORT" validate:"gt=0"`
	Database sql.Config      `json:"database" yaml:"database" env:",dive"`
	LogLevel zap.AtomicLevel `json:"log_level" yaml:"log_level" env:"LOG_LEVEL"`
}

func init() {
	pflag.StringVarP(&configPath, "config", "c", defaultConfigPath, "Path to configuration file")
	pflag.BoolVarP(&integrationTests, "integration", "i", false, "Mocked environment for integration tests")
	pflag.Parse()
}

func parseConfig(path string) (cfg config, err error) {
	viper.SetConfigFile(path)

	// Even if there is no config file, apply env variables and validate the config.
	defer func() {
		if err != nil {
			return
		}
		if envxErr := envx.Apply(&cfg); envxErr != nil {
			err = envxErr
			return
		}
		if validateErr := validator.New().Struct(&cfg); validateErr != nil {
			err = validateErr
			return
		}
	}()

	if err = viper.ReadInConfig(); err != nil {
		var e *fs.PathError
		if errors.As(err, &e) {
			return config{}, nil
		}
		return config{}, err
	}
	err = viper.Unmarshal(
		&cfg,
		func(dc *mapstructure.DecoderConfig) { dc.TagName = "yaml" },
		viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()),
	)
	if err != nil {
		return config{}, err
	}

	return cfg, nil
}

func getConnString() string {
	connString := "postgres://postgres:matchtheface123@localhost:5433/event-facematch?sslmode=disable"

	if integrationTests {
		// Integration tests need modified connection string without caching and with the special exec mode. This mode is slower than usual but it works well with the frequent DB schema changes
		connString += "&default_query_exec_mode=describe_exec"
	}

	return connString
}

func main() {
	// Parse config, set up logger.
	cfg, err := parseConfig(configPath)
	if err != nil {
		panic(fmt.Errorf("parse config: %w", err))
	}

	if integrationTests {
		_, _ = fmt.Println("Running in integration tests mode")
	}

	util.SetServerLogLevel(cfg.LogLevel)
	logger := zapx.MustCreateLogger(zapx.Config{
		Level:             cfg.LogLevel,
		DisableStacktrace: true,
	})

	ctx := context.Background()
	addr := fmt.Sprintf(":%d", cfg.Port)

	logger.With(
		zap.String("version", version),
		zap.String("addr", addr),
	).Info("starting application")

	connString := getConnString()

	database, err := sql.Open(ctx, connString)
	if err != nil {
		logger.Fatal("opening database", zap.Error(err))
	}

	authorFactory, err := domauthor.NewFactory(
		timesource.DefaultTimeSource{},
	)
	if err != nil {
		logger.Fatal("new author factory", zap.Error(err))
	}

	authorRepository, err := pgauthor.NewRepository(database, authorFactory)
	if err != nil {
		logger.Fatal("new author repository", zap.Error(err))
	}

	authorService, err := svcauthor.NewService(
		authorFactory,
		authorRepository,
	)
	if err != nil {
		logger.Fatal("new author service", zap.Error(err))
	}

	newsletterFactory, err := domnewsletter.NewFactory(
		timesource.DefaultTimeSource{},
	)
	if err != nil {
		logger.Fatal("new newsletter factory", zap.Error(err))
	}

	newsletterRepository, err := pgnewsletter.NewRepository(database, newsletterFactory)
	if err != nil {
		logger.Fatal("new newsletter repository", zap.Error(err))
	}

	newsletterService, err := svcnewsletter.NewService(
		newsletterFactory,
		newsletterRepository,
	)
	if err != nil {
		logger.Fatal("new newsletter service", zap.Error(err))
	}

	//var tokenParser middleware.TokenParser
	//if integrationTests {
	//	tokenParser = &mockTokenParser
	//} else {
	//	tokenParser = firebaseClient
	//}

	controller, err := httpapi.NewController(
		authorService,
		newsletterService,
		logger,
	)
	if err != nil {
		logger.Fatal("new HTTP controller", zap.Error(err))
	}

	// Run API server.
	serverConfig := httpx.ServerConfig{
		Addr:    addr,
		Handler: controller,
		Hooks: httpx.ServerHooks{
			BeforeShutdown: []httpx.ServerHookFunc{
				func(_ context.Context) {
					database.Close()
				},
			},
		},
		Limits: nil,
		Logger: util.NewServerLogger("httpx.Server"),
	}
	server := httpx.NewServer(&serverConfig)
	if err = server.Run(ctx); err != nil {
		logger.Fatal("HTTP server unexpectedly ended", zap.Error(err))
	}
}
