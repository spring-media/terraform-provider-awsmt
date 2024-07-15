resource "awsmt_vod_source" "test" {
  http_package_configurations = [{
    path = "/"
    source_group = "default"
    type = "HLS"
  }]
  source_location_name = awsmt_source_location.example.name
  name = "vod_source_example"
}

data "awsmt_vod_source" "data_test" {
  source_location_name = awsmt_source_location.example.name
  name = awsmt_vod_source.test.name
}

output "vod_source_out" {
  value = data.awsmt_vod_source.data_test
}

