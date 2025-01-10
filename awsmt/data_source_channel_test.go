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
				Config: basicChannelDSHLS(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "id", "test"),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "name", "test"),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "channel_state", "STOPPED"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "playback_mode", "LOOP"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "tier", "BASIC"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "tags.Environment", "dev"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.manifest_name", "default"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.source_group", "default"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.manifest_window_seconds", "30"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "enable_as_run_logs", "false"),
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
				Config: basicChannelDSHLSWithSlate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "id", "test"),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "name", "test"),
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
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "enable_as_run_logs", "false"),
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
				Config:      channelErrorDS(),
				ExpectError: regexp.MustCompile("Error while describing channel "),
			},
		},
	})
}

func basicChannelDSHLS() string {
	return `
				resource "awsmt_channel" "test"  {
  					name = "test"
  					channel_state = "STOPPED"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							ad_markup_type = ["DATERANGE"]
							manifest_window_seconds = 30
						}
  					}]
  					playback_mode = "LOOP"
  					policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:985600762523:channel/test\"}]}"
  					tier = "BASIC"
					tags = {"Environment": "dev"}
				}

				data "awsmt_channel" "test" {
  					name = awsmt_channel.test.name
				}
				output "channel_out" {
					value = data.awsmt_channel.test
				}
				`
}

func basicChannelDSHLSWithSlate() string {
	return `
				resource "awsmt_vod_source" "test" {
  					http_package_configurations = [{
						path = "/"
						source_group = "default"
    					type = "HLS"
  					}]
  					source_location_name = awsmt_source_location.test_source_location.name
  					name = "vod_source_example"
					tags = {"Environment": "dev"}
				}
				data "awsmt_vod_source" "data_test" {
  					source_location_name = awsmt_source_location.test_source_location.name
  					name = awsmt_vod_source.test.name
				}

				output "vod_source_out" {
  					value = data.awsmt_vod_source.data_test
				}
				resource "awsmt_source_location" "test_source_location"{
  					name = "test_source_location"
  					http_configuration = {
    					base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
  					}
  					default_segment_delivery_configuration = {
    					base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  					}
				}
				data "awsmt_source_location" "test" {
  					name = awsmt_source_location.test_source_location.name
				}
				output "awsmt_source_location" {
  					value = data.awsmt_source_location.test
				}
				resource "awsmt_channel" "test"  {
  					name = "test"
  					channel_state = "STOPPED"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							ad_markup_type = ["DATERANGE"]
							manifest_window_seconds = 30
						}
  					}]
  					playback_mode = "LINEAR"
					filler_slate = {
						source_location_name = awsmt_source_location.test_source_location.name
						vod_source_name = awsmt_vod_source.test.name
					}
  					policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:985600762523:channel/test\"}]}"
  					tier = "BASIC"
					tags = {"Environment": "dev"}
				}

				data "awsmt_channel" "test" {
  					name = awsmt_channel.test.name
				}
				output "channel_out" {
					value = data.awsmt_channel.test
				}
				`
}

func channelErrorDS() string {
	return `
				resource "awsmt_channel" "test"  {
  					name = "test"
  					channel_state = "STOPPED"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							ad_markup_type = ["DATERANGE"]
							manifest_window_seconds = 30
						}
  					}]
  					playback_mode = "LOOP"
  					policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:985600762523:channel/test\"}]}"
  					tier = "BASIC"
					tags = {"Environment": "dev"}
				}

				data "awsmt_channel" "test" {
  					name = "testingError"
				}
				output "channel_out" {
					value = data.awsmt_channel.test
				}
				`
}
