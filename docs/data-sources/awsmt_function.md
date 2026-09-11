# Data Source: awsmt_function

This data source provides details about a specific MediaTailor monetization function.

## Example Usage

```terraform
data "awsmt_function" "example" {
  function_id = "my-function-id"
}
```

## Arguments Reference

- `function_id` - (Required) The identifier of the function.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `arn` - The Amazon Resource Name (ARN) of the function.
- `function_type` - The type of the function (`CUSTOM_OUTPUT`, `HTTP_REQUEST`, or `SEQUENTIAL_EXECUTOR`).
- `description` - A description of the function.
- `custom_output_configuration` - Configuration for a CUSTOM_OUTPUT function.
  - `runtime` - The expression language.
  - `output` - Map of output bindings.
- `http_request_configuration` - Configuration for an HTTP_REQUEST function.
  - `method_type` - The HTTP method.
  - `request_timeout_milliseconds` - Maximum timeout in milliseconds.
  - `runtime` - The expression language.
  - `url` - The request URL expression.
  - `body` - The request body expression.
  - `headers` - Map of HTTP headers.
  - `output` - Map of output bindings.
- `sequential_executor_configuration` - Configuration for a SEQUENTIAL_EXECUTOR function.
  - `function_list` - Ordered list of child function steps.
    - `function_id` - The child function identifier.
    - `run_condition` - Expression controlling whether the step runs.
  - `runtime` - The expression language.
  - `timeout_milliseconds` - Maximum timeout for the sequence.
  - `output` - Map of output bindings.
- `tags` - Key-value mapping of resource tags.
