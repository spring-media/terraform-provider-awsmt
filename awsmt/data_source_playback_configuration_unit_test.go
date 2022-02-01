package awsmt

import (
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"reflect"
	"testing"
)

func TestFlattenPlaybackConfiguration(t *testing.T) {
	// arrange
	testString := "testString"
	var input = mediatailor.PlaybackConfiguration{
		AdDecisionServerUrl:                 &testString,
		CdnConfiguration:                    &mediatailor.CdnConfiguration{AdSegmentUrlPrefix: &testString, ContentSegmentUrlPrefix: &testString},
		DashConfiguration:                   &mediatailor.DashConfiguration{ManifestEndpointPrefix: &testString, MpdLocation: &testString, OriginManifestType: &testString},
		HlsConfiguration:                    &mediatailor.HlsConfiguration{ManifestEndpointPrefix: &testString},
		Name:                                &testString,
		PlaybackConfigurationArn:            &testString,
		PlaybackEndpointPrefix:              &testString,
		SessionInitializationEndpointPrefix: &testString,
		SlateAdUrl:                          &testString,
		Tags:                                map[string]*string{},
		TranscodeProfileName:                &testString,
		VideoContentSourceUrl:               &testString,
	}
	var expected = map[string]interface{}{
		"ad_decision_server_url": &testString,
		"cdn_configuration": []interface{}{map[string]interface{}{
			"ad_segment_url_prefix":      &testString,
			"content_segment_url_prefix": &testString,
		}},
		"dash_configuration": []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": &testString,
			"mpd_location":             &testString,
			"origin_manifest_type":     &testString,
		}},
		"hls_configuration": []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": &testString,
		}},
		"name":                                   &testString,
		"playback_configuration_arn":             &testString,
		"playback_endpoint_prefix":               &testString,
		"session_initialization_endpoint_prefix": &testString,
		"slate_ad_url":                           &testString,
		"tags":                                   map[string]*string{},
		"transcode_profile_name":                 &testString,
		"video_content_source_url":               &testString,
	}
	// act
	output := flatten(&input)
	// assert
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Not matching. Expected:\n%#v\nGot\n%#v", expected, output)
	}
}
