// Package aws implements the driftwatch provider interface for Amazon Web Services.
//
// Registration
//
// The package registers itself under the name "aws" via its init function,
// so importing it with a blank identifier is sufficient to make the provider
// available to the provider registry:
//
//	_ "github.com/driftwatch/driftwatch/internal/provider/aws"
//
// Configuration
//
// The following keys are recognised in the provider config map:
//
//	region   – (required) AWS region, e.g. "us-east-1"
//	profile  – (optional) AWS named profile from ~/.aws/credentials
//
// Supported resource types
//
//	ec2.instance  – EC2 virtual machine instance
//	s3.bucket     – S3 object storage bucket
package aws
