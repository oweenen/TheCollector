package datastore

import (
	"TheCollectorDG/types"
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vmihailenco/msgpack/v5"
)

func storeComp(matchId string, comp *types.Comp) {
	msgpackData, err := msgpack.Marshal(comp)
	if err != nil {
		fmt.Println("Error serializing data to MessagePack:", err)
		return
	}

	bucket := "tft-stats-match-data"
	key := fmt.Sprintf("%s_%s.msgpack", matchId, comp.Summoner.Puuid)

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(msgpackData),
	})

	if err != nil {
		fmt.Printf("Error uploading %v to S3: %v\n", key, err)
		return
	}
}
