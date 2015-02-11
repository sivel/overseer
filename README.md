# overseer

## Monitors

Monitor configurations live at `/etc/overseer/monitors` and each file must have a `.json` file extension.

In addition to the monitor configurations below, each monitor configuration can also name specific notifiers that it should use such as:

```json
{
    ...,
    "notifiers": [
        "Slack status channel",
        "email ops team"
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
    "method": "HEAD"
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
    "method": "HEAD"
}
```

### connect

```json
{
    "type": "connect",
    "name": "localhost port 80",
    "host": "localhost",
    "port": 80,
    "check_interval": "10s",
    "notification_interval": "30m",
    "timeout": "2s"
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
    "token": "oauth2_token_here",
    "channel": "status",
    "username": "overseer"
}
```

If you know your channel ID, you can replace `channel` with `channel_id` and use the channel ID instead of a channel name.

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
