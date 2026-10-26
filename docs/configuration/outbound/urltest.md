### Structure

```json
{
  "type": "urltest",
  "tag": "auto",
  
  "outbounds": [
    "proxy-a",
    "proxy-b",
    "proxy-c"
  ],
  "url": "",
  "interval": "",
  "tolerance": 0,
  "idle_timeout": "",
  "interrupt_exist_connections": false,
  "max_successive_failures": 0
}
```

### Fields

#### outbounds

==Required==

List of outbound tags to test.

#### url

The URL to test. `https://www.gstatic.com/generate_204` will be used if empty.

#### interval

The test interval. `3m` will be used if empty.

#### tolerance

The test tolerance in milliseconds. `50` will be used if empty.

#### idle_timeout

The idle timeout. `30m` will be used if empty.

#### interrupt_exist_connections

Interrupt existing connections when the selected outbound has changed.

Only inbound connections are affected by this setting, internal connections will always be interrupted.

#### max_successive_failures

Max successive dial failures before an extra test is enforced. Set to `0` to disable.

### Pinning

An outbound can be pinned at runtime by selecting it, through the Clash API (`PUT /proxies/{tag}`) or in the graphical clients. A pinned outbound is used in place of the fastest one for as long as it passes URL tests, and the group falls back to the fastest outbound while it doesn't. An outbound that is unavailable when selected is refused.

The pin is kept in the cache file, and is removed by a manual URL test of the group or by `DELETE /proxies/{tag}`. Pinned groups are shown with a 📌 next to their type in the graphical clients, or 📌⏳ while the group has fallen back from an unavailable pinned outbound.
