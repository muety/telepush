# Inlets
## `default`
`/api/inlets/default`

Forwards a basic text message or a file to a Telegram chat.  

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

Optionally, you can pass sending options with your message:
```json
{
    ...
    "options": {
        "disable_link_previews": true
    }
}
```

**Alternatively**, the default inlet's payload can be passed as _query parameters_ of a `GET` request (see [#29](https://github.com/muety/webhook2telegram/issues/29)), e.g.:
```
GET http://localhost:8080/api/messages \
    ?recipient_token=3edf633a-eab0-45ea-9721-16c07bb8f245 \
    &text=Just a test \
    &origin=Some Script \
    &type=TEXT \
    &disable_link_previews=true 
```

## `alertmanager_webhook`
`/api/inlets/alertmanager_webhook`

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
  receiver: 'webhook2telegram'

receivers:
- name: 'webhook2telegram'
  webhook_configs:
  - url: 'http://localhost:8080/api/inlets/alertmanager_webhook'
    http_config:
      bearer_token: '3edf633a-eab0-45ea-9721-16c07bb8f245'
```

## `bitbucket_webhook`
`/api/inlets/bitbucket_webhook`

Accepts, transforms and forwards events sent by [Bitbucket](https://bitbucket.org/) to a Telegram chat.

#### Parameters
Requires the recipient token to be included as a `token` URL query parameter and the `X-Event-Key` header to be set. 

### Body
See [Event Payloads](https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html).

## `webmentionio_webhook`
`/api/inlets/webmentionio_webhook`

Accepts, transforms and forwards notifications sent by [Webmention.io](https://webmention.io) to a Telegram chat.

#### Parameters
Requires the recipient token to be included as a `secret` parameter in the request's `application/json` body, which can be configured on the [settings page](https://webmention.io/settings/webhooks). 

### Body
An example payload looks as follows, however, only `secret`, `source` and `target` are utilized.
```json
{
  "secret": "1234abcd",
  "source": "http://rhiaro.co.uk/2015/11/1446953889",
  "target": "http://aaronparecki.com/notes/2015/11/07/4/indiewebcamp",
  "post": {
    "type": "entry",
    "author": {
      "name": "Amy Guy",
      "photo": "http://webmention.io/avatar/rhiaro.co.uk/829d3f6e7083d7ee8bd7b20363da84d88ce5b4ce094f78fd1b27d8d3dc42560e.png",
      "url": "http://rhiaro.co.uk/about#me"
    },
    "url": "http://rhiaro.co.uk/2015/11/1446953889",
    "published": "2015-11-08T03:38:09+00:00",
    "name": "repost of http://aaronparecki.com/notes/2015/11/07/4/indiewebcamp",
    "repost-of": "http://aaronparecki.com/notes/2015/11/07/4/indiewebcamp",
    "wm-property": "repost-of"
  }
}
```