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
	resourceName := "awsmt_channel.test"
	name := "test"
	stateStopped := "STOPPED"
	stateRunning := "RUNNING"
	manifestWindowSeconds := "30"
	manifestWindowSeconds2 := "40"
	minBufferTimeSeconds := "2"
	minBufferTimeSeconds2 := "3"
	minUpdatePeriodSeconds := "2"
	minUpdatePeriodSeconds2 := "3"
	PresentationDelaySeconds := "2"
	presentationDelaySeconds2 := "3"
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
				Config: basicChannel(name, stateStopped, manifestWindowSeconds, minBufferTimeSeconds, minUpdatePeriodSeconds, PresentationDelaySeconds, k1, v1, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test"),
					resource.TestMatchResourceAttr(resourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestMatchResourceAttr(resourceName, "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr(resourceName, "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(resourceName, "channel_state", "STOPPED"),
					resource.TestCheckResourceAttr(resourceName, "playback_mode", "LOOP"),
					resource.TestMatchResourceAttr(resourceName, "policy", regexp.MustCompile(`mediatailor:GetManifest`)),
					resource.TestCheckResourceAttr(resourceName, "tier", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "dev"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.manifest_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.source_group", "default"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.dash_playlist_settings.manifest_window_seconds", "30"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.dash_playlist_settings.min_buffer_time_seconds", "2"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.dash_playlist_settings.min_update_period_seconds", "2"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.dash_playlist_settings.suggested_presentation_delay_seconds", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: basicChannel(name, stateRunning, manifestWindowSeconds2, minBufferTimeSeconds2, minUpdatePeriodSeconds2, presentationDelaySeconds2, k3, v3, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test"),
					resource.TestMatchResourceAttr(resourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestCheckResourceAttr(resourceName, "name", "test"),
					resource.TestMatchResourceAttr(resourceName, "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr(resourceName, "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(resourceName, "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr(resourceName, "playback_mode", "LOOP"),
					resource.TestCheckResourceAttr(resourceName, "tier", "BASIC"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "prod"),
					resource.TestCheckResourceAttr(resourceName, "tags.Testing", "pass"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.manifest_name", "default"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.source_group", "default"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.dash_playlist_settings.manifest_window_seconds", "40"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.dash_playlist_settings.min_buffer_time_seconds", "3"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.dash_playlist_settings.min_update_period_seconds", "3"),
					resource.TestCheckResourceAttr(resourceName, "outputs.0.dash_playlist_settings.suggested_presentation_delay_seconds", "3"),
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

func TestAccChannelResourceNoState(t *testing.T) {
	noStateChannel := `
resource "awsmt_channel" "test"  {
	name = "test"
	outputs = [{
		manifest_name                = "default"
		source_group                 = "default"
		hls_playlist_settings = {
			ad_markup_type = ["DATERANGE"]
			manifest_window_seconds = "30"
		}
	}]
	playback_mode = "LOOP"
	tier = "BASIC"
	tags = {"Environment": "dev"}
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: noStateChannel,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "name", "test"),
				),
			},
		},
	})
}

func TestAccChannelResourceRunning(t *testing.T) {
	manifestWindowSeconds := "30"
	manifestWindowSeconds2 := "40"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: hlsChannel(manifestWindowSeconds),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.manifest_window_seconds", "30"),
				),
			},
			{
				Config: hlsChannel(manifestWindowSeconds2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "RUNNING"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.manifest_window_seconds", "40"),
				),
			},
		},
	})
}

func TestAccChannelResourceLoggingConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: logConfigChannel(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "enable_as_run_logs", "true"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "enable_as_run_logs", "true"),
				),
			},
			{
				Config: logConfigChannel(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "enable_as_run_logs", "false"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "enable_as_run_logs", "false"),
				),
			},
		},
	})
}

func TestAccChannelValuesNotFlickering(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: hlsChannelNoManifestWindowSeconds(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_channel.test", "channel_state", "STOPPED"),
					resource.TestCheckResourceAttr("data.awsmt_channel.test", "outputs.0.hls_playlist_settings.ad_markup_type.0", "DATERANGE"),
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

func basicChannel(name, state, manifestWindowSeconds, minBufferTimeSeconds, minUpdatePeriodSeconds, presentationDelaySeconds, k1, v1, k2, v2 string) string {
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
  					policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:985600762523:channel/test\"}]}"
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
				`, name, state, manifestWindowSeconds, minBufferTimeSeconds, minUpdatePeriodSeconds, presentationDelaySeconds, k1, v1, k2, v2,
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

func hlsChannelNoManifestWindowSeconds() string {
	return `
				resource "awsmt_channel" "test"  {
  					name = "test"
  					channel_state = "STOPPED"
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							ad_markup_type = ["DATERANGE"]
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

func logConfigChannel(enable bool) string {
	return fmt.Sprintf(`
				resource "awsmt_channel" "test"  {
  					name = "test"
  					channel_state = "RUNNING"
					enable_as_run_logs = %[1]v
  					outputs = [{
    					manifest_name                = "default"
						source_group                 = "default"
    					hls_playlist_settings = {
							ad_markup_type = ["DATERANGE"]
							manifest_window_seconds = "60"
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
				`, enable)
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
			policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:985600762523:channel/test\"}]}"
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
