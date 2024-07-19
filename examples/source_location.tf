resource "awsmt_source_location" "example" {
  name = "examplesourcelocation"
  http_configuration = {
    base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"
  }
  default_segment_delivery_configuration = {
    base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  }
}

data "awsmt_source_location" "read" {
  name = awsmt_source_location.example.name
}

output "awsmt_source_location" {
  value = data.awsmt_source_location.read
}



