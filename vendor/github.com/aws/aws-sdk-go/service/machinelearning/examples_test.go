// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

package machinelearning_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/service/machinelearning"
)

var _ time.Duration
var _ bytes.Buffer

func ExampleMachineLearning_CreateBatchPrediction() {
	svc := machinelearning.New(nil)

	params := &machinelearning.CreateBatchPredictionInput{
		BatchPredictionDataSourceID: aws.String("EntityId"), // Required
		BatchPredictionID:           aws.String("EntityId"), // Required
		MLModelID:                   aws.String("EntityId"), // Required
		OutputURI:                   aws.String("S3Url"),    // Required
		BatchPredictionName:         aws.String("EntityName"),
	}
	resp, err := svc.CreateBatchPrediction(params)

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

func ExampleMachineLearning_CreateDataSourceFromRDS() {
	svc := machinelearning.New(nil)

	params := &machinelearning.CreateDataSourceFromRDSInput{
		DataSourceID: aws.String("EntityId"), // Required
		RDSData: &machinelearning.RDSDataSpec{ // Required
			DatabaseCredentials: &machinelearning.RDSDatabaseCredentials{ // Required
				Password: aws.String("RDSDatabasePassword"), // Required
				Username: aws.String("RDSDatabaseUsername"), // Required
			},
			DatabaseInformation: &machinelearning.RDSDatabase{ // Required
				DatabaseName:       aws.String("RDSDatabaseName"),       // Required
				InstanceIdentifier: aws.String("RDSInstanceIdentifier"), // Required
			},
			ResourceRole:      aws.String("EDPResourceRole"), // Required
			S3StagingLocation: aws.String("S3Url"),           // Required
			SecurityGroupIDs: []*string{ // Required
				aws.String("EDPSecurityGroupId"), // Required
				// More values...
			},
			SelectSQLQuery:    aws.String("RDSSelectSqlQuery"), // Required
			ServiceRole:       aws.String("EDPServiceRole"),    // Required
			SubnetID:          aws.String("EDPSubnetId"),       // Required
			DataRearrangement: aws.String("DataRearrangement"),
			DataSchema:        aws.String("DataSchema"),
			DataSchemaURI:     aws.String("S3Url"),
		},
		RoleARN:           aws.String("RoleARN"), // Required
		ComputeStatistics: aws.Boolean(true),
		DataSourceName:    aws.String("EntityName"),
	}
	resp, err := svc.CreateDataSourceFromRDS(params)

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

func ExampleMachineLearning_CreateDataSourceFromRedshift() {
	svc := machinelearning.New(nil)

	params := &machinelearning.CreateDataSourceFromRedshiftInput{
		DataSourceID: aws.String("EntityId"), // Required
		DataSpec: &machinelearning.RedshiftDataSpec{ // Required
			DatabaseCredentials: &machinelearning.RedshiftDatabaseCredentials{ // Required
				Password: aws.String("RedshiftDatabasePassword"), // Required
				Username: aws.String("RedshiftDatabaseUsername"), // Required
			},
			DatabaseInformation: &machinelearning.RedshiftDatabase{ // Required
				ClusterIdentifier: aws.String("RedshiftClusterIdentifier"), // Required
				DatabaseName:      aws.String("RedshiftDatabaseName"),      // Required
			},
			S3StagingLocation: aws.String("S3Url"),                  // Required
			SelectSQLQuery:    aws.String("RedshiftSelectSqlQuery"), // Required
			DataRearrangement: aws.String("DataRearrangement"),
			DataSchema:        aws.String("DataSchema"),
			DataSchemaURI:     aws.String("S3Url"),
		},
		RoleARN:           aws.String("RoleARN"), // Required
		ComputeStatistics: aws.Boolean(true),
		DataSourceName:    aws.String("EntityName"),
	}
	resp, err := svc.CreateDataSourceFromRedshift(params)

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

func ExampleMachineLearning_CreateDataSourceFromS3() {
	svc := machinelearning.New(nil)

	params := &machinelearning.CreateDataSourceFromS3Input{
		DataSourceID: aws.String("EntityId"), // Required
		DataSpec: &machinelearning.S3DataSpec{ // Required
			DataLocationS3:       aws.String("S3Url"), // Required
			DataRearrangement:    aws.String("DataRearrangement"),
			DataSchema:           aws.String("DataSchema"),
			DataSchemaLocationS3: aws.String("S3Url"),
		},
		ComputeStatistics: aws.Boolean(true),
		DataSourceName:    aws.String("EntityName"),
	}
	resp, err := svc.CreateDataSourceFromS3(params)

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

func ExampleMachineLearning_CreateEvaluation() {
	svc := machinelearning.New(nil)

	params := &machinelearning.CreateEvaluationInput{
		EvaluationDataSourceID: aws.String("EntityId"), // Required
		EvaluationID:           aws.String("EntityId"), // Required
		MLModelID:              aws.String("EntityId"), // Required
		EvaluationName:         aws.String("EntityName"),
	}
	resp, err := svc.CreateEvaluation(params)

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

func ExampleMachineLearning_CreateMLModel() {
	svc := machinelearning.New(nil)

	params := &machinelearning.CreateMLModelInput{
		MLModelID:            aws.String("EntityId"),    // Required
		MLModelType:          aws.String("MLModelType"), // Required
		TrainingDataSourceID: aws.String("EntityId"),    // Required
		MLModelName:          aws.String("EntityName"),
		Parameters: map[string]*string{
			"Key": aws.String("StringType"), // Required
			// More values...
		},
		Recipe:    aws.String("Recipe"),
		RecipeURI: aws.String("S3Url"),
	}
	resp, err := svc.CreateMLModel(params)

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

func ExampleMachineLearning_CreateRealtimeEndpoint() {
	svc := machinelearning.New(nil)

	params := &machinelearning.CreateRealtimeEndpointInput{
		MLModelID: aws.String("EntityId"), // Required
	}
	resp, err := svc.CreateRealtimeEndpoint(params)

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

func ExampleMachineLearning_DeleteBatchPrediction() {
	svc := machinelearning.New(nil)

	params := &machinelearning.DeleteBatchPredictionInput{
		BatchPredictionID: aws.String("EntityId"), // Required
	}
	resp, err := svc.DeleteBatchPrediction(params)

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

func ExampleMachineLearning_DeleteDataSource() {
	svc := machinelearning.New(nil)

	params := &machinelearning.DeleteDataSourceInput{
		DataSourceID: aws.String("EntityId"), // Required
	}
	resp, err := svc.DeleteDataSource(params)

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

func ExampleMachineLearning_DeleteEvaluation() {
	svc := machinelearning.New(nil)

	params := &machinelearning.DeleteEvaluationInput{
		EvaluationID: aws.String("EntityId"), // Required
	}
	resp, err := svc.DeleteEvaluation(params)

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

func ExampleMachineLearning_DeleteMLModel() {
	svc := machinelearning.New(nil)

	params := &machinelearning.DeleteMLModelInput{
		MLModelID: aws.String("EntityId"), // Required
	}
	resp, err := svc.DeleteMLModel(params)

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

func ExampleMachineLearning_DeleteRealtimeEndpoint() {
	svc := machinelearning.New(nil)

	params := &machinelearning.DeleteRealtimeEndpointInput{
		MLModelID: aws.String("EntityId"), // Required
	}
	resp, err := svc.DeleteRealtimeEndpoint(params)

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

func ExampleMachineLearning_DescribeBatchPredictions() {
	svc := machinelearning.New(nil)

	params := &machinelearning.DescribeBatchPredictionsInput{
		EQ:             aws.String("ComparatorValue"),
		FilterVariable: aws.String("BatchPredictionFilterVariable"),
		GE:             aws.String("ComparatorValue"),
		GT:             aws.String("ComparatorValue"),
		LE:             aws.String("ComparatorValue"),
		LT:             aws.String("ComparatorValue"),
		Limit:          aws.Long(1),
		NE:             aws.String("ComparatorValue"),
		NextToken:      aws.String("StringType"),
		Prefix:         aws.String("ComparatorValue"),
		SortOrder:      aws.String("SortOrder"),
	}
	resp, err := svc.DescribeBatchPredictions(params)

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

func ExampleMachineLearning_DescribeDataSources() {
	svc := machinelearning.New(nil)

	params := &machinelearning.DescribeDataSourcesInput{
		EQ:             aws.String("ComparatorValue"),
		FilterVariable: aws.String("DataSourceFilterVariable"),
		GE:             aws.String("ComparatorValue"),
		GT:             aws.String("ComparatorValue"),
		LE:             aws.String("ComparatorValue"),
		LT:             aws.String("ComparatorValue"),
		Limit:          aws.Long(1),
		NE:             aws.String("ComparatorValue"),
		NextToken:      aws.String("StringType"),
		Prefix:         aws.String("ComparatorValue"),
		SortOrder:      aws.String("SortOrder"),
	}
	resp, err := svc.DescribeDataSources(params)

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

func ExampleMachineLearning_DescribeEvaluations() {
	svc := machinelearning.New(nil)

	params := &machinelearning.DescribeEvaluationsInput{
		EQ:             aws.String("ComparatorValue"),
		FilterVariable: aws.String("EvaluationFilterVariable"),
		GE:             aws.String("ComparatorValue"),
		GT:             aws.String("ComparatorValue"),
		LE:             aws.String("ComparatorValue"),
		LT:             aws.String("ComparatorValue"),
		Limit:          aws.Long(1),
		NE:             aws.String("ComparatorValue"),
		NextToken:      aws.String("StringType"),
		Prefix:         aws.String("ComparatorValue"),
		SortOrder:      aws.String("SortOrder"),
	}
	resp, err := svc.DescribeEvaluations(params)

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

func ExampleMachineLearning_DescribeMLModels() {
	svc := machinelearning.New(nil)

	params := &machinelearning.DescribeMLModelsInput{
		EQ:             aws.String("ComparatorValue"),
		FilterVariable: aws.String("MLModelFilterVariable"),
		GE:             aws.String("ComparatorValue"),
		GT:             aws.String("ComparatorValue"),
		LE:             aws.String("ComparatorValue"),
		LT:             aws.String("ComparatorValue"),
		Limit:          aws.Long(1),
		NE:             aws.String("ComparatorValue"),
		NextToken:      aws.String("StringType"),
		Prefix:         aws.String("ComparatorValue"),
		SortOrder:      aws.String("SortOrder"),
	}
	resp, err := svc.DescribeMLModels(params)

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

func ExampleMachineLearning_GetBatchPrediction() {
	svc := machinelearning.New(nil)

	params := &machinelearning.GetBatchPredictionInput{
		BatchPredictionID: aws.String("EntityId"), // Required
	}
	resp, err := svc.GetBatchPrediction(params)

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

func ExampleMachineLearning_GetDataSource() {
	svc := machinelearning.New(nil)

	params := &machinelearning.GetDataSourceInput{
		DataSourceID: aws.String("EntityId"), // Required
		Verbose:      aws.Boolean(true),
	}
	resp, err := svc.GetDataSource(params)

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

func ExampleMachineLearning_GetEvaluation() {
	svc := machinelearning.New(nil)

	params := &machinelearning.GetEvaluationInput{
		EvaluationID: aws.String("EntityId"), // Required
	}
	resp, err := svc.GetEvaluation(params)

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

func ExampleMachineLearning_GetMLModel() {
	svc := machinelearning.New(nil)

	params := &machinelearning.GetMLModelInput{
		MLModelID: aws.String("EntityId"), // Required
		Verbose:   aws.Boolean(true),
	}
	resp, err := svc.GetMLModel(params)

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

func ExampleMachineLearning_Predict() {
	svc := machinelearning.New(nil)

	params := &machinelearning.PredictInput{
		MLModelID:       aws.String("EntityId"), // Required
		PredictEndpoint: aws.String("VipURL"),   // Required
		Record: map[string]*string{ // Required
			"Key": aws.String("VariableValue"), // Required
			// More values...
		},
	}
	resp, err := svc.Predict(params)

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

func ExampleMachineLearning_UpdateBatchPrediction() {
	svc := machinelearning.New(nil)

	params := &machinelearning.UpdateBatchPredictionInput{
		BatchPredictionID:   aws.String("EntityId"),   // Required
		BatchPredictionName: aws.String("EntityName"), // Required
	}
	resp, err := svc.UpdateBatchPrediction(params)

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

func ExampleMachineLearning_UpdateDataSource() {
	svc := machinelearning.New(nil)

	params := &machinelearning.UpdateDataSourceInput{
		DataSourceID:   aws.String("EntityId"),   // Required
		DataSourceName: aws.String("EntityName"), // Required
	}
	resp, err := svc.UpdateDataSource(params)

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

func ExampleMachineLearning_UpdateEvaluation() {
	svc := machinelearning.New(nil)

	params := &machinelearning.UpdateEvaluationInput{
		EvaluationID:   aws.String("EntityId"),   // Required
		EvaluationName: aws.String("EntityName"), // Required
	}
	resp, err := svc.UpdateEvaluation(params)

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

func ExampleMachineLearning_UpdateMLModel() {
	svc := machinelearning.New(nil)

	params := &machinelearning.UpdateMLModelInput{
		MLModelID:      aws.String("EntityId"), // Required
		MLModelName:    aws.String("EntityName"),
		ScoreThreshold: aws.Double(1.0),
	}
	resp, err := svc.UpdateMLModel(params)

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
