package datastore

import (
	"TheCollectorDG/types"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func storeComp(matchId string, comp *types.Comp) {
	jsonData, err := json.Marshal(comp)
	if err != nil {
		fmt.Println("Error marshaling data to json:", err)
		return
	}

	bucket := "tft-stats-comps"
	key := fmt.Sprintf("%s/%s.json", matchId, comp.Summoner.Puuid)

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		Body:     bytes.NewReader(jsonData),
		ACL:      aws.String("public-read"),
		Metadata: map[string]*string{},
	})

	if err != nil {
		fmt.Printf("Error uploading %v to spaces: %v\n", key, err)
		return
	}
}
