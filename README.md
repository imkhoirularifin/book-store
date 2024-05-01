# Book Store

Personal project for learning purpose

## Getting Started

1. To generate swagger docs, Install swag by using:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Run `swag init --parseDependency --parseInternal` in the project's root folder which contains the `main.go` file. This will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).

```sh
swag init --parseDependency --parseInternal
```

3. To generate JWT Private and Public Key, use the following command:

```sh
ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key
# Don't add passphrase
openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub
```

4. Copy Private Key and Public Key, and encode with base64 in [here](https://www.base64encode.org/)

## Environment

Environment variables:

| Name            | Description                    | Default Value                |
| --------------- | ------------------------------ | ---------------------------- |
| HOST            | Hostname                       | localhost                    |
| PORT            | Port                           | 8080                         |
| IS_DEVELOPMENT  | Is Development                 | true                         |
| PROXY_HEADER    | Proxy Header                   | X-Real-IP                    |
| LOG_FIELDS      | Log Fields                     | level, time, logger, message |
| DB_DRIVER       | Database Driver                | sqlite                       |
| DB_DSN          | Database DSN                   | file::memory:?cache=shared   |
| JWT_PRIVATE_KEY | Base64 Encoded JWT Private Key |                              |
| JWT_PUBLIC_KEY  | Base64 Encoded JWT Public Key  |                              |
| JWT_EXPIRES_IN  | JWT Expires In                 | 24h                          |

## Run Command

```sh
go run main.go
```
