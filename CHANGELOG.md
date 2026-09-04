# Changelog

## Unreleased

### Features

- **New Resource: `awsmt_function`** — Manage MediaTailor monetization functions (CUSTOM_OUTPUT, HTTP_REQUEST, SEQUENTIAL_EXECUTOR). Functions define reusable logic that MediaTailor executes at lifecycle hooks during ad insertion.
- **New Data Source: `awsmt_function`** — Read existing MediaTailor monetization functions.
- **playback_configuration**: Added `ad_conditioning_configuration` attribute to control ad transcoding behavior (`TRANSCODE` or `NONE`).
- **playback_configuration**: Added `ad_decision_server_configuration` attribute with nested `http_request` block for customizing HTTP requests to the ADS (method, headers, body, compression).
- **playback_configuration**: Added `function_mapping` attribute to map lifecycle hooks (`PRE_SESSION_INITIALIZATION`, `PRE_ADS_REQUEST`) to MediaTailor function IDs.
- **playback_configuration**: Added `insertion_mode` attribute to control stitched vs guided ad insertion (`STITCHED_ONLY`, `PLAYER_SELECT`).
- **playback_configuration**: Added `log_configuration_ads_interaction_log` attribute for fine-grained ADS log event filtering (exclude/publish opt-in event types).
- **playback_configuration**: Added `log_configuration_manifest_service_interaction_log` attribute for fine-grained manifest service log event filtering (exclude/publish opt-in event types).

### Dependencies

- Upgraded `github.com/aws/aws-sdk-go-v2/service/mediatailor` from v1.47.0 to v1.58.0
- Upgraded `github.com/aws/aws-sdk-go-v2` from v1.36.3 to v1.41.7
- Upgraded `github.com/aws/smithy-go` from v1.22.3 to v1.25.1

### Breaking Changes

None. All new fields are optional and do not alter existing behavior.

### Design Notes

- `ad_decision_server_configuration` is a separate optional nested attribute alongside the existing `ad_decision_server_url`. This matches the SDK structure where the URL and HTTP request configuration are independent fields.
- Log interaction configs (`ads_interaction_log`, `manifest_service_interaction_log`) are flattened to the top level with a `log_configuration_` prefix, consistent with the existing `log_configuration_percent_enabled` and `log_configuration_enabled_logging_strategies` pattern. This follows the existing @ADR decision that the Provider Framework does not allow computed blocks, so log configuration is flattened into the resource.
- The `function_mapping` attribute requires that referenced functions already exist in the account (created via the MediaTailor Functions API).
