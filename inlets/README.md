# Inlets
## `default`
`/api/inlets/default`

Forwards a basic text message or a file to a Telegram chat.  

#### Headers
None

#### Body (Text Message)
```json
{
    "recipient_token": "<string>",
    "text": "<string>",
    "type": "<string: TEXT|FILE>",
    "origin": "<string>"
}
```

#### Body (File Message)
```json
{
    "recipient_token": "<string>",
    "file": "<base64string>",
    "filename": "<string>",
    "type": "<string: TEXT|FILE>",
    "origin": "<string>"
}
```

## `alertmanager`
`/api/inlets/alertmanager`

Accepts, transforms and forwards alerts sent by [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/) to a Telegram chat.

#### Headers
Requires the recipient token to be included as a bearer token in the request headers.
```
Authorization: Bearer <string>
```

### Body
See [webhook_config](https://prometheus.io/docs/alerting/configuration/#webhook_config).

### Example Configuration
```yaml
# alertmanager.yml

global:
  resolve_timeout: 5m

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'middleman'

receivers:
- name: 'middleman'
  webhook_configs:
  - url: 'http://localhost:8080/api/inlets/alertmanager'
    http_config:
      bearer_token: '3edf633a-eab0-45ea-9721-16c07bb8f245'
```