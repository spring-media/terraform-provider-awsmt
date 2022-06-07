package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccLiveSourceResource_basic(t *testing.T) {
	rName := "live_source_test_basic"
	resourceName := "awsmt_live_source.test"
	SourceLocationName := "live_source_basic_sl"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckLiveSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveSourceConfig(SourceLocationName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "live_source_name", rName),
					resource.TestCheckResourceAttr(resourceName, "source_location_name", SourceLocationName),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.type", "HLS"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateVerify: true,
				ImportState:       true,
			},
		},
	})
}

func TestAccLiveSourceResource_update(t *testing.T) {
	rName := "live_source_test_basic"
	resourceName := "awsmt_live_source.test"
	SourceLocationName := "live_source_update_sl"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckLiveSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveSourceConfig_update(SourceLocationName, rName, "/"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "live_source_name", rName),
					resource.TestCheckResourceAttr(resourceName, "source_location_name", SourceLocationName),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.type", "HLS"),
				),
			},
			{
				Config: testAccLiveSourceConfig_update(SourceLocationName, rName, "/test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "live_source_name", rName),
					resource.TestCheckResourceAttr(resourceName, "source_location_name", SourceLocationName),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.path", "/test"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.type", "HLS"),
				),
			},
		},
	})
}

func TestAccLiveSourceResource_tags(t *testing.T) {
	rName := "live_source_test_basic"
	resourceName := "awsmt_live_source.test"
	SourceLocationName := "live_source_tags_sl"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckLiveSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveSourceConfig_tags(SourceLocationName, rName, "a", "b", "c", "d"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "live_source_name", rName),
					resource.TestCheckResourceAttr(resourceName, "source_location_name", SourceLocationName),
					resource.TestCheckResourceAttr(resourceName, "tags.a", "b"),
					resource.TestCheckResourceAttr(resourceName, "tags.c", "d"),
				),
			},
			{
				Config: testAccLiveSourceConfig_tags(SourceLocationName, rName, "e", "f", "g", "h"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "live_source_name", rName),
					resource.TestCheckResourceAttr(resourceName, "source_location_name", SourceLocationName),
					resource.TestCheckResourceAttr(resourceName, "tags.e", "f"),
					resource.TestCheckResourceAttr(resourceName, "tags.g", "h"),
				),
			},
		},
	})
}

func testAccCheckLiveSourceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*mediatailor.MediaTailor)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "awsmt_live_source" {
			continue
		}

		var resourceName string

		if arn.IsARN(rs.Primary.ID) {
			resourceArn, err := arn.Parse(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("error parsing resource arn: %s.\n%s", err, rs.Primary.ID)
			}
			arnSections := strings.Split(resourceArn.Resource, "/")
			resourceName = arnSections[len(arnSections)-1]
		} else {
			resourceName = rs.Primary.ID
		}

		input := &mediatailor.DescribeLiveSourceInput{LiveSourceName: aws.String(resourceName), SourceLocationName: aws.String("source-location-testacc")}
		_, err := conn.DescribeLiveSource(input)

		if err != nil && strings.Contains(err.Error(), "NotFound") {
			continue
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func testAccLiveSourceConfig(sourceLocationName, liveSourceName string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "example"{
  source_location_name = "%[1]s"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
}

resource "awsmt_live_source" "test" {
  http_package_configurations {
    path = "/"
    source_group = "default"
    type = "HLS"
  }
  source_location_name = awsmt_source_location.example.source_location_name
  live_source_name = "%[2]s"
}
`, sourceLocationName, liveSourceName)
}

func testAccLiveSourceConfig_update(sourceLocationName, liveSourceName, path string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "example"{
  source_location_name = "%[1]s"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
}

resource "awsmt_live_source" "test" {
  http_package_configurations {
    path = "%[3]s"
    source_group = "default"
    type = "HLS"
  }
  source_location_name = awsmt_source_location.example.source_location_name
  live_source_name = "%[2]s"
}
`, sourceLocationName, liveSourceName, path)
}

func testAccLiveSourceConfig_tags(sourceLocationName, liveSourceName, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "example"{
  source_location_name = "%[1]s"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
}

resource "awsmt_live_source" "test" {
  http_package_configurations {
    path = "%[3]s"
    source_group = "default"
    type = "HLS"
  }
  source_location_name = awsmt_source_location.example.source_location_name
  tags = {
    "%[3]s": "%[4]s",
	"%[5]s": "%[6]s",
  }
  live_source_name = "%[2]s"
}
`, sourceLocationName, liveSourceName, k1, v1, k2, v2)
}
