# CONTCP

这是一个简单的 Go 应用程序，用于检查服务器状态并返回结果。

## 使用

    - 返回JSON格式：
        ```
        http://localhost:3344/{server}
        ```
    - 直接返回状态true或false：
        ```
        http://localhost:3344/status/{server}
        ```
    其中 `{server}` 不加端口为imcp ping 加端口为tcp ping
