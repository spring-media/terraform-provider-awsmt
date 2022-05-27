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

func TestAccSourceLocationResource_basic(t *testing.T) {
	rName := "source_location_test_basic"
	resourceName := "awsmt_source_location.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckSourceLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceLocationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_location_name", rName),
					resource.TestCheckResourceAttr(resourceName, "http_configuration_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
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

func TestAccSourceLocationResource_recreate(t *testing.T) {
	rName := "source_location_test_recreate"
	resourceName := "awsmt_source_location.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckSourceLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceLocationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_location_name", rName),
					resource.TestCheckResourceAttr(resourceName, "http_configuration_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
				),
			},
			{
				Taint:  []string{resourceName},
				Config: testAccSourceLocationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_location_name", rName),
					resource.TestCheckResourceAttr(resourceName, "http_configuration_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
				),
			},
		},
	})
}

func TestAccSourceLocationResource_update(t *testing.T) {
	rName := "source_location_test_update"
	resourceName := "awsmt_source_location.test_update"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckSourceLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceLocationConfig_update(rName, "example", "https://example.com", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_location_name", rName),
					resource.TestCheckResourceAttr(resourceName, "default_segment_delivery_configuration_url", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "segment_delivery_configurations.0.name", "example"),
					resource.TestCheckResourceAttr(resourceName, "segment_delivery_configurations.0.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
					resource.TestCheckResourceAttr(resourceName, "http_configuration_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
				),
			},
			{
				Config: testAccSourceLocationConfig_update(rName, "test", "https://test.com", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_location_name", rName),
					resource.TestCheckResourceAttr(resourceName, "default_segment_delivery_configuration_url", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "segment_delivery_configurations.0.name", "test"),
					resource.TestCheckResourceAttr(resourceName, "segment_delivery_configurations.0.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
					resource.TestCheckResourceAttr(resourceName, "http_configuration_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
				),
			},
		},
	})
}

func TestAccSourceLocationResource_tags(t *testing.T) {
	rName := "source_location_test_tags"
	resourceName := "awsmt_source_location.test_tags"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckSourceLocationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceLocationConfig_tags(rName, "a", "b", "c", "d"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_location_name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.a", "b"),
					resource.TestCheckResourceAttr(resourceName, "tags.c", "d"),
				),
			},
			{
				Config: testAccSourceLocationConfig_tags(rName, "e", "f", "g", "h"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "source_location_name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.e", "f"),
					resource.TestCheckResourceAttr(resourceName, "tags.g", "h"),
				),
			},
		},
	})
}

func testAccCheckSourceLocationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*mediatailor.MediaTailor)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "awsmt_source_location" {
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

		input := &mediatailor.DescribeSourceLocationInput{SourceLocationName: aws.String(resourceName)}
		_, err := conn.DescribeSourceLocation(input)

		if err != nil && strings.Contains(err.Error(), "NotFound") {
			continue
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func testAccSourceLocationConfig(rName string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "test"{
  source_location_name = "%[1]s"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
}
`, rName)
}

func testAccSourceLocationConfig_update(rName, exampleString, exampleUrl, baseUrl string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "test_update"{
  access_configuration {
    access_type = "S3_SIGV4"
  }
  default_segment_delivery_configuration_url = "%[3]s"
  http_configuration_url = "%[4]s"
  source_location_name = "%[1]s"
  segment_delivery_configurations {
    base_url = "%[4]s"
    name =     "%[2]s"
  }
}
`, rName, exampleString, exampleUrl, baseUrl)
}

func testAccSourceLocationConfig_tags(rName, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "test_tags"{
  source_location_name = "%[1]s"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  tags = {
    "%[2]s": "%[3]s",
	"%[4]s": "%[5]s",
  }
}
`, rName, k1, v1, k2, v2)
}
