package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"flag"
)

// Creates a S3 Bucket in the region configured in the shared config
// or AWS_REGION environment variable.
//
// Usage:
//    go run s3_upload_object.go BUCKET_NAME FILENAME

var aws_region = flag.String("region", "eu-west-3", "AWS Region of bucket")
var aws_profile = flag.String("profile", "test", "AWS Profile") // os.Args[1] //"hktech"
var bucket = flag.String("bucket", "testbucket", "Bucket Name") // os.Args[2]
var filename = flag.String("file", "", "File to Upload") //os.Args[3]
var acl = flag.String("acl", "public-read", "ACL on Bucket") 

func main() {

	flag.Parse()
	file, err := os.Open(*filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", filename, err)
	}

	defer file.Close()

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: *aws_profile,
		Config: aws.Config{
			Region: aws.String(*aws_region),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})

	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*filename),
		Body:   file,
		ACL:    aws.String(*acl),
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", *filename, *bucket, err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", *filename, *bucket)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
