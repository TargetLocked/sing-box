### 结构

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
  "tolerance": 50,
  "idle_timeout": "",
  "interrupt_exist_connections": false
}
```

### 字段

#### outbounds

==必填==

用于测试的出站标签列表。

#### url

用于测试的链接。默认使用 `https://www.gstatic.com/generate_204`。

#### interval

测试间隔。 默认使用 `3m`。

#### tolerance

以毫秒为单位的测试容差。 默认使用 `50`。

#### idle_timeout

空闲超时。默认使用 `30m`。

#### interrupt_exist_connections

当选定的出站发生更改时，中断现有连接。

仅入站连接受此设置影响，内部连接将始终被中断。

#### max_successive_failures

发起额外测试前的最大连续失败次数。设置为`0`以禁用。

### 固定

可在运行时通过 Clash API（`PUT /proxies/{tag}`）或图形客户端选择出站来固定该出站。被固定的出站在通过 URL 测试时将代替最快的出站被使用，测试未通过时组将回退到最快的出站。选择时不可用的出站将被拒绝。

固定状态保存在缓存文件中，手动对该组进行 URL 测试或调用 `DELETE /proxies/{tag}` 将解除固定。在图形客户端中，被固定的组会在其类型旁显示 📌；当组因固定的出站不可用而回退时显示 📌⏳。
