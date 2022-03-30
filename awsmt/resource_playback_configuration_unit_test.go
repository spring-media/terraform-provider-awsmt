package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"testing"
)

var sess, _ = session.NewSession(&aws.Config{Region: aws.String("eu-central-11")})
var c = mediatailor.New(sess)

func TestGetPlaybackConfigurationError(t *testing.T) {
	v, err := getSinglePlaybackConfiguration(c, "not-a-configuration")

	if err == nil {
		t.Fatalf("expected error, got: %v", v)
	}
}
