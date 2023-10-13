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

	bucket := "tft-stats-comps"
	key := fmt.Sprintf("%s/%s.msgpack", matchId, comp.Summoner.Puuid)

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		Body:     bytes.NewReader(msgpackData),
		ACL:      aws.String("public-read"),
		Metadata: map[string]*string{},
	})

	if err != nil {
		fmt.Printf("Error uploading %v to spaces: %v\n", key, err)
		return
	}
}
