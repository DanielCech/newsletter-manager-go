package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	ErrNoSecretValue    = errors.New("no secret value")
	ErrMissingPathOrARN = errors.New("invalid config: at least one of path or arn must be present")
)

type SecretsManager struct {
	client *secretsmanager.Client
}

func NewSecretsManager(awsCfg aws.Config) SecretsManager {
	return SecretsManager{
		client: secretsmanager.NewFromConfig(awsCfg),
	}
}

// GetAndParseSecretValue retrieves and parses secret value (JSON object) from AWS secrets manager.
func (sm SecretsManager) GetAndParseSecretValue(ctx context.Context, secretARN string, target any) error {
	result, err := sm.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretARN),
	})
	if err != nil {
		return fmt.Errorf("getting secret value: %w", err)
	}
	if result.SecretString == nil {
		return ErrNoSecretValue
	}
	if err = json.Unmarshal([]byte(*result.SecretString), target); err != nil {
		return fmt.Errorf("unmarshaling secret data: %w", err)
	}
	return nil
}
