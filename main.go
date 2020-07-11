package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

const AWS_ID = "<Your AWS ID>"
const AWS_SECRET = "<Your AWS Secret>"
const AWS_TOKEN = "<Your AWS Token>"
const AWS_REGION = "<Your AWS Region>"

func main() {
	svc := createAwsSession()

	fileName := "zippedFile.zip"
	zipFileContents, err := getZipFileContents(fileName)
	if err != nil {
		return
	}

	functionName := "<Your function name>"

	err = createLambdaFunction(svc, zipFileContents, functionName)
	if err != nil {
		fmt.Println("Failed to create lambda function")
	}
	response, err := invokeLambdaFunction(svc, functionName)
	if err != nil {
		fmt.Println("Failed to execute lambda function")
		return
	}

	fmt.Println(response)

}

func createAwsSession() *lambda.Lambda {
	creds := credentials.NewStaticCredentials(AWS_ID, AWS_SECRET, AWS_TOKEN)

	svc := lambda.New(session.New(
		&aws.Config{
			Credentials: creds,
			Region:      aws.String(AWS_REGION),
		}))

	return svc
}

func getZipFileContents(zipFileName string) ([]byte, error) {
	file, err := os.Open(zipFileName)
	if err != nil {
		return nil, err
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileContents, nil
}

func invokeLambdaFunction(svc *lambda.Lambda, functionName string) ([]byte, error) {
	input := &lambda.InvokeInput{
		FunctionName: aws.String(functionName),
		Payload:      []byte("{}"), // Can modify the payload to your usecase.
	}

	result, err := svc.Invoke(input)
	if err != nil {
		return nil, err
	}
	return result.Payload, err
}

func createLambdaFunction(svc *lambda.Lambda, zipFileContents []byte, functionName string) error {
	input := &lambda.CreateFunctionInput{
		Code: &lambda.FunctionCode{
			ZipFile: zipFileContents,
		},
		Description: aws.String("Your Code description"),
		Environment: &lambda.Environment{
			Variables: map[string]*string{
				"BUCKET": aws.String("<Yout Bucket Name>"),
				"PREFIX": aws.String("inbound"),
			},
		},
		FunctionName: aws.String(functionName),
		// The code file should export handler method from index file.
		Handler:    aws.String("index.handler"),
		MemorySize: aws.Int64(256),
		Publish:    aws.Bool(true),
		Role:       aws.String("<Your AWS Role>"),
		Runtime:    aws.String("nodejs12.x"),
		Tags: map[string]*string{
			"DEPARTMENT": aws.String("Assets"),
		},
		Timeout: aws.Int64(15),
		TracingConfig: &lambda.TracingConfig{
			Mode: aws.String("Active"),
		},
	}

	_, err := svc.CreateFunction(input)
	return err
}
