```bash
require go version 1.2x ++
```

```bash
# setting
# vscode:Cmd+Shift+P
"[go]": {
    "editor.formatOnSave": true
},
"go.formatTool": "gofmt"

```

```bash
#lib
go mod init fiber-mongo-api
go get -u github.com/gofiber/fiber/v2 go.mongodb.org/mongo-driver/mongo github.com/joho/godotenv
go get github.com/klauspost/compress
```

```bash
# Live Reload
go install github.com/cosmtrek/air@latest
air init
```

```bash
# RUN
air
# if you need to run a command without live reload ::>> go run main.go
```

```bash
#unit test
go get github.com/stretchr/testify
go get github.com/stretchr/testify/assert
```

```bash
# clean cache
go clean -testcache # or restart IDE
```
