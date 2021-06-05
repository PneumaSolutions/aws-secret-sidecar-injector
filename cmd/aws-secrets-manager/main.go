package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {
	secretArn := os.Getenv("SECRET_ARN")
	secretFilename := os.Getenv("SECRET_FILENAME")
	var AWSRegion string

	if arn.IsARN(secretArn) {
		arnobj, _ := arn.Parse(secretArn)
		AWSRegion = arnobj.Region
	} else {
		log.Println("Not a valid ARN")
		os.Exit(1)
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Panic(err)
	}
	svc := secretsmanager.New(sess, &aws.Config{
		Region: aws.String(AWSRegion),
	})

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretArn),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeDecryptionFailure:
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return
	}
	// Decrypts secret using the associated KMS CMK.
	if result.SecretString != nil {
		writeStringOutput(*result.SecretString, secretFilename)
	} else {
		writeBinaryOutput(result.SecretBinary, secretFilename)
	}
}
func writeStringOutput(output string, name string) error {
	mountPoint := "/tmp"
	dir, file := filepath.Split(name)
	if file == "" {
		file = "secret"
	}
	err := os.MkdirAll(mountPoint + dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory, %w", err)
	}
	if filepath.IsAbs(filepath.Join(mountPoint + dir, file)) {
		f, err := os.Create(filepath.Join(mountPoint + dir, file))
		defer f.Close()
		if err != nil {
			return fmt.Errorf("error creating file, %w", err)
		}
		f.WriteString(output)
		return nil
	}
	return fmt.Errorf("not a valid file path")
}
func writeBinaryOutput(output []byte, name string) error {
	mountPoint := "/tmp"
	dir, file := filepath.Split(name)
	if file == "" {
		file = "secret"
	}
	err := os.MkdirAll(mountPoint + dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory, %w", err)
	}
	if filepath.IsAbs(filepath.Join(mountPoint + dir, file)) {
		f, err := os.Create(filepath.Join(mountPoint + dir, file))
		defer f.Close()
		if err != nil {
			return fmt.Errorf("error creating file, %w", err)
		}
		f.Write(output)
		return nil
	}
	return fmt.Errorf("not a valid file path")
}
