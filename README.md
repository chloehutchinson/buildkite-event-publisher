### Buildkite event publisher
This event publisher is designed to be run as a lambda behind API gateway or similar. It responds to Buildkite webhook events and forwards them to New Relic.

Read about Buildkite webhooks here: https://buildkite.com/docs/apis/webhooks

The lambda contains support for the following Buildkite events and correspondingly publishes custom events to New Relic:
 - `build.running` (as BuildEvent)
 - `build.finished` (as BuildEvent)
 - `job.activated` (as JobEvent)
 - `job.started` (as JobEvent)
 - `job.finished` (as JobEvent)

