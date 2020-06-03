# ws2redis

websocket 에서 redis 연결하는 interface server

## 사용법
sockjs로 연결 후 `<ip:port> <authKey>` 형식으로 최초 요구 후 REDIS command를 보낼 수 있다.

## Docker
* image를 생성한다.
    ```bash
    docker build . -t ws2redis
    ```
* container를 시작한다.
    ```bash
    docker run -d -p 8090:8090 --name ws2redis ws2redis
    ```