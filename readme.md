# CONTCP

这是一个简单的 Go 应用程序，用于检查服务器状态并返回结果。

## 使用

- 返回 JSON 格式：
  ```
  http://localhost:3344/{server}
  ```
- 直接返回状态 true 或 false：
  `    http://localhost:3344/status/{server}
   `
  其中 `{server}` 不加端口为 imcp ping 加端口为 tcp ping
