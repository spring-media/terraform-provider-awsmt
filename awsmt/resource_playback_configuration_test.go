package awsmt

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccPlaybackConfigurationMinimal(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r2"
	name := "test-acc-playback-configuration-minimal"
	adUrl := "https://www.foo.de/"
	videoSourceUrl := "https://www.bar.at"
	adUrl2 := "https://www.biz.ch"
	videoSourceUrl2 := "https://www.buzz.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: minimalPlaybackConfiguration(name, adUrl, videoSourceUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ad_decision_server_url", adUrl),
					resource.TestCheckResourceAttr(resourceName, "video_content_source_url", videoSourceUrl),
				),
			},
			{
				Config: minimalPlaybackConfiguration(name, adUrl2, videoSourceUrl2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ad_decision_server_url", adUrl2),
					resource.TestCheckResourceAttr(resourceName, "video_content_source_url", videoSourceUrl2),
				),
			},
		},
	})
}

func TestAccPlaybackConfigurationConfigurationAliases(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r4"
	name := "test-acc-playback-configuration-minimal"
	adUrl := "https://www.foo.de/"
	videoSourceUrl := "https://www.bar.at"
	aliases := map[string]map[string]string{
		"player_params.foo": {
			"player_params.bar": "player_params.boo",
		},
	}
	aliases2 := map[string]map[string]string{
		"player_params.foo": {
			"player_params.bar": "player_params.buzz",
		},
	}

	aliasesJson, _ := json.Marshal(aliases)
	aliases2Json, _ := json.Marshal(aliases2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configurationAliasesPlaybackConfiguration(name, adUrl, videoSourceUrl, string(aliasesJson)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ad_decision_server_url", adUrl),
					resource.TestCheckResourceAttr(resourceName, "video_content_source_url", videoSourceUrl),
					resource.TestCheckResourceAttr(resourceName, "configuration_aliases.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration_aliases.player_params.foo.player_params.bar", "player_params.boo"),
				),
			},
			{
				Config: configurationAliasesPlaybackConfiguration(name, adUrl, videoSourceUrl, string(aliases2Json)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ad_decision_server_url", adUrl),
					resource.TestCheckResourceAttr(resourceName, "video_content_source_url", videoSourceUrl),
					resource.TestCheckResourceAttr(resourceName, "configuration_aliases.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration_aliases.player_params.foo.player_params.bar", "player_params.buzz"),
				),
			},
		},
	})
}

func TestAccPlaybackConfigurationLogPercentage(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r3"
	name := "test-acc-playback-configuration-log-percentage"
	adUrl := "https://www.foo.de/"
	videoSourceUrl := "https://www.bar.at"
	p1 := 5
	p2 := 8
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: logPercentagePlaybackConfiguration(name, adUrl, videoSourceUrl, p1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "log_configuration_percent_enabled", fmt.Sprintf("%d", p1)),
				),
			},
			{
				Config: logPercentagePlaybackConfiguration(name, adUrl, videoSourceUrl, p2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "log_configuration_percent_enabled", fmt.Sprintf("%d", p2)),
				),
			},
		},
	})
}

func TestAccPlaybackConfigurationLoggingStrategies(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r3"
	name := "test-acc-playback-configuration-enabled-logging-strategies"
	adUrl := "https://www.foo.de/"
	videoSourceUrl := "https://www.bar.at"
	p1 := `["LEGACY_CLOUDWATCH"]`
	p2 := `["LEGACY_CLOUDWATCH", "VENDED_LOGS"]`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loggingStrategiesPlaybackConfiguration(name, adUrl, videoSourceUrl, p1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "log_configuration_enabled_logging_strategies.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_configuration_enabled_logging_strategies.0", "LEGACY_CLOUDWATCH"),
				),
			},
			{
				Config: loggingStrategiesPlaybackConfiguration(name, adUrl, videoSourceUrl, p2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "log_configuration_enabled_logging_strategies.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "log_configuration_enabled_logging_strategies.0", "LEGACY_CLOUDWATCH"),
					resource.TestCheckResourceAttr(resourceName, "log_configuration_enabled_logging_strategies.1", "VENDED_LOGS"),
				),
			},
		},
	})
}

func TestAccPlaybackConfigurationCreationFail(t *testing.T) {
	name := "test-acc-playback-configuration-delete"
	adUrl := "invalid"
	videoSourceUrl := "https://www.bar.at"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      minimalPlaybackConfiguration(name, adUrl, videoSourceUrl),
				ExpectError: regexp.MustCompile(`.*AdDecisionServerUrl requires a valid URL.*`),
			},
		},
	})
}

func TestAccPlaybackConfigurationUpdateFail(t *testing.T) {
	resourceName := "awsmt_playback_configuration.r2"
	name := "test-acc-playback-configuration-delete"
	adUrl := "https://www.buzz.com"
	adUrl2 := "invalid"
	videoSourceUrl := "https://www.bar.at"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: minimalPlaybackConfiguration(name, adUrl, videoSourceUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", name),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "ad_decision_server_url", adUrl),
				),
			},
			{
				Config:      minimalPlaybackConfiguration(name, adUrl2, videoSourceUrl),
				ExpectError: regexp.MustCompile(`.*AdDecisionServerUrl requires a valid URL.*`),
			},
		},
	})
}

func TestAccPlaybackConfigurationResource(t *testing.T) {
	name := "example-playback-configuration-awsmt"
	adUrl := "https://exampleurl.com/"
	adUrl2 := "https://exampleurl2.com/"
	bumperE := "https://wxample.com/endbumper"
	bumperE2 := "https://wxample.com/endbumper2"
	bumperS := "https://wxample.com/startbumper"
	bumperS2 := "https://wxample.com/startbumper2"
	cdnUrl := "https://exampleurl.com/"
	cdnUrl2 := "https://exampleurl2.com/"
	maxD := "2"
	maxD2 := "3"
	pS := "2"
	pS2 := "3"
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
			// Create and Read testing
			{
				Config: completePlaybackConfiguration(name, adUrl, bumperE, bumperS, cdnUrl, maxD, pS, k1, v1, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "id", name),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "ad_decision_server_url", adUrl),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "avail_suppression.fill_policy", "FULL_AVAIL_ONLY"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "avail_suppression.mode", "BEHIND_LIVE_EDGE"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "avail_suppression.value", "00:00:00"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "bumper.end_url", bumperE),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "bumper.start_url", bumperS),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "cdn_configuration.ad_segment_url_prefix", cdnUrl),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "dash_configuration.mpd_location", "DISABLED"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "dash_configuration.origin_manifest_type", "SINGLE_PERIOD"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "live_pre_roll_configuration.ad_decision_server_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "live_pre_roll_configuration.max_duration_seconds", maxD),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "manifest_processing_rules.ad_marker_passthrough.enabled", "false"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "name", name),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "personalization_threshold_seconds", pS),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "slate_ad_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "tags.Environment", "dev"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "tags.Testing", "pass"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "video_content_source_url", "https://exampleurl.com/"),
				),
			},
			// ImportState testing
			{
				ResourceName: "awsmt_playback_configuration.r1",
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: completePlaybackConfiguration(name, adUrl2, bumperE2, bumperS2, cdnUrl2, maxD2, pS2, k3, v3, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "name", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "personalization_threshold_seconds", "3"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "ad_decision_server_url", "https://exampleurl2.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "bumper.end_url", "https://wxample.com/endbumper2"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "bumper.start_url", "https://wxample.com/startbumper2"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "cdn_configuration.ad_segment_url_prefix", "https://exampleurl2.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "tags.Environment", "prod"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "tags.Testing", "pass"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func minimalPlaybackConfiguration(name, adUrl, videoSourceUrl string) string {
	return fmt.Sprintf(`
		resource "awsmt_playback_configuration" "r2" {
			ad_decision_server_url = "%[2]s"
			name = "%[1]s"
			video_content_source_url = "%[3]s"
		}
		`, name, adUrl, videoSourceUrl,
	)
}

func configurationAliasesPlaybackConfiguration(name, adUrl, videoSourceUrl, configAliases string) string {
	return fmt.Sprintf(`
		resource "awsmt_playback_configuration" "r4" {
			ad_decision_server_url = "%[2]s"
			name = "%[1]s"
			video_content_source_url = "%[3]s"
			configuration_aliases = %[4]s
		}
		`, name, adUrl, videoSourceUrl, configAliases,
	)
}

func logPercentagePlaybackConfiguration(name, adUrl, videoSourceUrl string, logPercentage int) string {
	return fmt.Sprintf(`
		resource "awsmt_playback_configuration" "r3" {
			ad_decision_server_url = "%[2]s"
			name = "%[1]s"
			video_content_source_url = "%[3]s"
			log_configuration_percent_enabled = %[4]d
		}
		`, name, adUrl, videoSourceUrl, logPercentage,
	)
}

func loggingStrategiesPlaybackConfiguration(name, adUrl, videoSourceUrl, loggingStrategies string) string {
	return fmt.Sprintf(`
		resource "awsmt_playback_configuration" "r3" {
			ad_decision_server_url = "%[2]s"
			name = "%[1]s"
			video_content_source_url = "%[3]s"
			log_configuration_enabled_logging_strategies = %[4]s
		}
		`, name, adUrl, videoSourceUrl, loggingStrategies,
	)
}

func completePlaybackConfiguration(name, adUrl, bumperE, bumperS, cdnUrl, maxD, pS, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`
		resource "awsmt_playback_configuration" "r1" {
			ad_decision_server_url = "%[2]s"
			avail_suppression = {
				fill_policy = "FULL_AVAIL_ONLY"
				mode = "BEHIND_LIVE_EDGE"
				value = "00:00:00"
			}
			bumper = {
				end_url = "%[3]s"
				start_url = "%[4]s"
			}
			cdn_configuration = {
				ad_segment_url_prefix = "%[5]s"
			}
			dash_configuration = {
				mpd_location = "DISABLED",
				origin_manifest_type = "SINGLE_PERIOD"
			}
			live_pre_roll_configuration = {
				ad_decision_server_url = "%[2]s"
				max_duration_seconds = "%[6]s"
			}
			manifest_processing_rules = {
				ad_marker_passthrough = {
					enabled = "false"
				}
			}
			name = "%[1]s"
			personalization_threshold_seconds = "%[7]s"
			slate_ad_url = "https://exampleurl.com/"
			tags = {
				"%[8]s": "%[9]s",
				"%[10]s": "%[11]s"
			}
			video_content_source_url = "%[2]s"
		}
		`, name, adUrl, bumperE, bumperS, cdnUrl, maxD, pS, k1, v1, k2, v2,
	)
}