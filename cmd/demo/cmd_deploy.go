package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVar(&deployRegion, "region", "eu-central-1", "AWS region")
	deployCmd.Flags().StringVar(&deployFunctionName, "function", "demo-app-apiv1", "Lambda function name")
	deployCmd.Flags().BoolVar(&deploySkipBuild, "skip-build", false, "Skip building the Lambda function")
}

var (
	deployRegion       string
	deployFunctionName string
	deploySkipBuild    bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the application to AWS Lambda",
	Long:  `Build and deploy the demo application to AWS Lambda environment`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Get project root (two levels up from cmd/demo)
		executable, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		projectRoot := filepath.Join(executable)
		lambdaDir := filepath.Join(projectRoot, "cmd", "demo-aws-lambda")

		if !deploySkipBuild {
			slog.Info("Building Lambda function...")
			if err := buildLambdaFunction(lambdaDir); err != nil {
				return fmt.Errorf("failed to build Lambda function: %w", err)
			}
			slog.Info("Build completed successfully")
		} else {
			slog.Info("Skipping build as requested")
		}

		zipFile := filepath.Join(lambdaDir, "demo-aws-lambda.zip")
		if !deploySkipBuild {
			slog.Info("Creating deployment package...")
			if err := createDeploymentPackage(lambdaDir, zipFile); err != nil {
				return fmt.Errorf("failed to create deployment package: %w", err)
			}
			slog.Info("Deployment package created")
		}

		slog.Info(fmt.Sprintf("Deploying to AWS Lambda function: %s in region: %s", deployFunctionName, deployRegion))
		functionURL, err := deployToAWS(ctx, zipFile, deployFunctionName, deployRegion)
		if err != nil {
			return fmt.Errorf("failed to deploy to AWS: %w", err)
		}

		slog.Info(fmt.Sprintf("Deployment successful! Function URL: %s", functionURL))
		return nil
	},
}

func buildLambdaFunction(lambdaDir string) error {
	cmd := exec.Command("go", "build", "-v", "-tags", "lambda.norpc", "-o", "bootstrap")
	cmd.Dir = lambdaDir
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=0",
		"GOOS=linux",
		"GOARCH=amd64",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func createDeploymentPackage(lambdaDir, zipFile string) error {
	// Create zip file
	zipFileHandle, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zipFileHandle.Close()

	zipWriter := zip.NewWriter(zipFileHandle)
	defer zipWriter.Close()

	// Add bootstrap binary
	if err := addFileToZip(zipWriter, filepath.Join(lambdaDir, "bootstrap"), "bootstrap"); err != nil {
		return fmt.Errorf("failed to add bootstrap to zip: %w", err)
	}

	// Add env.txt if it exists
	envFile := filepath.Join(lambdaDir, "env.txt")
	if _, err := os.Stat(envFile); err == nil {
		if err := addFileToZip(zipWriter, envFile, "env.txt"); err != nil {
			return fmt.Errorf("failed to add env.txt to zip: %w", err)
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename, nameInZip string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = nameInZip
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

func deployToAWS(ctx context.Context, zipFile, functionName, region string) (string, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Lambda client
	lambdaSvc := lambda.NewFromConfig(cfg)

	// Read zip file
	zipData, err := os.ReadFile(zipFile)
	if err != nil {
		return "", fmt.Errorf("failed to read zip file: %w", err)
	}

	// Update function code
	updateOutput, err := lambdaSvc.UpdateFunctionCode(ctx, &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(functionName),
		ZipFile:      zipData,
	})
	if err != nil {
		return "", fmt.Errorf("failed to update Lambda function code: %w", err)
	}

	slog.Info(fmt.Sprintf("Function updated: %s (version: %s)", *updateOutput.FunctionName, *updateOutput.Version))

	// Get function URL
	urlConfig, err := lambdaSvc.GetFunctionUrlConfig(ctx, &lambda.GetFunctionUrlConfigInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		slog.Warn(fmt.Sprintf("Could not retrieve function URL: %v", err))
		return "", nil
	}

	return *urlConfig.FunctionUrl, nil
}
