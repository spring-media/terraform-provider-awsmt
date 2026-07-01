package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
	"time"
)

func TestAccPrefetchScheduleSingle(t *testing.T) {
	resourceName := "awsmt_prefetch_schedule.test"
	name := "test_acc_prefetch_single"
	playbackConfigName := "test_acc_prefetch_pc"

	consumptionEnd := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	retrievalEnd := time.Now().Add(12 * time.Hour).UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: singlePrefetchScheduleConfig(name, playbackConfigName, consumptionEnd, retrievalEnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "playback_configuration_name", playbackConfigName),
					resource.TestCheckResourceAttr(resourceName, "schedule_type", "SINGLE"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "consumption.end_time", consumptionEnd),
					resource.TestCheckResourceAttr(resourceName, "retrieval.end_time", retrievalEnd),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s/%s", playbackConfigName, name),
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccPrefetchScheduleWithAvailMatching(t *testing.T) {
	resourceName := "awsmt_prefetch_schedule.test"
	name := "test_acc_prefetch_avail"
	playbackConfigName := "test_acc_prefetch_pc2"

	consumptionEnd := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	retrievalEnd := time.Now().Add(12 * time.Hour).UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: prefetchScheduleWithAvailMatchingConfig(name, playbackConfigName, consumptionEnd, retrievalEnd),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consumption.avail_matching_criteria.0.dynamic_variable", "scte.event_id"),
					resource.TestCheckResourceAttr(resourceName, "consumption.avail_matching_criteria.0.operator", "EQUALS"),
					resource.TestCheckResourceAttr(resourceName, "retrieval.dynamic_variables.scte.event_id", "12345"),
				),
			},
		},
	})
}

func TestAccPrefetchScheduleRecurring(t *testing.T) {
	resourceName := "awsmt_prefetch_schedule.test"
	name := "test_acc_prefetch_recurring"
	playbackConfigName := "test_acc_prefetch_pc3"

	endTime := time.Now().Add(23 * time.Hour).UTC().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: recurringPrefetchScheduleConfig(name, playbackConfigName, endTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "schedule_type", "RECURRING"),
					resource.TestCheckResourceAttr(resourceName, "recurring_prefetch_configuration.end_time", endTime),
					resource.TestCheckResourceAttr(resourceName, "recurring_prefetch_configuration.recurring_consumption.retrieved_ad_expiration_seconds", "300"),
					resource.TestCheckResourceAttr(resourceName, "recurring_prefetch_configuration.recurring_retrieval.delay_after_avail_end_seconds", "60"),
				),
			},
		},
	})
}

func TestAccPrefetchScheduleErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      prefetchScheduleErrorConfig(),
				ExpectError: regexp.MustCompile("Error creating prefetch schedule"),
			},
		},
	})
}

func minimalPlaybackConfig(name string) string {
	return fmt.Sprintf(`
resource "awsmt_playback_configuration" "test" {
  ad_decision_server_url   = "https://exampleurl.com/"
  name                     = "%[1]s"
  video_content_source_url = "https://exampleurl.com/"
}
`, name)
}

func singlePrefetchScheduleConfig(name, playbackConfigName, consumptionEnd, retrievalEnd string) string {
	return minimalPlaybackConfig(playbackConfigName) + fmt.Sprintf(`
resource "awsmt_prefetch_schedule" "test" {
  name                        = "%[1]s"
  playback_configuration_name = awsmt_playback_configuration.test.name
  schedule_type               = "SINGLE"
  consumption = {
    end_time = "%[2]s"
  }
  retrieval = {
    end_time = "%[3]s"
  }
}
`, name, consumptionEnd, retrievalEnd)
}

func prefetchScheduleWithAvailMatchingConfig(name, playbackConfigName, consumptionEnd, retrievalEnd string) string {
	return minimalPlaybackConfig(playbackConfigName) + fmt.Sprintf(`
resource "awsmt_prefetch_schedule" "test" {
  name                        = "%[1]s"
  playback_configuration_name = awsmt_playback_configuration.test.name
  schedule_type               = "SINGLE"
  consumption = {
    end_time = "%[2]s"
    avail_matching_criteria = [{
      dynamic_variable = "scte.event_id"
      operator         = "EQUALS"
    }]
  }
  retrieval = {
    end_time = "%[3]s"
    dynamic_variables = {
      "scte.event_id" = "12345"
    }
  }
}
`, name, consumptionEnd, retrievalEnd)
}

func recurringPrefetchScheduleConfig(name, playbackConfigName, endTime string) string {
	return minimalPlaybackConfig(playbackConfigName) + fmt.Sprintf(`
resource "awsmt_prefetch_schedule" "test" {
  name                        = "%[1]s"
  playback_configuration_name = awsmt_playback_configuration.test.name
  schedule_type               = "RECURRING"
  recurring_prefetch_configuration = {
    end_time = "%[2]s"
    recurring_consumption = {
      retrieved_ad_expiration_seconds = 300
    }
    recurring_retrieval = {
      delay_after_avail_end_seconds = 60
    }
  }
}
`, name, endTime)
}

func prefetchScheduleErrorConfig() string {
	return `
resource "awsmt_prefetch_schedule" "test" {
  name                        = "test_acc_prefetch_error"
  playback_configuration_name = "nonexistent_playback_config"
  schedule_type               = "SINGLE"
  consumption = {
    end_time = "2099-01-01T00:00:00Z"
  }
  retrieval = {
    end_time = "2099-01-01T00:00:00Z"
  }
}
`
}
