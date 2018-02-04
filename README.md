# go-auto-cfg
Access environmental variables as go-funcs. It is a wrapper over viper. It takes a schema file as input and generates a go file with functions named after the env variables.

[![Build Status](https://travis-ci.org/sidtharthanan/go-auto-cfg.svg?branch=master)](https://travis-ci.org/sidtharthanan/go-auto-cfg)

Sample schema file:
```yaml
APP_HOST: string
APP_PORT: integer

AUTH_FEATURE_ON: bool
AUTH_TIMEOUT_IN_HOURS: float

CLIENT_ID: integer
CLIENT_PASS: string

DB_URL: string
LOG_LEVEL: string
OTHER_FEATURES: strings
```

Sample config reader file generated:
```go
...
func GetAPP_HOST() string { ... }
func GetAPP_PORT() int { ... }

func GetAUTH_FEATURE_ON() bool { ... }
func GetAUTH_TIMEOUT_IN_HOURS() float64 { ... }

func GetCLIENT_ID() int { ... }
func GetCLIENT_PASS() string { ... }

func GetDB_URL() string { ... }
func GetLOG_LEVEL() string { ... }
func GetOTHER_FEATURES() []string { ... }
```

Sample configuration file:
```yaml
APP_HOST: localhost
APP_PORT: 8050

AUTH_FEATURE_ON: false
AUTH_TIMEOUT_IN_HOURS: 2.5

CLIENT_ID: 20082
CLIENT_PASS: qyUswmix82sw

DB_URL: postgres://postgres:password@localhost:5432/db
LOG_LEVEL: DEBUG
OTHER_FEATURES: OTP,SMS
```

Now you could easily do:
```go
 import "config"
 ...
 if config.GetAUTH_FEATURE_ON() { // false
   middleware.authorize(user)
 }
```

Run the following command to generate the go config file:
```bash
go run $GOPATH/src/github.com/sidtharthanan/go-auto-cfg/cfg.go parse schema.yml config configuration
```
The above command will parse `schema.yml` file and generates `config.auto.go` file.
This file will be located at `configuration` directory under `config` package.

To automate the code generation part, you could use `go generate` tool.
