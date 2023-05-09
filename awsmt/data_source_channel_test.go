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

func TestAccChannelDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "awsmt_channel" "test"  {
  name = "test"
  channel_state = "STOPPED"
  outputs = [{
    manifest_name                = "default"
    source_group                 = "default"
    hls_playlist_settings = [{manifest_windows_seconds = 30}]
  }]
  playback_mode = "LOOP"
  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:319158032161:channel/test\"}]}"
  tier = "BASIC"
}

data "awsmt_channel" "test" {
  name = awsmt_channel.test.name
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr("data.awsmt_channel.test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "name", "test"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "id", "placeholder"),
				),
			},
		},
	})
}
