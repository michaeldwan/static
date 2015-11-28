package dest

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

var cf *cloudfront.CloudFront

func cloudFrontClent() *cloudfront.CloudFront {
	if cf == nil {
		cf = cloudfront.New(&aws.Config{
			Credentials: getCredentials(),
			// LogLevel: 1,
		})
	}
	return cf
}

func CreateInvalidations(bucket *Bucket, keys []string) {
	distroIds := findDistributionIdsForBucket(bucket)
	if len(distroIds) == 0 {
		return
	}
	fmt.Println("keys:", keys)
	paths := determineInvalidationPaths(keys)
	if len(paths) == 0 {
		return
	}
	fmt.Println("paths:", paths)

	for _, id := range distroIds {
		invalidateDistribution(id, paths)
	}
}

func findDistributionIdsForBucket(bucket *Bucket) []string {
	var ids []string
	params := &cloudfront.ListDistributionsInput{}
	response, err := cloudFrontClent().ListDistributions(params)
	if err != nil {
		log.Fatal(err)
	}
	for _, summary := range response.DistributionList.Items {
		if isOriginForDistro(bucket, summary) {
			ids = append(ids, *summary.ID)
		}
	}
	return ids
}

func isOriginForDistro(bucket *Bucket, distro *cloudfront.DistributionSummary) bool {
	// TODO: detect & warn if origins are configured w/ s3 url rather than s3 website url
	for _, origin := range distro.Origins.Items {
		return *origin.DomainName == bucket.WebsiteEndpoint()
	}
	return false
}

func invalidateDistribution(id string, paths []string) {
	log.Println("Invalidate", id, " paths:", paths)
	items := make([]*string, 0)

	for _, path := range paths {
		items = append(items, aws.String(path))
	}

	fmt.Println("invalidationid", invalidationId())

	params := &cloudfront.CreateInvalidationInput{
		DistributionID: aws.String(id),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(invalidationId()),
			Paths: &cloudfront.Paths{
				Quantity: aws.Long(int64(len(items))),
				Items:    items,
			},
		},
	}
	resp, err := cloudFrontClent().CreateInvalidation(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS error with Code, Message, and original error (if any)
			fmt.Println(awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				fmt.Println(reqErr.Code(), reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())
			}
		} else {
			// This case should never be hit, the SDK should always return an
			// error which satisfies the awserr.Error interface.
			fmt.Println(err.Error())
		}
	}

	// Pretty-print the response data.
	fmt.Println(awsutil.StringValue(resp))
}

func invalidationId() string {
	epoch := time.Now().Unix()
	// TODO: use a timestamp and digest of paths for the id
	return fmt.Sprintf("static-%d", epoch)
}

func determineInvalidationPaths(keys []string) []string {
	pathSet := make(map[string]bool)
	// TODO: modes <precise, directory>
	for _, key := range keys {
		path := "/" + key + "*"
		pathSet[path] = true
	}
	paths := make([]string, 0)
	for key := range pathSet {
		paths = append(paths, key)
	}
	return paths
}
