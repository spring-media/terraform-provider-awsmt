resource "awsmt_channel" "test"  {
  name = "test"
  channel_state = "RUNNING"
  outputs = [{
    manifest_name                = "default"
    source_group                 = "default"
    hls_playlist_settings = {
      ad_markup_type = ["DATERANGE"]
      manifest_window_seconds = 30
    }
  }]
  playback_mode = "LOOP"
  tier = "BASIC"
  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:985600762523:channel/test\"}]}"
  tags = {"Environment": "dev"}
}

data "awsmt_channel" "test" {
  name = awsmt_channel.test.name
}

output "channel_out" {
  value = data.awsmt_channel.test
}