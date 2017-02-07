# overseer

## Monitors

Monitor configurations live at `/etc/overseer/monitors` and each file must have a `.json` file extension.

In addition to the monitor configurations below, each monitor configuration can also name specific notifiers and loggers that it should use such as:

```json
{
    ...,
    "notifiers": [
        "Slack status channel",
        "email ops team"
    ],
    "loggers": [
        "mongodb"
    ]
}
```

The values for the `notifiers` key can either be the notifiers `name` or `type`.

### http-status

```json
{
    "type": "http-status",
    "name": "localhost status",
    "url": "http://localhost/",
    "codes": [
        200,
        400
    ],
    "check_interval": "10s",
    "notification_interval": "30m",
    "timeout": "2s",
    "verify": false,
    "method": "HEAD",
    "retries": 3
}
```

### http-content

```json
{
    "type": "http-content",
    "name": "localhost content",
    "url": "http://localhost/",
    "content": "Welcome!",
    "check_interval": "10s",
    "notification_interval": "30m",
    "timeout": "2s",
    "verify": false,
    "method": "HEAD",
    "retries": 3
}
```

### connect

```json
{
    "type": "connect",
    "name": "localhost port 80",
    "host": "localhost",
    "protocol": "tcp",
    "port": 80,
    "check_interval": "10s",
    "notification_interval": "30m",
    "timeout": "2s",
    "retries": 3
}
```

## Notifiers

Notifier configurations live at `/etc/overseer/notifiers` and each file must have a `.json` file extension.

### stderr

```json
{
    "type": "stderr",
    "name": "stderr notifier"
}
```

### slack

```json
{
    "type": "slack",
    "name": "Slack status channel",
    "webhook_url": "https://hooks.slack.com/services/stuff/goes/here",
    "channel": "#status",
    "username": "overseer"
}
```

### mailgun

```json
{
    "type": "mailgun",
    "name": "email ops team",
    "domain": "overseer.mailgun.org",
    "apikey": "key-goes-here",
    "from": "Overseer <overseer@overseer.mailgun.org>",
    "to": [
        "oncall@someplace.com",
        "ops@someplace.com",
        "bob@someplace.com"
    ]
}
```

## Loggers

Logger configurations live at `/etc/overseer/loggers` and each file must have a `.json` file extension.

Loggers differ from notifiers in that they log every status check result, as opposed to just results
from status changing state.  This is useful for creating historical response time graphs.

### stderr

```json
{
    "type": "stderr",
    "name": "stderr logger"
}
```

### mongodb

```json
{
    "type": "mongodb",
    "name": "mongodb logger",
    "mongodb_uri": "mongodb://localhost/overseer"
}
```

### elasticsearch

```json
{
    "type": "elasticsearch",
    "name": "elasticsearch logger",
    "servers": [
        "http://127.0.0.1:9200"
    ],
    "username": "elastic_user",
    "password": "elastic_pass",
    "index": "overseer",
    "doc_type": "log"
}
```
