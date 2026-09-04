resource "awsmt_function" "identity_resolver" {
  function_id   = "identity-resolver"
  function_type = "HTTP_REQUEST"
  description   = "Resolves hashed email to identity envelope via LiveRamp"

  http_request_configuration = {
    method_type                  = "GET"
    request_timeout_milliseconds = 1500
    runtime                      = "JSONATA"
    url                          = "https://identity.example.com/resolve?id={%session.player_params.hashed_email%}"
    output = {
      "player_params.identity_envelope" = "response.body.envelope"
    }
  }

  tags = {
    Environment = "production"
  }
}

resource "awsmt_function" "ad_enrichment" {
  function_id   = "ad-enrichment"
  function_type = "CUSTOM_OUTPUT"
  description   = "Enriches ad request with device metadata"

  custom_output_configuration = {
    runtime = "JSONATA"
    output = {
      "player_params.device_type" = "$lowercase(session.player_params.device)"
      "player_params.app_version" = "session.player_params.version"
    }
  }
}

resource "awsmt_function" "pre_ads_pipeline" {
  function_id   = "pre-ads-pipeline"
  function_type = "SEQUENTIAL_EXECUTOR"
  description   = "Runs identity resolution then ad enrichment before ADS request"

  sequential_executor_configuration = {
    runtime              = "JSONATA"
    timeout_milliseconds = 3000
    function_list = [
      {
        function_id = awsmt_function.identity_resolver.function_id
      },
      {
        function_id   = awsmt_function.ad_enrichment.function_id
        run_condition = "$exists(temp.identity)"
      }
    ]
  }
}

data "awsmt_function" "test" {
  function_id = awsmt_function.identity_resolver.function_id
}

output "function_arn" {
  value = data.awsmt_function.test.arn
}
