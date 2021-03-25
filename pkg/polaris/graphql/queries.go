// Code generated by queries_gen.go DO NOT EDIT

// MIT License
//
// Copyright (c) 2021 Rubrik
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package graphql

// awsCloudAccountDeleteInitiate GraphQL query
var awsCloudAccountDeleteInitiateQuery = `mutation SdkGolang_awsCloudAccountDeleteInitiate($polarisAccountId: UUID!) {
    awsCloudAccountDeleteInitiate(cloudAccountUuid: $polarisAccountId, awsCloudAccountDeleteInitiateArg: {feature: CLOUD_NATIVE_PROTECTION}) {
        cloudFormationUrl
    }
}`

// awsCloudAccountDeleteProcess GraphQL query
var awsCloudAccountDeleteProcessQuery = `mutation SdkGolang_awsCloudAccountDeleteProcess($polarisAccountId: UUID!) {
    awsCloudAccountDeleteProcess(cloudAccountUuid: $polarisAccountId, awsCloudAccountDeleteProcessArg: {feature: CLOUD_NATIVE_PROTECTION}) {
        message
    }
}`

// awsCloudAccountSave GraphQL query
var awsCloudAccountSaveQuery = `mutation SdkGolang_awsCloudAccountSave($polarisAccountId: UUID!, $awsRegions: [AwsCloudAccountRegionEnum!]!) {
  awsCloudAccountSave(cloudAccountUuid: $polarisAccountId, awsCloudAccountSaveArg: {action: UPDATE_REGIONS, feature: CLOUD_NATIVE_PROTECTION, awsRegions: $awsRegions}) {
    message
  }
}`

// awsCloudAccountUpdateFeatureInitiate GraphQL query
var awsCloudAccountUpdateFeatureInitiateQuery = `mutation SdkGolang_awsCloudAccountUpdateFeatureInitiate($polarisAccountId: UUID!) {
    awsCloudAccountUpdateFeatureInitiate(cloudAccountUuid: $polarisAccountId, features: [CLOUD_NATIVE_PROTECTION]) {
        cloudFormationUrl
        templateUrl
    }
}`

// awsCloudAccounts GraphQL query
var awsCloudAccountsQuery = `query SdkGolang_awsCloudAccounts($columnFilter: String = "") {
    awsCloudAccounts(awsCloudAccountsArg: {columnSearchFilter: $columnFilter, statusFilters: [], feature: CLOUD_NATIVE_PROTECTION}) {
        awsCloudAccounts {
            awsCloudAccount {
                id
                nativeId
                message
                accountName
            }
            featureDetails {
                feature
                roleArn
                stackArn
                status
                awsRegions
            }
        }
    }
}`

// awsDeleteNativeAccount GraphQL query
var awsDeleteNativeAccountQuery = `mutation SdkGolang_awsDeleteNativeAccount($polarisAccountId: UUID!, $deleteNativeSnapshots: Boolean = false, $awsNativeProtectionFeature: AwsNativeProtectionFeatureEnum = EC2) {
    deleteAwsNativeAccount(awsNativeAccountId: $polarisAccountId, deleteNativeSnapshots: $deleteNativeSnapshots, awsNativeProtectionFeature: $awsNativeProtectionFeature) {
        taskchainUuid
    }
}`

// awsNativeAccountConnection GraphQL query
var awsNativeAccountConnectionQuery = `query SdkGolang_awsNativeAccountConnection($awsNativeProtectionFeature: AwsNativeProtectionFeatureEnum = EC2, $nameFilter: String = "") {
	awsNativeAccountConnection(awsNativeProtectionFeature: $awsNativeProtectionFeature, accountFilters: {nameSubstringFilter: {nameSubstring: $nameFilter}}) {
		count
		edges {
			node {
				id
				regions
				status
				name
				slaAssignment
				configuredSlaDomain {
					id
					name
				}
				effectiveSlaDomain {
					id
					name
				}
			}
		}
		pageInfo {
			endCursor
			hasNextPage
		}
	}
}`

// awsNativeProtectionAccountAdd GraphQL query
var awsNativeProtectionAccountAddQuery = `mutation SdkGolang_awsNativeProtectionAccountAdd($accountId: String!, $accountName: String!, $regions: [String!]!) {
    awsNativeProtectionAccountAdd(awsNativeProtectionAccountAddArg: {accountId: $accountId, name: $accountName, regions: $regions}) {
       cloudFormationName
       cloudFormationUrl
       cloudFormationTemplateUrl
       errorMessage
    }
}`

// coreTaskchainStatus GraphQL query
var coreTaskchainStatusQuery = `query SdkGolang_coreTaskchainStatus($taskchainId: String!){
    getKorgTaskchainStatus(taskchainId: $taskchainId){
        taskchain {
            id
            state
            taskchainUuid
            ... on Taskchain{
                progressedAt
            }
        }
    }
}`
