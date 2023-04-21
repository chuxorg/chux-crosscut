package secrets

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Get returns the value of the secret with the given name.
// Example:
// secretName := "dev/secrets/SOME_SECRET_KEY"
//
//	secretValue, err := getSecret(secretName)
//	if err != nil {
//		log.Fatalf("Error getting secret: %v", err)
//	}
//	fmt.Printf("Secret value: %s\n", secretValue)
func GetSecret(secretName string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %v", err)
	}

	smClient := secretsmanager.New(sess)

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

// GetAll returns a map of all secrets in the secrets manager.
// Example:
// secrets, err := GetAll()
//
//	if err != nil {
//		log.Fatalf("Error getting secrets: %v", err)
//	}
//
//	for name, value := range secrets {
//		fmt.Printf("Secret name: %s, value: %s\n", name, value)
//	}
func GetAll() (map[string]string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	smClient := secretsmanager.New(sess)

	input := &secretsmanager.ListSecretsInput{}

	result, err := smClient.ListSecrets(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %v", err)
	}

	secrets := make(map[string]string)

	for _, secret := range result.SecretList {
		secretValue, err := GetSecret(aws.StringValue(secret.Name))
		if err != nil {
			return nil, fmt.Errorf("failed to get secret value: %v", err)
		}

		secrets[aws.StringValue(secret.Name)] = secretValue
	}

	return secrets, nil
}
