resource "awsmt_channel" "test"  {
  channel_name = "test"
  channel_state = "RUNNING"
  outputs = [{
    manifest_name                = "default"
    source_group                 = "default"
    dash_playlist_settings = {
      manifest_window_seconds = 30
      min_buffer_time_seconds = 2
      min_update_period_seconds = 2
      suggested_presentation_delay_seconds = 2
    }
  }]
  playback_mode = "LOOP"
  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:985600762523:channel/test\"}]}"
  tier = "BASIC"
  tags = {"Environment": "dev", "Name": "test"}
}

data "awsmt_channel" "test" {
  channel_name = awsmt_channel.test.channel_name
}
output "channel_out" {
  value = data.awsmt_channel.test
}