resource "awsmt_source_location" "example" {
  name = "example_source_location"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
}

data "awsmt_source_location" "example" {
  name = awsmt_source_location.example.name
}

output "source_location_out" {
  value = data.awsmt_source_location.example
}
