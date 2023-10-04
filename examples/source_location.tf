resource "awsmt_source_location" "example" {
  source_location_name = "examplesourcelocation"
  http_configuration = {
    hc_base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"
  }
  default_segment_delivery_configuration = {
    dsdc_base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  }
}
data "awsmt_source_location" "read" {
  source_location_name = awsmt_source_location.example.source_location_name
}

output "awsmt_source_location" {
  value = data.awsmt_source_location.read
}



