package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

func testAccPreCheck(t *testing.T) string {
	if a, b, c := os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_PROFILE"); (a == "" || b == "") && c == "" {
		t.Fatal("Either AWS_PROFILE or both AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
	return fmt.Sprintf(
		`resource "awsmt_source_location" "test_source_location_vod"{
		source_location_name = "source_location_vod"
		http_configuration = {
			hc_base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
		}
	}`)

}

func TestAccChannelResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "awsmt_channel" "resourcetest"  {
  name = "resourcetest"
  channel_state = "RUNNING"
  outputs = [{
    manifest_name                = "default"
    source_group                 = "default"
    hls_playlist_settings = [{manifest_windows_seconds = 30}]
  }]
  playback_mode = "LOOP"
  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:319158032161:channel/resourcetest\"}]}"
  tier = "BASIC"
}
data "awsmt_channel" "resourcetest" {
  name = awsmt_channel.resourcetest.name
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("awsmt_channel.resourcetest", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestMatchResourceAttr("awsmt_channel.resourcetest", "playback_mode", regexp.MustCompile("LOOP")),
					resource.TestMatchResourceAttr("awsmt_channel.resourcetest", "tier", regexp.MustCompile("BASIC")),
					resource.TestCheckResourceAttr("awsmt_channel.resourcetest", "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr("awsmt_channel.resourcetest", "name", "resourcetest"),
					resource.TestCheckResourceAttr("awsmt_channel.resourcetest", "outputs.0.manifest_name", "default"),
				),
			},
		},
	})
}
