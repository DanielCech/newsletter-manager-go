package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	httpapi "strv-template-backend-go-api/api/rest"
	"strv-template-backend-go-api/api/rest/middleware"
	"strv-template-backend-go-api/crypto"
	"strv-template-backend-go-api/database/sql"
	domsession "strv-template-backend-go-api/domain/session"
	pgsession "strv-template-backend-go-api/domain/session/postgres"
	domuser "strv-template-backend-go-api/domain/user"
	pguser "strv-template-backend-go-api/domain/user/postgres"
	"strv-template-backend-go-api/metrics"
	secret "strv-template-backend-go-api/secret/aws"
	svcsession "strv-template-backend-go-api/service/session"
	svcuser "strv-template-backend-go-api/service/user"
	"strv-template-backend-go-api/util"
	"strv-template-backend-go-api/util/timesource"

	envx "go.strv.io/env"
	zapx "go.strv.io/logging/zap"
	httpx "go.strv.io/net/http"
	timex "go.strv.io/time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// version is set during the build.
	version    = "0.0.0"
	configPath string
)

const (
	defaultConfigPath = "./config.yaml"
)

type config struct {
	Port       uint       `json:"port" yaml:"port" env:"PORT" validate:"gt=0"`
	Database   sql.Config `json:"database" yaml:"database" env:",dive"`
	HashPepper string     `json:"hash_pepper" yaml:"hash_pepper" env:"HASH_PEPPER" validate:"gte=64"`
	AuthSecret string     `json:"auth_secret" yaml:"auth_secret" env:"AUTH_SECRET" validate:"gte=64"`
	Session    struct {
		AccessTokenExpiration  timex.Duration `json:"access_token_expiration" yaml:"access_token_expiration" env:"SESSION_ACCESS_TOKEN_EXPIRATION" validate:"required"`
		RefreshTokenExpiration timex.Duration `json:"refresh_token_expiration" yaml:"refresh_token_expiration" env:"SESSION_REFRESH_TOKEN_EXPIRATION" validate:"required"`
	} `json:"session" yaml:"session" env:",dive"`
	LogLevel       zap.AtomicLevel       `json:"log_level" yaml:"log_level" env:"LOG_LEVEL"`
	Metrics        metrics.Config        `json:"metrics" yaml:"metrics" env:",dive"`
	CORS           middleware.CORSConfig `json:"cors" yaml:"cors" env:",dive"`
	AWSEndpointURL *string               `json:"aws_endpoint_url" yaml:"aws_endpoint_url" env:"AWS_ENDPOINT_URL" validate:"omitempty,url"`
}

func init() {
	pflag.StringVarP(&configPath, "config", "c", defaultConfigPath, "Path to configuration file")
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

func loadAWSConfig(ctx context.Context, awsEndpointURL *string) (aws.Config, error) {
	return awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(string, string, ...any) (aws.Endpoint, error) {
				if awsEndpointURL != nil {
					return aws.Endpoint{
						URL:           *awsEndpointURL,
						PartitionID:   "aws-local-stack",
						SigningRegion: "us-east-1",
					}, nil
				}
				// Fallback to default AWS endpoint.
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			}),
		),
	)
}

func setupDatabase(ctx context.Context, cfg config) (sql.Database, error) {
	awsCfg, err := loadAWSConfig(ctx, cfg.AWSEndpointURL)
	if err != nil {
		return sql.Database{}, fmt.Errorf("loading aws config: %w", err)
	}

	secretsManager := secret.NewSecretsManager(awsCfg)
	databaseSecret := secret.NewSecret[sql.DSNValues](cfg.Database.Secret)
	dsnValues, err := databaseSecret.Resolve(ctx, secretsManager)
	if err != nil {
		return sql.Database{}, fmt.Errorf("resolving database configuration: %w", err)
	}

	connString, err := dsnValues.ConnString()
	if err != nil {
		return sql.Database{}, fmt.Errorf("parsing database connection string: %w", err)
	}
	database, err := sql.Open(ctx, connString)
	if err != nil {
		return sql.Database{}, fmt.Errorf("opening database: %w", err)
	}

	return database, nil
}

func setupUserServiceDeps(database sql.Database, cfg config) (domuser.Factory, domuser.Repository, error) {
	userFactory, err := domuser.NewFactory(
		crypto.NewDefaultBcryptHasher([]byte(cfg.HashPepper)),
		timesource.DefaultTimeSource{},
	)
	if err != nil {
		return domuser.Factory{}, nil, fmt.Errorf("new user factory: %w", err)
	}

	userRepository, err := pguser.NewRepository(database, userFactory)
	if err != nil {
		return domuser.Factory{}, nil, fmt.Errorf("new user repository: %w", err)
	}

	return userFactory, userRepository, nil
}

func setupSessionServiceDeps(database sql.Database, cfg config) (domsession.Factory, domsession.Repository, error) {
	sessionFactory, err := domsession.NewFactory(
		[]byte(cfg.AuthSecret),
		timesource.DefaultTimeSource{},
		cfg.Session.AccessTokenExpiration.Duration(),
		cfg.Session.RefreshTokenExpiration.Duration(),
	)
	if err != nil {
		return domsession.Factory{}, nil, fmt.Errorf("new session factory: %w", err)
	}

	sessionRepository, err := pgsession.NewRepository(database, sessionFactory)
	if err != nil {
		return domsession.Factory{}, nil, fmt.Errorf("new session repository: %w", err)
	}

	return sessionFactory, sessionRepository, nil
}

// TODO: As can be seen, dependency injection is sort of hacked. Waiting for resolving this big issue.
func setupController(database sql.Database, logger *zap.Logger, cfg config) (*httpapi.Controller, error) {
	userService := new(svcuser.Service)
	sessionService := new(svcsession.Service)

	userFactory, userRepository, err := setupUserServiceDeps(database, cfg)
	if err != nil {
		return nil, fmt.Errorf("setup user service dependencies: %w", err)
	}
	userServiceTmp, err := svcuser.NewService(
		userFactory,
		userRepository,
		sessionService,
	)
	if err != nil {
		return nil, fmt.Errorf("new user service: %w", err)
	}
	*userService = *userServiceTmp

	sessionFactory, sessionRepository, err := setupSessionServiceDeps(database, cfg)
	if err != nil {
		return nil, fmt.Errorf("setup session service dependencies: %w", err)
	}
	sessionServiceTmp, err := svcsession.NewService(
		sessionFactory,
		sessionRepository,
		userService,
	)
	if err != nil {
		return nil, fmt.Errorf("new session service: %w", err)
	}
	*sessionService = *sessionServiceTmp

	controller, err := httpapi.NewController(
		userService,
		sessionService,
		sessionFactory,
		cfg.CORS,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("new http controller: %w", err)
	}

	metrics.MustRegister(userService, sessionService)
	return controller, nil
}

func main() {
	// Parse config, set up logger.
	cfg, err := parseConfig(configPath)
	if err != nil {
		panic(fmt.Errorf("parse config: %w", err))
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

	database, err := setupDatabase(ctx, cfg)
	if err != nil {
		logger.Fatal("setup database", zap.Error(err))
	}

	controller, err := setupController(database, logger, cfg)
	if err != nil {
		logger.Fatal("setup controller", zap.Error(err))
	}

	go func() {
		if err := metrics.NewServer(cfg.Metrics).Run(ctx); err != nil {
			logger.Error("metrics HTTP server unexpectedly ended", zap.Error(err))
		}
	}()

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
