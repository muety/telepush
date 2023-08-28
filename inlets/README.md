# Inlets
Inlets are endpoints that accept messages in different formats. Most of the time, when sending messages, you will use the `default` inlet. However, Telepush also supports to accept webhook payloads from third-party applications, such as Prometheus Alertmanager, Bitbucket, and others. The according inlets are responsible for receiving messages in the respective formats and translating them into Telegram text messages.

Telepush comes with a couple of pre-configured inlets. However, you can easily **define your own**.

## Creating custom inlets
Inlets can be written in code (Go) (e.g. see the [`default`](default)) inlet or in a config-based fashion with YAML (recommended). Inlet configs are placed inside the [inlets.d](../inlets.d) folder and require a couple of different properties, such as a name and a template. See [`example.yaml`](../inlets.d) for an easy to understand example inlet definition. 

For each inlet, you will write a [Go template](https://pkg.go.dev/text/template) to define the resulting Telegram message's text. Inside the template, you'll have access to `.Message`, containing the incoming requests' body payload (as plain text or a nested map in case of JSON content).

Whenever adding or updating an inlet, Telepush will automatically reload its config.

---

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
Not yet implemented.

## `plain`
`/api/inlets/plain/<recipient>`

Forwards a plain string (`text/plain`) of text to a Telegram chat.

#### Body (Text Message)
```
<string>
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

## `grafana`
`/api/inlets/grafana/<recipient>`

Accepts, transforms and forwards alerts sent by [Grafana](https://grafana.com/docs/grafana/latest/alerting/) to a Telegram chat.

Create a new contact point with `POST` method and URL `https://telepush.dev/api/inlets/grafana/<recipient>`. Also see [webhook-notifier](https://grafana.com/docs/grafana/latest/alerting/contact-points/notifiers/webhook-notifier/).

## `bitbucket`
`/api/inlets/bitbucket/<recipient>`

Accepts, transforms and forwards events sent by [Bitbucket](https://bitbucket.org/) to a Telegram chat.

#### Limitations
* Currently, only these events are implemented:
  * `repo:commit_status_created`
  * `repo:commit_status_updated`

#### Parameters
Requires the `X-Event-Key` header to be set. 

### Body
See [Events](https://confluence.atlassian.com/bitbucket/event-payloads-740262817.html).

## `stripe`
`/api/inlets/stripe/<recipient>`

Accepts, transforms and forwards events sent by [Stripe](https://stripe.com/docs/webhooks) to a Telegram chat.

#### Limitations
* Currently, only these events are implemented:
    * `customer.subscription.created`
    * `customer.subscription.updated`

### Body
See [Event](https://stripe.com/docs/api/events).

## `webmentionio`
`/api/inlets/webmentionio/<recipient>`

Accepts, transforms and forwards notifications sent by [Webmention.io](https://webmention.io) to a Telegram chat.

### Body
An example payload looks as follows, however, only `source` and `target` are utilized.
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