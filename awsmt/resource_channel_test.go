package awsmt

import (
	"fmt"
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
	name := "test"
	state_stopped := "STOPPED"
	state_running := "RUNNING"
	mw_s := "30"
	mw_s2 := "40"
	mbt_s := "2"
	mbt_s2 := "3"
	mup_s := "2"
	mup_s2 := "3"
	spd_s := "2"
	spd_s2 := "3"
	k1 := "Environment"
	v1 := "dev"
	k2 := "Testing"
	v2 := "pass"
	k3 := "Environment"
	v3 := "prod"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: basicChannel(name, state_stopped, mw_s, mbt_s, mup_s, spd_s, k1, v1, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "id", "test"),
					resource.TestMatchResourceAttr("awsmt_channel.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr("awsmt_channel.test", "name", "test"),
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
				Config: basicChannel(name, state_running, mw_s2, mbt_s2, mup_s2, spd_s2, k3, v3, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "id", "test"),
					resource.TestMatchResourceAttr("awsmt_channel.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr("awsmt_channel.test", "name", "test"),
					resource.TestMatchResourceAttr("awsmt_channel.test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr("awsmt_channel.test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "playback_mode", "LOOP"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "tier", "BASIC"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "tags.Environment", "prod"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "tags.Testing", "pass"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.manifest_name", "default"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.source_group", "default"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.manifest_window_seconds", "40"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.min_buffer_time_seconds", "3"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.min_update_period_seconds", "3"),
					resource.TestCheckResourceAttr("awsmt_channel.test", "outputs.0.dash_playlist_settings.suggested_presentation_delay_seconds", "3"),
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
				Config:      errorChannel(),
				ExpectError: regexp.MustCompile("Error while creating channel "),
			},
		},
	})
}

func TestAccChannelResourceRunning(t *testing.T) {
	mw_s := "30"
	mw_s2 := "40"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: hlsChannel(mw_s),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.manifest_window_seconds", "30"),
				),
			},
			{
				Config: hlsChannel(mw_s2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.manifest_window_seconds", "40"),
				),
			},
		},
	})
}

func TestAccChannelResourceSTANDARD(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: standardTierChannel(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "tier", "STANDARD"),
				),
			},
		},
	})
}

func basicChannel(name, state, mw_s, mbt_s, mup_s, spd_s, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(
		`
				resource "awsmt_channel" "test"  {
  					name = "%[1]s"
  					channel_state = "%[2]s"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					dash_playlist_settings = {
							manifest_window_seconds = "%[3]s"
							min_buffer_time_seconds = "%[4]s"
							min_update_period_seconds = "%[5]s"
							suggested_presentation_delay_seconds = "%[6]s"
						}
  					}]
  					playback_mode = "LOOP"
  					policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:319158032161:channel/test\"}]}"
  					tier = "BASIC"
					tags = {
   		 						"%[7]s": "%[8]s",
								"%[9]s": "%[10]s"
							}
					}

				data "awsmt_channel" "test" {
  					name = awsmt_channel.test.name
				}
				output "channel_out" {
					value = data.awsmt_channel.test
				}
				`, name, state, mw_s, mbt_s, mup_s, spd_s, k1, v1, k2, v2,
	)
}

func errorChannel() string {
	return `
				resource "awsmt_channel" "test"  {
  					name = "test"
  					channel_state = "RUNNING"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							ad_markup_type = ["DATERANGE"]
							manifest_window_seconds = 30
						}
  					}]
  					playback_mode = "LINEAR"
  					policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:319158032161:channel/test\"}]}"
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

func hlsChannel(mw_s string) string {
	return fmt.Sprintf(`
				resource "awsmt_channel" "test"  {
  					name = "test"
  					channel_state = "RUNNING"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							ad_markup_type = ["DATERANGE"]
							manifest_window_seconds = "%[1]s"
						}
  					}]
  					playback_mode = "LOOP"
  					tier = "BASIC"
					tags = {"Environment": "dev"}
					}

				data "awsmt_channel" "test" {
  					name = awsmt_channel.test.name
				}
				output "channel_out" {
					value = data.awsmt_channel.test
				}
				`, mw_s)
}
func standardTierChannel() string {
	return `resource "awsmt_vod_source" "test" {
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
policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:319158032161:channel/test\"}]}"
tier = "STANDARD"
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
