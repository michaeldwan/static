// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package directoryservice_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/directoryservice"
)

var _ time.Duration
var _ bytes.Buffer

func ExampleDirectoryService_ConnectDirectory() {
	svc := directoryservice.New(nil)

	params := &directoryservice.ConnectDirectoryInput{
		ConnectSettings: &directoryservice.DirectoryConnectSettings{ // Required
			CustomerDNSIPs: []*string{ // Required
				aws.String("IpAddr"), // Required
				// More values...
			},
			CustomerUserName: aws.String("UserName"), // Required
			SubnetIDs: []*string{ // Required
				aws.String("SubnetId"), // Required
				// More values...
			},
			VPCID: aws.String("VpcId"), // Required
		},
		Name:        aws.String("DirectoryName"),   // Required
		Password:    aws.String("ConnectPassword"), // Required
		Size:        aws.String("DirectorySize"),   // Required
		Description: aws.String("Description"),
		ShortName:   aws.String("DirectoryShortName"),
	}
	resp, err := svc.ConnectDirectory(params)

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

func ExampleDirectoryService_CreateAlias() {
	svc := directoryservice.New(nil)

	params := &directoryservice.CreateAliasInput{
		Alias:       aws.String("AliasName"),   // Required
		DirectoryID: aws.String("DirectoryId"), // Required
	}
	resp, err := svc.CreateAlias(params)

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

func ExampleDirectoryService_CreateComputer() {
	svc := directoryservice.New(nil)

	params := &directoryservice.CreateComputerInput{
		ComputerName: aws.String("ComputerName"),     // Required
		DirectoryID:  aws.String("DirectoryId"),      // Required
		Password:     aws.String("ComputerPassword"), // Required
		ComputerAttributes: []*directoryservice.Attribute{
			{ // Required
				Name:  aws.String("AttributeName"),
				Value: aws.String("AttributeValue"),
			},
			// More values...
		},
		OrganizationalUnitDistinguishedName: aws.String("OrganizationalUnitDN"),
	}
	resp, err := svc.CreateComputer(params)

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

func ExampleDirectoryService_CreateDirectory() {
	svc := directoryservice.New(nil)

	params := &directoryservice.CreateDirectoryInput{
		Name:        aws.String("DirectoryName"), // Required
		Password:    aws.String("Password"),      // Required
		Size:        aws.String("DirectorySize"), // Required
		Description: aws.String("Description"),
		ShortName:   aws.String("DirectoryShortName"),
		VPCSettings: &directoryservice.DirectoryVPCSettings{
			SubnetIDs: []*string{ // Required
				aws.String("SubnetId"), // Required
				// More values...
			},
			VPCID: aws.String("VpcId"), // Required
		},
	}
	resp, err := svc.CreateDirectory(params)

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

func ExampleDirectoryService_CreateSnapshot() {
	svc := directoryservice.New(nil)

	params := &directoryservice.CreateSnapshotInput{
		DirectoryID: aws.String("DirectoryId"), // Required
		Name:        aws.String("SnapshotName"),
	}
	resp, err := svc.CreateSnapshot(params)

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

func ExampleDirectoryService_DeleteDirectory() {
	svc := directoryservice.New(nil)

	params := &directoryservice.DeleteDirectoryInput{
		DirectoryID: aws.String("DirectoryId"), // Required
	}
	resp, err := svc.DeleteDirectory(params)

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

func ExampleDirectoryService_DeleteSnapshot() {
	svc := directoryservice.New(nil)

	params := &directoryservice.DeleteSnapshotInput{
		SnapshotID: aws.String("SnapshotId"), // Required
	}
	resp, err := svc.DeleteSnapshot(params)

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

func ExampleDirectoryService_DescribeDirectories() {
	svc := directoryservice.New(nil)

	params := &directoryservice.DescribeDirectoriesInput{
		DirectoryIDs: []*string{
			aws.String("DirectoryId"), // Required
			// More values...
		},
		Limit:     aws.Long(1),
		NextToken: aws.String("NextToken"),
	}
	resp, err := svc.DescribeDirectories(params)

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

func ExampleDirectoryService_DescribeSnapshots() {
	svc := directoryservice.New(nil)

	params := &directoryservice.DescribeSnapshotsInput{
		DirectoryID: aws.String("DirectoryId"),
		Limit:       aws.Long(1),
		NextToken:   aws.String("NextToken"),
		SnapshotIDs: []*string{
			aws.String("SnapshotId"), // Required
			// More values...
		},
	}
	resp, err := svc.DescribeSnapshots(params)

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

func ExampleDirectoryService_DisableRadius() {
	svc := directoryservice.New(nil)

	params := &directoryservice.DisableRadiusInput{
		DirectoryID: aws.String("DirectoryId"), // Required
	}
	resp, err := svc.DisableRadius(params)

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

func ExampleDirectoryService_DisableSSO() {
	svc := directoryservice.New(nil)

	params := &directoryservice.DisableSSOInput{
		DirectoryID: aws.String("DirectoryId"), // Required
		Password:    aws.String("ConnectPassword"),
		UserName:    aws.String("UserName"),
	}
	resp, err := svc.DisableSSO(params)

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

func ExampleDirectoryService_EnableRadius() {
	svc := directoryservice.New(nil)

	params := &directoryservice.EnableRadiusInput{
		DirectoryID: aws.String("DirectoryId"), // Required
		RadiusSettings: &directoryservice.RadiusSettings{ // Required
			AuthenticationProtocol: aws.String("RadiusAuthenticationProtocol"),
			DisplayLabel:           aws.String("RadiusDisplayLabel"),
			RadiusPort:             aws.Long(1),
			RadiusRetries:          aws.Long(1),
			RadiusServers: []*string{
				aws.String("Server"), // Required
				// More values...
			},
			RadiusTimeout:   aws.Long(1),
			SharedSecret:    aws.String("RadiusSharedSecret"),
			UseSameUsername: aws.Boolean(true),
		},
	}
	resp, err := svc.EnableRadius(params)

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

func ExampleDirectoryService_EnableSSO() {
	svc := directoryservice.New(nil)

	params := &directoryservice.EnableSSOInput{
		DirectoryID: aws.String("DirectoryId"), // Required
		Password:    aws.String("ConnectPassword"),
		UserName:    aws.String("UserName"),
	}
	resp, err := svc.EnableSSO(params)

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

func ExampleDirectoryService_GetDirectoryLimits() {
	svc := directoryservice.New(nil)

	var params *directoryservice.GetDirectoryLimitsInput
	resp, err := svc.GetDirectoryLimits(params)

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

func ExampleDirectoryService_GetSnapshotLimits() {
	svc := directoryservice.New(nil)

	params := &directoryservice.GetSnapshotLimitsInput{
		DirectoryID: aws.String("DirectoryId"), // Required
	}
	resp, err := svc.GetSnapshotLimits(params)

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

func ExampleDirectoryService_RestoreFromSnapshot() {
	svc := directoryservice.New(nil)

	params := &directoryservice.RestoreFromSnapshotInput{
		SnapshotID: aws.String("SnapshotId"), // Required
	}
	resp, err := svc.RestoreFromSnapshot(params)

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

func ExampleDirectoryService_UpdateRadius() {
	svc := directoryservice.New(nil)

	params := &directoryservice.UpdateRadiusInput{
		DirectoryID: aws.String("DirectoryId"), // Required
		RadiusSettings: &directoryservice.RadiusSettings{ // Required
			AuthenticationProtocol: aws.String("RadiusAuthenticationProtocol"),
			DisplayLabel:           aws.String("RadiusDisplayLabel"),
			RadiusPort:             aws.Long(1),
			RadiusRetries:          aws.Long(1),
			RadiusServers: []*string{
				aws.String("Server"), // Required
				// More values...
			},
			RadiusTimeout:   aws.Long(1),
			SharedSecret:    aws.String("RadiusSharedSecret"),
			UseSameUsername: aws.Boolean(true),
		},
	}
	resp, err := svc.UpdateRadius(params)

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
