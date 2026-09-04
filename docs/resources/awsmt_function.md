# Resource: awsmt_function

Use this resource to manage a MediaTailor monetization function. Functions define reusable logic that MediaTailor executes at lifecycle hooks during ad insertion.

## Example Usage

### CUSTOM_OUTPUT Function

```terraform
resource "awsmt_function" "custom_output" {
  function_id   = "my-custom-output-function"
  function_type = "CUSTOM_OUTPUT"
  description   = "Resolves device type for ad targeting"

  custom_output_configuration = {
    runtime = "JSONATA"
    output = {
      "player_params.device_type" = "$lowercase(session.player_params.device)"
    }
  }

  tags = {
    Environment = "production"
  }
}
```

### HTTP_REQUEST Function

```terraform
resource "awsmt_function" "http_request" {
  function_id   = "my-http-request-function"
  function_type = "HTTP_REQUEST"
  description   = "Calls identity resolution service"

  http_request_configuration = {
    method_type                  = "GET"
    request_timeout_milliseconds = 1500
    runtime                      = "JSONATA"
    url                          = "https://identity.example.com/resolve?id={%session.player_params.user_id%}"
    output = {
      "player_params.identity_envelope" = "response.body.envelope"
    }
  }
}
```

### SEQUENTIAL_EXECUTOR Function

```terraform
resource "awsmt_function" "sequential" {
  function_id   = "my-sequential-function"
  function_type = "SEQUENTIAL_EXECUTOR"
  description   = "Runs identity resolution then enrichment"

  sequential_executor_configuration = {
    runtime              = "JSONATA"
    timeout_milliseconds = 3000
    function_list = [
      {
        function_id = "my-http-request-function"
      },
      {
        function_id   = "my-custom-output-function"
        run_condition = "$exists(temp.identity)"
      }
    ]
  }
}
```

### Using with PlaybackConfiguration

```terraform
resource "awsmt_playback_configuration" "example" {
  name                   = "my-config"
  ad_decision_server_url = "https://ads.example.com/"
  video_content_source_url = "https://content.example.com/"

  function_mapping = {
    "PRE_ADS_REQUEST"            = awsmt_function.sequential.function_id
    "PRE_SESSION_INITIALIZATION" = awsmt_function.custom_output.function_id
  }
}
```

## Arguments Reference

- `function_id` - (Required, Forces new resource) The unique identifier for the function.
- `function_type` - (Required) The type of the function. Valid values: `CUSTOM_OUTPUT`, `HTTP_REQUEST`, `SEQUENTIAL_EXECUTOR`.
- `description` - (Optional) A description of the function.
- `custom_output_configuration` - (Optional, required when function_type is CUSTOM_OUTPUT) Configuration for a CUSTOM_OUTPUT function.
  - `runtime` - (Required) The expression language. Must be `JSONATA`.
  - `output` - (Optional) A map of output bindings. Keys are namespaced output paths, values are expressions.
- `http_request_configuration` - (Optional, required when function_type is HTTP_REQUEST) Configuration for an HTTP_REQUEST function.
  - `method_type` - (Required) The HTTP method. Valid values: `GET`, `POST`.
  - `request_timeout_milliseconds` - (Required) Maximum time in milliseconds to wait for a response. Valid values: 100 to 2000.
  - `runtime` - (Required) The expression language. Must be `JSONATA`.
  - `url` - (Required) Expression that evaluates to the request URL. Use `{%...%}` for dynamic expressions.
  - `body` - (Optional) Expression that evaluates to the request body. Used with POST requests.
  - `headers` - (Optional) Map of HTTP header names to expression values.
  - `output` - (Optional) A map of output bindings that can reference the HTTP response.
- `sequential_executor_configuration` - (Optional, required when function_type is SEQUENTIAL_EXECUTOR) Configuration for a SEQUENTIAL_EXECUTOR function.
  - `function_list` - (Required) Ordered list of 1 to 10 steps.
    - `function_id` - (Optional) The identifier of the child function to execute.
    - `run_condition` - (Optional) Expression that controls whether the step runs.
  - `runtime` - (Required) The expression language. Must be `JSONATA`.
  - `timeout_milliseconds` - (Required) Maximum time in milliseconds for the entire sequence.
  - `output` - (Optional) Map of output bindings committed after all steps complete.
- `tags` - (Optional) Key-value mapping of resource tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `arn` - The Amazon Resource Name (ARN) of the function.

## Import

`awsmt_function` resources can be imported using their function_id. For example:

```sh
  $ terraform import awsmt_function.example my-function-id
```
