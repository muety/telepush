# Inlets
## `default`
`/api/inlets/default/<recipient>`

Forwards a basic text message or a file to a Telegram chat.  

#### Body (Text Message)
```json
{
    "text": "<string>",
    "type": "<string: TEXT|FILE>",
    "origin": "<string>"
}
```

#### Body (File Message)
```json
{
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

**Alternatively**, the default inlet's payload can be passed as _query parameters_ of a `GET` request (see [#29](https://github.com/muety/telepush/issues/29)), e.g.:
```
GET http://localhost:8080/api/messages/<recipient> \
    ?text=Just a test \
    &origin=Some Script \
    &type=TEXT \
    &disable_link_previews=true 
```

## `alertmanager`
`/api/inlets/alertmanager/<recipient>`

Accepts, transforms and forwards alerts sent by [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/) to a Telegram chat.

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
  receiver: 'telepush'

receivers:
- name: 'telepush'
  webhook_configs:
  - url: 'http://localhost:8080/api/inlets/alertmanager_webhook/5hd9mx'
```

## `bitbucket`
`/api/inlets/bitbucket/<recipient>`

Accepts, transforms and forwards events sent by [Bitbucket](https://bitbucket.org/) to a Telegram chat.

#### Parameters
Requires the `X-Event-Key` header to be set. 

### Body
See [Event Payloads](https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html).

## `webmentionio`
`/api/inlets/webmentionio/<recipient>`

Accepts, transforms and forwards notifications sent by [Webmention.io](https://webmention.io) to a Telegram chat.

### Body
An example payload looks as follows, however, only `secret`, `source` and `target` are utilized.
```json
{
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