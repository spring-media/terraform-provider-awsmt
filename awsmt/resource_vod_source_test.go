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

func init() {
	resource.AddTestSweepers("test_vod_source", &resource.Sweeper{
		Name: "test_vod_source",
		F: func(region string) error {
			client, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("error getting client: %s", err)
			}
			conn := client.(*mediatailor.MediaTailor)
			names := map[string]string{"test_source_location_basic": "vod_source_test_basic", "test_source_location_update": "vod_source_test_basic", "test_source_location_tags": "vod_source_test_basic", "vod_basic_sl": "vod_source_data_source_test"}
			for k, v := range names {
				_, err = conn.DeleteVodSource(&mediatailor.DeleteVodSourceInput{SourceLocationName: &k, VodSourceName: &v})
				if err != nil {
					if !strings.Contains(err.Error(), "NotFound") {
						return err
					}
				}
				_, err = conn.DeleteSourceLocation(&mediatailor.DeleteSourceLocationInput{SourceLocationName: &k})
				if err != nil {
					if !strings.Contains(err.Error(), "NotFound") {
						return err
					}
				}
			}
			return nil
		},
	})
}

func TestAccVodSourceResource_basic(t *testing.T) {
	rName := "vod_source_test_basic"
	resourceName := "awsmt_vod_source.test"
	SourceLocationName := "test_source_location_basic"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckVodSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSourceConfig(SourceLocationName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func TestAccVodSourceResource_update(t *testing.T) {
	rName := "vod_source_test_basic"
	resourceName := "awsmt_vod_source.test"
	SourceLocationName := "test_source_location_update"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckVodSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSourceConfig_update(SourceLocationName, rName, "/"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "source_location_name", SourceLocationName),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.path", "/"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.type", "HLS"),
				),
			},
			{
				Config: testAccVodSourceConfig_update(SourceLocationName, rName, "/test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "source_location_name", SourceLocationName),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.path", "/test"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr(resourceName, "http_package_configurations.0.type", "HLS"),
				),
			},
		},
	})
}

func TestAccVodSourceResource_tags(t *testing.T) {
	rName := "vod_source_test_basic"
	resourceName := "awsmt_vod_source.test"
	SourceLocationName := "test_source_location_tags"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		CheckDestroy:      testAccCheckVodSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSourceConfig_tags(SourceLocationName, rName, "a", "b", "c", "d"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "source_location_name", SourceLocationName),
					resource.TestCheckResourceAttr(resourceName, "tags.a", "b"),
					resource.TestCheckResourceAttr(resourceName, "tags.c", "d"),
				),
			},
			{
				Config: testAccVodSourceConfig_tags(SourceLocationName, rName, "e", "f", "g", "h"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "source_location_name", SourceLocationName),
					resource.TestCheckResourceAttr(resourceName, "tags.e", "f"),
					resource.TestCheckResourceAttr(resourceName, "tags.g", "h"),
				),
			},
		},
	})
}

func testAccCheckVodSourceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*mediatailor.MediaTailor)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "awsmt_vod_source" {
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

		input := &mediatailor.DescribeVodSourceInput{VodSourceName: aws.String(resourceName), SourceLocationName: aws.String("vod-source-test")}
		_, err := conn.DescribeVodSource(input)

		if err != nil && strings.Contains(err.Error(), "NotFound") {
			continue
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func testAccVodSourceConfig(sourceLocationName, vodSourceName string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "example"{
  name = "%[1]s"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
}

resource "awsmt_vod_source" "test" {
  http_package_configurations {
    path = "/"
    source_group = "default"
    type = "HLS"
  }
  source_location_name = awsmt_source_location.example.name
  name = "%[2]s"
}
`, sourceLocationName, vodSourceName)
}

func testAccVodSourceConfig_update(sourceLocationName, vodSourceName, path string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "example"{
  name = "%[1]s"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
}

resource "awsmt_vod_source" "test" {
  http_package_configurations {
    path = "%[3]s"
    source_group = "default"
    type = "HLS"
  }
  source_location_name = awsmt_source_location.example.name
  name = "%[2]s"
}
`, sourceLocationName, vodSourceName, path)
}

func testAccVodSourceConfig_tags(sourceLocationName, vodSourceName, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "example"{
  name = "%[1]s"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
}

resource "awsmt_vod_source" "test" {
  http_package_configurations {
    path = "%[3]s"
    source_group = "default"
    type = "HLS"
  }
  source_location_name = awsmt_source_location.example.name
  tags = {
    "%[3]s": "%[4]s",
	"%[5]s": "%[6]s",
  }
  name = "%[2]s"
}
`, sourceLocationName, vodSourceName, k1, v1, k2, v2)
}
