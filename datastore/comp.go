package datastore

import (
	"TheCollectorDG/types"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func storeComp(matchId string, comp *types.Comp) error {
	jsonData, err := json.Marshal(comp)
	if err != nil {
		return err
	}

	bucket := "tft-stats-comps"
	key := fmt.Sprintf("%s/%s.json", matchId, comp.SummonerPuuid)

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(jsonData),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("application/json"),
		Metadata:    map[string]*string{},
	})

	return err
}
