package awsmt

import (
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"reflect"
	"testing"
)

func TestFlattenPlaybackConfiguration(t *testing.T) {
	// arrange
	testString := "testString"
	var testNumber int64 = 10
	var testBool = true
	var input = mediatailor.PlaybackConfiguration{
		AdDecisionServerUrl:                 &testString,
		AvailSuppression:                    &mediatailor.AvailSuppression{Mode: &testString, Value: &testString},
		Bumper:                              &mediatailor.Bumper{EndUrl: &testString, StartUrl: &testString},
		CdnConfiguration:                    &mediatailor.CdnConfiguration{AdSegmentUrlPrefix: &testString, ContentSegmentUrlPrefix: &testString},
		ConfigurationAliases:                map[string]map[string]*string{},
		DashConfiguration:                   &mediatailor.DashConfiguration{ManifestEndpointPrefix: &testString, MpdLocation: &testString, OriginManifestType: &testString},
		HlsConfiguration:                    &mediatailor.HlsConfiguration{ManifestEndpointPrefix: &testString},
		LivePreRollConfiguration:            &mediatailor.LivePreRollConfiguration{AdDecisionServerUrl: &testString, MaxDurationSeconds: &testNumber},
		LogConfiguration:                    &mediatailor.LogConfiguration{PercentEnabled: &testNumber},
		ManifestProcessingRules:             &mediatailor.ManifestProcessingRules{AdMarkerPassthrough: &mediatailor.AdMarkerPassthrough{Enabled: &testBool}},
		Name:                                &testString,
		PersonalizationThresholdSeconds:     &testNumber,
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
		"avail_suppression": []interface{}{map[string]interface{}{
			"mode":  &testString,
			"value": &testString,
		}},
		"bumper": []interface{}{map[string]interface{}{
			"end_url":   &testString,
			"start_url": &testString,
		}},
		"cdn_configuration": []interface{}{map[string]interface{}{
			"ad_segment_url_prefix":      &testString,
			"content_segment_url_prefix": &testString,
		}},
		"configuration_aliases": map[string]map[string]*string{},
		"dash_configuration": []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": &testString,
			"mpd_location":             &testString,
			"origin_manifest_type":     &testString,
		}},
		"hls_configuration": []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": &testString,
		}},
		"live_pre_roll_configuration": []interface{}{map[string]interface{}{
			"ad_decision_server_url": &testString,
			"max_duration_seconds":   &testNumber,
		}},
		"log_configuration": []interface{}{map[string]interface{}{
			"percent_enabled": &testNumber,
		}},
		"manifest_processing_rules": []interface{}{map[string]interface{}{
			"ad_marker_passthrough": []interface{}{map[string]interface{}{
				"enabled": &testBool,
			}},
		}},
		"name":                                   &testString,
		"personalization_threshold_seconds":      &testNumber,
		"playback_configuration_arn":             &testString,
		"playback_endpoint_prefix":               &testString,
		"session_initialization_endpoint_prefix": &testString,
		"slate_ad_url":                           &testString,
		"transcode_profile_name":                 &testString,
		"video_content_source_url":               &testString,
	}
	// act
	output := flattenPlaybackConfiguration(&input)
	// assert
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Not matching. Expected:\n%#v\nGot\n%#v", expected, output)
	}
}

func TestFlattenPlaybackConfigurationNil(t *testing.T) {
	expected := map[string]interface{}{}
	output := flattenPlaybackConfiguration(nil)
	if !reflect.DeepEqual(expected, output) {
		t.Fatalf("Not matching. Expected:\n%#v\nGot\n%#v", expected, output)
	}
}
