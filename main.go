package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
	"log"
	"io/ioutil"
	"fmt"
)

var (
	originRegion     = os.Getenv("ORIGIN_REGION")
	accessKeyId      = os.Getenv("ACCESS_KEY_ID")
	accessKeySecret  = os.Getenv("ACCESS_KEY_SECRET")
	originBucket     = os.Getenv("ORIGIN_BUCKET")
	targetRegion     = os.Getenv("TARGET_REGION")
	targetBucket     = os.Getenv("TARGET_BUCKET")
	targetEndpoint   = os.Getenv("TARGET_ENDPOINT")
)

func main() {
	originClient, err := oss.New(fmt.Sprintf("https://oss-%s.aliyuncs.com", originRegion), accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatalf("Error creating origin OSS client: %v", err)
	}

	targetClient, err := oss.New(fmt.Sprintf("https://oss-%s.aliyuncs.com", targetRegion), accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatalf("Error creating target OSS client: %v", err)
	}

	originBucket, err := originClient.Bucket(originBucket)
	if err != nil {
		log.Fatalf("Error getting origin bucket: %v", err)
	}

	targetBucket, err := targetClient.Bucket(targetBucket)
	if err != nil {
		log.Fatalf("Error getting target bucket: %v", err)
	}

	compareObjects(originBucket, targetBucket)
}

func compareObjects(originBucket, targetBucket *oss.Bucket) {
	marker := ""
	for {
		lor, err := originBucket.ListObjects(oss.Marker(marker), oss.MaxKeys(100))
		if err != nil {
			log.Fatalf("Error listing objects in origin bucket: %v", err)
		}

		for _, object := range lor.Objects {
			if object.Size > 0 {
				targetObject, err := targetBucket.GetObjectMeta(object.Key)
				if err != nil {
					if ossErr, ok := err.(oss.ServiceError); ok && ossErr.StatusCode == 404 {
						continue
					}
					log.Fatalf("Error getting object metadata from target bucket: %v", err)
				}

				originObject, err := originBucket.GetObjectMeta(object.Key)
				if err != nil {
					log.Fatalf("Error getting object metadata from origin bucket: %v", err)
				}

				if targetObject.Get("Content-Md5") != originObject.Get("Content-Md5") {
					log.Printf("Different content for object: %s", object.Key)
					writeToFile(fmt.Sprintf("Source URL: %s\nTarget URL: %s/%s\n\n", object.Key, targetEndpoint, object.Key))
				} else {
					log.Printf("Same content for object: %s", object.Key)
				}
			}
		}

		if !lor.IsTruncated {
			break
		}
		marker = lor.NextMarker
	}
}

func writeToFile(data string) {
	f, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	if _, err = f.WriteString(data); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}
