package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func init() {
	resource.AddTestSweepers("test_playback_configuration", &resource.Sweeper{
		Name: "test_playback_configuration",
		F: func(region string) error {
			client, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("error getting client: %s", err)
			}
			conn := client.(*mediatailor.MediaTailor)
			name := "test-playback-configuration-awsmt"
			_, err = conn.DeletePlaybackConfiguration(&mediatailor.DeletePlaybackConfigurationInput{Name: &name})
			if err != nil {
				return err
			}
			return nil
		},
	})
}

func TestAccPlaybackConfigurationResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckPlaybackConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "name", "test-playback-configuration-awsmt"),
				),
			},
			{
				Config: testAccPlaybackConfigurationUpdateResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "name", "test-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "slate_ad_url", "https://exampleurl.com/updated"),
				),
			},
		},
	})
}

func TestAccPlaybackConfigurationResourceImport(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r2"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { importPreCheck(t, "eu-central-1") },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckPlaybackConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationImportResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-playback-configuration-awsmt"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func importPreCheck(t *testing.T, region string) {
	testAccPreCheck(t)

	client, err := sharedClientForRegion(region)
	if err != nil {
		t.Fatalf("error getting client: %s", err)
	}
	conn := client.(*mediatailor.MediaTailor)

	exampleUrl := "https://exampleurl.com/"
	mpdLocation := "DISABLED"
	manifestType := "SINGLE_PERIOD"
	name := "test-playback-configuration-awsmt"
	env := "dev"
	input := mediatailor.PutPlaybackConfigurationInput{
		AdDecisionServerUrl:   &exampleUrl,
		CdnConfiguration:      &mediatailor.CdnConfiguration{AdSegmentUrlPrefix: &exampleUrl},
		DashConfiguration:     &mediatailor.DashConfigurationForPut{MpdLocation: &mpdLocation, OriginManifestType: &manifestType},
		Name:                  &name,
		VideoContentSourceUrl: &exampleUrl,
		Tags:                  map[string]*string{"Environment": &env},
	}

	_, err = conn.PutPlaybackConfiguration(&input)
	if err != nil {
		t.Fatal(diag.FromErr(err))
	}
}

func testAccCheckPlaybackConfigurationDestroy(_ *terraform.State) error {
	c := testAccProvider.Meta().(*mediatailor.MediaTailor)
	name := "test-playback-configuration-awsmt"
	_, err := c.DeletePlaybackConfiguration(&mediatailor.DeletePlaybackConfigurationInput{Name: &name})
	if err != nil {
		return err
	}
	return nil
}

func testAccPlaybackConfigurationResource() string {
	return `
resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  cdn_configuration {
    ad_segment_url_prefix = "test"
    content_segment_url_prefix = "test"
  }
  dash_configuration {
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  name = "test-playback-configuration-awsmt"
  slate_ad_url = "https://exampleurl.com/"
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/"
}

`
}
func testAccPlaybackConfigurationUpdateResource() string {
	return `
resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  cdn_configuration {
    ad_segment_url_prefix = "test-updated"
    content_segment_url_prefix = "test-updated"
  }
  dash_configuration {
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  name = "test-playback-configuration-awsmt"
  slate_ad_url = "https://exampleurl.com/updated"
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/updated"
}

`
}
func testAccPlaybackConfigurationImportResource() string {
	return `
resource "awsmt_playback_configuration" "r2" {
  ad_decision_server_url = "https://exampleurl.com/"
  cdn_configuration {
    ad_segment_url_prefix = "https://exampleurl.com/"
  }
  dash_configuration {
    mpd_location = "DISABLED"
	origin_manifest_type = "SINGLE_PERIOD"
  }
  name = "test-playback-configuration-awsmt"
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/"
}
`
}
