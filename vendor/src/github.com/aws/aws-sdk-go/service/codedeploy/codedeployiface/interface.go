// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

// Package codedeployiface provides an interface for the AWS CodeDeploy.
package codedeployiface

import (
	"github.com/aws/aws-sdk-go/service/codedeploy"
)

// CodeDeployAPI is the interface type for codedeploy.CodeDeploy.
type CodeDeployAPI interface {
	AddTagsToOnPremisesInstances(*codedeploy.AddTagsToOnPremisesInstancesInput) (*codedeploy.AddTagsToOnPremisesInstancesOutput, error)

	BatchGetApplications(*codedeploy.BatchGetApplicationsInput) (*codedeploy.BatchGetApplicationsOutput, error)

	BatchGetDeployments(*codedeploy.BatchGetDeploymentsInput) (*codedeploy.BatchGetDeploymentsOutput, error)

	BatchGetOnPremisesInstances(*codedeploy.BatchGetOnPremisesInstancesInput) (*codedeploy.BatchGetOnPremisesInstancesOutput, error)

	CreateApplication(*codedeploy.CreateApplicationInput) (*codedeploy.CreateApplicationOutput, error)

	CreateDeployment(*codedeploy.CreateDeploymentInput) (*codedeploy.CreateDeploymentOutput, error)

	CreateDeploymentConfig(*codedeploy.CreateDeploymentConfigInput) (*codedeploy.CreateDeploymentConfigOutput, error)

	CreateDeploymentGroup(*codedeploy.CreateDeploymentGroupInput) (*codedeploy.CreateDeploymentGroupOutput, error)

	DeleteApplication(*codedeploy.DeleteApplicationInput) (*codedeploy.DeleteApplicationOutput, error)

	DeleteDeploymentConfig(*codedeploy.DeleteDeploymentConfigInput) (*codedeploy.DeleteDeploymentConfigOutput, error)

	DeleteDeploymentGroup(*codedeploy.DeleteDeploymentGroupInput) (*codedeploy.DeleteDeploymentGroupOutput, error)

	DeregisterOnPremisesInstance(*codedeploy.DeregisterOnPremisesInstanceInput) (*codedeploy.DeregisterOnPremisesInstanceOutput, error)

	GetApplication(*codedeploy.GetApplicationInput) (*codedeploy.GetApplicationOutput, error)

	GetApplicationRevision(*codedeploy.GetApplicationRevisionInput) (*codedeploy.GetApplicationRevisionOutput, error)

	GetDeployment(*codedeploy.GetDeploymentInput) (*codedeploy.GetDeploymentOutput, error)

	GetDeploymentConfig(*codedeploy.GetDeploymentConfigInput) (*codedeploy.GetDeploymentConfigOutput, error)

	GetDeploymentGroup(*codedeploy.GetDeploymentGroupInput) (*codedeploy.GetDeploymentGroupOutput, error)

	GetDeploymentInstance(*codedeploy.GetDeploymentInstanceInput) (*codedeploy.GetDeploymentInstanceOutput, error)

	GetOnPremisesInstance(*codedeploy.GetOnPremisesInstanceInput) (*codedeploy.GetOnPremisesInstanceOutput, error)

	ListApplicationRevisions(*codedeploy.ListApplicationRevisionsInput) (*codedeploy.ListApplicationRevisionsOutput, error)

	ListApplications(*codedeploy.ListApplicationsInput) (*codedeploy.ListApplicationsOutput, error)

	ListDeploymentConfigs(*codedeploy.ListDeploymentConfigsInput) (*codedeploy.ListDeploymentConfigsOutput, error)

	ListDeploymentGroups(*codedeploy.ListDeploymentGroupsInput) (*codedeploy.ListDeploymentGroupsOutput, error)

	ListDeploymentInstances(*codedeploy.ListDeploymentInstancesInput) (*codedeploy.ListDeploymentInstancesOutput, error)

	ListDeployments(*codedeploy.ListDeploymentsInput) (*codedeploy.ListDeploymentsOutput, error)

	ListOnPremisesInstances(*codedeploy.ListOnPremisesInstancesInput) (*codedeploy.ListOnPremisesInstancesOutput, error)

	RegisterApplicationRevision(*codedeploy.RegisterApplicationRevisionInput) (*codedeploy.RegisterApplicationRevisionOutput, error)

	RegisterOnPremisesInstance(*codedeploy.RegisterOnPremisesInstanceInput) (*codedeploy.RegisterOnPremisesInstanceOutput, error)

	RemoveTagsFromOnPremisesInstances(*codedeploy.RemoveTagsFromOnPremisesInstancesInput) (*codedeploy.RemoveTagsFromOnPremisesInstancesOutput, error)

	StopDeployment(*codedeploy.StopDeploymentInput) (*codedeploy.StopDeploymentOutput, error)

	UpdateApplication(*codedeploy.UpdateApplicationInput) (*codedeploy.UpdateApplicationOutput, error)

	UpdateDeploymentGroup(*codedeploy.UpdateDeploymentGroupInput) (*codedeploy.UpdateDeploymentGroupOutput, error)
}
