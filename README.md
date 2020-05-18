# ws2redis

websocket 에서 redis 연결하는 interface server

## 사용법
sockjs로 연결 후 `REQCON <ip:port> <authKey>` 형식으로 최초 요구 후 REDIS command를 보낼 수 있다.