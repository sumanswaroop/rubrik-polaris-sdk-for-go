package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kr/pretty"
	"golang.org/x/sync/errgroup"

	"github.com/rubrikinc/rubrik-polaris-sdk-for-go/internal/testsetup"
	"github.com/rubrikinc/rubrik-polaris-sdk-for-go/pkg/polaris"
	"github.com/rubrikinc/rubrik-polaris-sdk-for-go/pkg/polaris/aws"
	"github.com/rubrikinc/rubrik-polaris-sdk-for-go/pkg/polaris/azure"
	"github.com/rubrikinc/rubrik-polaris-sdk-for-go/pkg/polaris/gcp"
	"github.com/rubrikinc/rubrik-polaris-sdk-for-go/pkg/polaris/graphql"
	"github.com/rubrikinc/rubrik-polaris-sdk-for-go/pkg/polaris/graphql/core"
	polaris_log "github.com/rubrikinc/rubrik-polaris-sdk-for-go/pkg/polaris/log"
)

func main() {
	cleanup := flag.Bool("cleanup", false, "Perform cleanup tasks post CI")
	precheck := flag.Bool("precheck", false, "Check pre-requirements for running CI tests")
	flag.Parse()
	if *cleanup == *precheck {
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	// Load configuration and create client
	polAccount, err := polaris.DefaultServiceAccount(true)
	if err != nil {
		log.Fatal(err)
	}
	logger := polaris_log.NewStandardLogger()
	logger.SetLogLevel(polaris_log.Info)
	client, err := polaris.NewClient(ctx, polAccount, logger)
	if err != nil {
		log.Fatal(err)
	}

	if *precheck {
		err = check(ctx, client)
	} else {
		err = clean(ctx, client)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func check(ctx context.Context, client *polaris.Client) error {
	var g errgroup.Group

	// AWS
	g.Go(func() error {
		testAcc, err := testsetup.AWSAccount()
		if err != nil {
			return err
		}
		awsAccount, err := client.AWS().Account(ctx, aws.AccountID(testAcc.AccountID), core.FeatureAll)
		switch {
		case err == nil:
			return fmt.Errorf("found pre-existing AWS account: %v", pretty.Sprint(awsAccount))
		case !errors.Is(err, graphql.ErrNotFound):
			return fmt.Errorf("failed to check AWS account: %v", err)
		}
		return nil
	})

	// Azure
	g.Go(func() error {
		testSub, err := testsetup.AzureSubscription()
		if err != nil {
			return err
		}
		azureAcc, err := client.Azure().Subscription(ctx, azure.SubscriptionID(testSub.SubscriptionID), core.FeatureAll)
		switch {
		case err == nil:
			return fmt.Errorf("found pre-existing Azure subscription: %v", pretty.Sprint(azureAcc))
		case !errors.Is(err, graphql.ErrNotFound):
			return fmt.Errorf("failed to check Azure account: %v", err)
		}
		return nil
	})

	// GCP
	g.Go(func() error {
		testProj, err := testsetup.GCPProject()
		if err != nil {
			return err
		}
		proj, err := client.GCP().Project(ctx, gcp.ProjectID(testProj.ProjectID), core.FeatureAll)
		switch {
		case err == nil:
			return fmt.Errorf("found pre-existing GCP projects: %v", pretty.Sprint(proj))
		case !errors.Is(err, graphql.ErrNotFound):
			return fmt.Errorf("failed to check GCP project: %v", err)
		}
		return nil
	})

	return g.Wait()
}

func clean(ctx context.Context, client *polaris.Client) error {
	var g errgroup.Group

	// AWS
	g.Go(func() error {
		testAcc, err := testsetup.AWSAccount()
		if err != nil {
			return err
		}
		awsAccount, err := client.AWS().Account(ctx, aws.AccountID(testAcc.AccountID), core.FeatureAll)
		switch {
		case errors.Is(err, graphql.ErrNotFound):
			return nil
		case err != nil:
			return fmt.Errorf("failed to check AWS account: %v", err)
		}
		if awsAccount.NativeID != testAcc.AccountID {
			return fmt.Errorf("existing AWS account %q isn't expected test account %q, won't remove",
				awsAccount.NativeID, testAcc.AccountID)
		}
		// TODO: we might need to iterate over awsAccount.Features to remove
		// all of them in the future
		return client.AWS().RemoveAccount(ctx, aws.Profile(testAcc.Profile), core.FeatureCloudNativeProtection, false)
	})

	// Azure
	g.Go(func() error {
		testSub, err := testsetup.AzureSubscription()
		if err != nil {
			return err
		}
		azureAcc, err := client.Azure().Subscription(ctx, azure.SubscriptionID(testSub.SubscriptionID), core.FeatureAll)
		switch {
		case errors.Is(err, graphql.ErrNotFound):
			return nil
		case err != nil:
			return fmt.Errorf("failed to check Azure subscription: %v", err)
		}
		if azureAcc.NativeID != testSub.SubscriptionID {
			return fmt.Errorf("existing Azure subscription %q isn't the expected test subscription %q, won't remove",
				azureAcc.NativeID, testSub.SubscriptionID)
		}
		// Azure doesn't automatically remove exocompute configs when removing
		// the subscription so we need to do it manually here
		exoCfgs, err := client.Azure().ExocomputeConfigs(ctx, azure.SubscriptionID(testSub.SubscriptionID))
		if err != nil {
			return err
		}
		for i := range exoCfgs {
			if err := client.Azure().RemoveExocomputeConfig(ctx, exoCfgs[i].ID); err != nil {
				return fmt.Errorf("failed to remove Azure ExocomputeConfig: %v", pretty.Sprint(exoCfgs[i]))
			}
		}

		sub := azure.Subscription(testSub.SubscriptionID, testSub.TenantDomain)
		// TODO: we might need to iterate over azureAcc.Features to remove
		// all of them in the future
		return client.Azure().RemoveSubscription(ctx, azure.ID(sub), core.FeatureCloudNativeProtection, false)
	})

	// GCP
	g.Go(func() error {
		testProj, err := testsetup.GCPProject()
		if err != nil {
			return err
		}
		proj, err := client.GCP().Project(ctx, gcp.ProjectID(testProj.ProjectID), core.FeatureAll)
		switch {
		case errors.Is(err, graphql.ErrNotFound):
			return nil
		case err != nil:
			return fmt.Errorf("failed to check GCP project: %v", err)
		}
		if pn := proj.ProjectNumber; pn != testProj.ProjectNumber {
			return fmt.Errorf("existing GCP project %q isn't expected test project %q, won't remove",
				pn, testProj.ProjectNumber)
		}
		return client.GCP().RemoveProject(ctx, gcp.ProjectNumber(testProj.ProjectNumber), core.FeatureCloudNativeProtection, false)
	})

	return g.Wait()
}
