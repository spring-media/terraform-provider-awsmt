package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

func testAccPreCheck(t *testing.T) {
	if a, b, c := os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_PROFILE"); (a == "" || b == "") && c == "" {
		t.Fatal("Either AWS_PROFILE or both AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
}

func TestAccChannelResourceBasic(t *testing.T) {
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
    					dash_playlist_settings = {
							manifest_window_seconds = 30
							min_buffer_time_seconds = 2
							min_update_period_seconds = 2
							suggested_presentation_delay_seconds = 2
						}
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
					resource.TestCheckResourceAttr("awsmt_channel.test", "id", "test"),
					resource.TestMatchResourceAttr("awsmt_channel.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_name", "test"),
					resource.TestMatchResourceAttr("awsmt_channel.test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr("awsmt_channel.test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "STOPPED"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "playback_mode", "LOOP"),
					resource.TestMatchResourceAttr("awsmt_channel.test", "policy", regexp.MustCompile(`mediatailor:GetManifest`)),
					resource.TestCheckResourceAttr("awsmt_channel.test", "tier", "BASIC"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "tags.Environment", "dev"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.manifest_name", "default"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.source_group", "default"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.manifest_window_seconds", "30"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.min_buffer_time_seconds", "2"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.min_update_period_seconds", "2"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.suggested_presentation_delay_seconds", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName: "awsmt_channel.test",
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: `
				resource "awsmt_channel" "test"  {
  					channel_name = "test"
  					channel_state = "RUNNING"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					dash_playlist_settings = {
							manifest_window_seconds = 30
							min_buffer_time_seconds = 2
							min_update_period_seconds = 2
							suggested_presentation_delay_seconds = 2
						}
  					}]
  					playback_mode = "LOOP"
					tier = "BASIC"
					tags = {"Environment": "dev", "Name": "test"}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "id", "test"),
					resource.TestMatchResourceAttr("awsmt_channel.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_name", "test"),
					resource.TestMatchResourceAttr("awsmt_channel.test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr("awsmt_channel.test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "playback_mode", "LOOP"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "tier", "BASIC"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "tags.Environment", "dev"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "tags.Name", "test"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.manifest_name", "default"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.source_group", "default"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.manifest_window_seconds", "30"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.min_buffer_time_seconds", "2"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.min_update_period_seconds", "2"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.suggested_presentation_delay_seconds", "2"),
				),
			},
		},
	})
}

func TestAccChannelResourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "awsmt_channel" "test"  {
  					channel_name = "test"
  					channel_state = "RUNNING"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							manifest_window_seconds = 40
						}
  					}]
  					playback_mode = "LINEAR"
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
				ExpectError: regexp.MustCompile("Error while creating channel "),
			},
		},
	})
}

func TestAccChannelResourceRunning(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "awsmt_channel" "test"  {
  					channel_name = "test"
  					channel_state = "RUNNING"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							manifest_window_seconds = 40
						}
  					}]
  					playback_mode = "LOOP"
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
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.manifest_window_seconds", "40"),
				),
			},
			{
				Config: `
				resource "awsmt_channel" "test"  {
  					channel_name = "test"
  					channel_state = "RUNNING"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							manifest_window_seconds = 30
						}
  					}]
  					playback_mode = "LOOP"
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
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.manifest_window_seconds", "30"),
				),
			},
		},
	})
}
