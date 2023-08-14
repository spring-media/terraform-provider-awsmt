package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccChannelDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "awsmt_channel" "test"  {
  					channel_name = "test"
  					channel_state = "STOPPED"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {manifest_window_seconds = 30}
  					}]
  					playback_mode = "LOOP"
  					policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:319158032161:channel/test\"}]}"
  					tier = "BASIC"
					tags = {"Environment": "dev"}
					}

				data "awsmt_channel" "test" {
  					channel_name = awsmt_channel.test.channel_name
				}
				output "channel_out" {
					value = data.awsmt_channel.test
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "id", "test"),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "channel_name", "test"),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "channel_state", "STOPPED"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "playback_mode", "LOOP"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "tier", "BASIC"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "tags.Environment", "dev"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.manifest_name", "default"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.source_group", "default"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.manifest_window_seconds", "30"),
				),
			},
		},
	})
}

func TestAccChannelDataSourceFillerSlateLinear(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "awsmt_vod_source" "test" {
  					http_package_configurations = [{
						path = "/"
						source_group = "default"
    					type = "HLS"
  					}]
  					source_location_name = awsmt_source_location.test_source_location.source_location_name
  					vod_source_name = "vod_source_example"
					tags = {"Environment": "dev"}
				}
				data "awsmt_vod_source" "data_test" {
  					source_location_name = awsmt_source_location.test_source_location.source_location_name
  					vod_source_name = awsmt_vod_source.test.vod_source_name
				}

				output "vod_source_out" {
  					value = data.awsmt_vod_source.data_test
				}
				resource "awsmt_source_location" "test_source_location"{
  					source_location_name = "test_source_location"
  					http_configuration = {
    					hc_base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
  					}
  					default_segment_delivery_configuration = {
    					dsdc_base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  					}
				}
				data "awsmt_source_location" "test" {
  					source_location_name = awsmt_source_location.test_source_location.source_location_name
				}
				output "awsmt_source_location" {
  					value = data.awsmt_source_location.test
				}
				resource "awsmt_channel" "test"  {
  					channel_name = "test"
  					channel_state = "STOPPED"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {manifest_window_seconds = 30}
  					}]
  					playback_mode = "LINEAR"
					filler_slate = {
						source_location_name = awsmt_source_location.test_source_location.source_location_name
						vod_source_name = awsmt_vod_source.test.vod_source_name
					}
  					policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:319158032161:channel/test\"}]}"
  					tier = "BASIC"
					tags = {"Environment": "dev"}
					}

				data "awsmt_channel" "test" {
  					channel_name = awsmt_channel.test.channel_name
				}
				output "channel_out" {
					value = data.awsmt_channel.test
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "id", "test"),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "channel_name", "test"),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "channel_state", "STOPPED"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "playback_mode", "LINEAR"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "filler_slate.source_location_name", "test_source_location"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "filler_slate.vod_source_name", "vod_source_example"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "tier", "BASIC"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "tags.Environment", "dev"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.manifest_name", "default"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.source_group", "default"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.manifest_window_seconds", "30"),
				),
			},
		},
	})
}

func TestAccChannelDataSourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "awsmt_channel" "test"  {
  					channel_name = "test"
  					channel_state = "STOPPED"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {manifest_window_seconds = 30}
  					}]
  					playback_mode = "LOOP"
  					policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:319158032161:channel/test\"}]}"
  					tier = "BASIC"
					tags = {"Environment": "dev"}
					}

				data "awsmt_channel" "test" {
  					channel_name = "testingError"
				}
				output "channel_out" {
					value = data.awsmt_channel.test
				}
				`,
				ExpectError: regexp.MustCompile("Error while describing channel "),
			},
		},
	})
}
