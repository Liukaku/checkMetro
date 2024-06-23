package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func checkArr(result *s3.ListBucketsOutput) bool {
	if len(result.Buckets) == 0 {
		fmt.Println("You don't have any buckets!")
		return false
	} else {
		for _, bucket := range result.Buckets {
			if *bucket.Name == "metro-stops" {
				return true
			}
		}
		return false
	}
}

func createBucket(s3Client *s3.Client, bucketName string){
	createRes, err := s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintEuWest2,
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(createRes)
}

func loadDefaultSdkConfig() aws.Config{
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	
	if err != nil {
		panic(err)
	}

	return sdkConfig
}


func LoadConfig() {
	sdkConfig := loadDefaultSdkConfig()
	
	s3Client := s3.NewFromConfig(sdkConfig)

	result, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	if err != nil {
		panic(err)
	}

	hasMetroBucket := checkArr(result)
	fmt.Println(hasMetroBucket)

	if !hasMetroBucket {
		createBucket(s3Client, "metro-stops")
	}

}

func GetFromBucket(bucketName string, keyName string) io.ReadCloser{
	sdkConfig := loadDefaultSdkConfig()
	s3Client := s3.NewFromConfig(sdkConfig)
	key := keyName

	getResult, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key: &key,
	})

	if err != nil {
		panic(err)
	}

	return getResult.Body
}

func SaveToBucket(bucketName string, stopsObject []byte, keyName string){
	sdkConfig := loadDefaultSdkConfig()
	s3Client := s3.NewFromConfig(sdkConfig)
	createKey := keyName

	putResult, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key: &createKey,
		Body: bytes.NewReader(stopsObject),
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(putResult)
}