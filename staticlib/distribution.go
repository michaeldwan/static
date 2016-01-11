package staticlib

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/michaeldwan/static/printer"
)

type Distribution struct {
	ID         string
	DomainName string
	Aliases    []*string
}

func newDistributionFromSummary(summary *cloudfront.DistributionSummary) Distribution {
	return Distribution{
		ID:         *summary.ID,
		DomainName: *summary.DomainName,
		Aliases:    summary.Aliases.Items,
	}
}

func findDistributionsForOrigin(originHost string) []Distribution {
	var out []Distribution

	cfClient.ListDistributionsPages(nil, func(page *cloudfront.ListDistributionsOutput, lastPage bool) bool {
		for _, distroSummary := range page.DistributionList.Items {
			for _, origin := range distroSummary.Origins.Items {
				if *origin.DomainName == originHost {
					out = append(out, newDistributionFromSummary(distroSummary))
				}
			}
		}
		return true
	})

	return out
}

func (d Distribution) invalidate(m *Manifest) error {
	var items []*string
	for _, path := range invalidationPathsFromManifest(m) {
		items = append(items, aws.String(path))
	}
	invalidationID := fmt.Sprintf("static-%x", m.digest)
	input := &cloudfront.CreateInvalidationInput{
		DistributionID: aws.String(d.ID),
		InvalidationBatch: &cloudfront.InvalidationBatch{
			CallerReference: aws.String(invalidationID),
			Paths: &cloudfront.Paths{
				Quantity: aws.Long(int64(len(items))),
				Items:    items,
			},
		},
	}
	_, err := cfClient.CreateInvalidation(input)
	if err != nil {
		// printAWSError(err)
		return err
	}

	printer.Infof("Invalidated %d paths for %v\n", len(items), d.DomainName)
	for _, cname := range d.Aliases {
		printer.Infoln("  -> ", *cname)
	}
	return nil
}

func invalidationPathsFromManifest(m *Manifest) []string {
	pathSet := make(map[string]bool)
	// TODO: modes <precise, directory>
	for _, e := range m.entriesForPushActions(Update, ForceUpdate, Create) {
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
