package secrets

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/stretchr/testify/assert"
)

type mockSecretsManagerClient struct {
	secretsmanageriface.SecretsManagerAPI
	getSecretValue func(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
}

func (m *mockSecretsManagerClient) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return m.getSecretValue(input)
}

func TestGetSecret(t *testing.T) {
	// Test successful case
	mockClient := &mockSecretsManagerClient{
		getSecretValue: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String("test-secret-value"),
			}, nil
		},
	}

	secretValue, err := getSecretWithClient("test-secret", mockClient)
	assert.NoError(t, err)
	assert.Equal(t, "test-secret-value", secretValue)

	// Test error case
	mockClient = &mockSecretsManagerClient{
		getSecretValue: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
			return nil, errors.New("test error")
		},
	}

	_, err = getSecretWithClient("test-secret", mockClient)
	assert.Error(t, err)
}

func getSecretWithClient(secretName string, smClient secretsmanageriface.SecretsManagerAPI) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := smClient.GetSecretValue(input)
	if err != nil {
		return "", fmt.Errorf("failed to get secret value: %v", err)
	}

	secretValue := aws.StringValue(result.SecretString)

	return secretValue, nil
}
