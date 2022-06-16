resource "awsmt_channel" "test" {
  name = "dev-channel"
  outputs {
    manifest_name                = "default"
    source_group                 = "default"
    hls_manifest_windows_seconds = 30
  }
  playback_mode = "LOOP"
  policy = "{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"AllowAnonymous\", \"Effect\": \"Allow\", \"Principal\": \"*\", \"Action\": \"mediatailor:GetManifest\", \"Resource\": \"arn:aws:mediatailor:eu-central-1:319158032161:channel/dev-channel\"}]}"
  tier = "BASIC"
}

data "awsmt_channel" "read" {
  name = awsmt_channel.test.name
}

output "channel_out" {
  value = data.awsmt_channel.read
}