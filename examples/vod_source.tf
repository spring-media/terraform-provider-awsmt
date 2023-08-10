resource "awsmt_vod_source" "test" {
  http_package_configurations = [{
    path = "/"
    source_group = "default"
    type = "HLS"
  }]
  source_location_name = awsmt_source_location.example.source_location_name
  vod_source_name = "vod_source_example"
}

data "awsmt_vod_source" "data_test" {
  source_location_name = awsmt_source_location.example.source_location_name
  vod_source_name = awsmt_vod_source.test.vod_source_name
}

output "vod_source_out" {
  value = data.awsmt_vod_source.data_test
}

