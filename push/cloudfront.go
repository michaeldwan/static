package push

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/michaeldwan/static/printer"
)

func createInvalidations(cf *cloudfront.CloudFront, region, bucket string, manifest *Manifest) {
	keys := invalidationPathsFromManifest(manifest)
	if len(keys) == 0 {
		return
	}
	for _, distroSummary := range distributionsForBucket(cf, region, bucket) {
		if err := invalidateDistribution(cf, distroSummary, manifest.digest(), keys); err != nil {
			panic(err)
		}
	}
}

func distributionsForBucket(cf *cloudfront.CloudFront, region, bucket string) []*cloudfront.DistributionSummary {
	var distroSummaries []*cloudfront.DistributionSummary
	s3Endpoint := s3WebsiteEndpoint(region, bucket)
	cf.ListDistributionsPages(nil, func(page *cloudfront.ListDistributionsOutput, lastPage bool) bool {
		for _, distroSummary := range page.DistributionList.Items {
			for _, origin := range distroSummary.Origins.Items {
				if *origin.DomainName == s3Endpoint {
					distroSummaries = append(distroSummaries, distroSummary)
				}
			}
		}
		return true
	})
	return distroSummaries
}

func invalidationPathsFromManifest(m *Manifest) []string {
	pathSet := make(map[string]bool)
	// TODO: modes <precise, directory>
	for _, e := range m.entriesForOperations(Update, ForceUpdate, Create) {
		// path := "/" + e.Key + "*"
		path := "/" + e.Key
		pathSet[path] = true
	}
	var paths []string
	for key := range pathSet {
		paths = append(paths, key)
	}
	return paths
}

func invalidateDistribution(cf *cloudfront.CloudFront, distro *cloudfront.DistributionSummary, digest [16]byte, paths []string) error {
	var items []*string
	for _, path := range paths {
		items = append(items, aws.String(path))
	}
	invalidationID := fmt.Sprintf("static-%x", digest)
	input := &cloudfront.CreateInvalidationInput{
		DistributionID: distro.ID,
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(invalidationID),
			Paths: &cloudfront.Paths{
				Quantity: aws.Long(int64(len(items))),
				Items:    items,
			},
		},
	}
	_, err := cf.CreateInvalidation(input)
	if err != nil {
		// printAWSError(err)
		return err
	}

	printer.Infof("Invalidated %d paths for %v\n", len(paths), *distro.DomainName)
	for _, cname := range distro.Aliases.Items {
		printer.Infoln("  -> ", *cname)
	}
	return nil
}
