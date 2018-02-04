# go-auto-cfg
Access environmental variables as go-funcs. It is a wrapper over viper. It takes a schema file as input and generates a go file with functions named after the env variables.

[![Build Status](https://travis-ci.org/sidtharthanan/go-auto-cfg.svg?branch=master)](https://travis-ci.org/sidtharthanan/go-auto-cfg)

Sample **schema.yml** file:
```yaml
AUTH_FEATURE_ON: bool
AUTH_TIMEOUT_IN_HOURS: float

CLIENT_ID: integer
CLIENT_PASS: string

DB_URL: string
LOG_LEVEL: string
OTHER_FEATURES: strings
```

Run the following command to generate the go config file:
```bash
#go run $GOPATH/src/github.com/sidtharthanan/go-auto-cfg/cfg.go parse <schema-file> <package-name> <output-dir>
go run $GOPATH/src/github.com/sidtharthanan/go-auto-cfg/cfg.go parse schema.yml cfg config
```
The above command will parse `schema.yml` file and generate `config.auto.go` file, located at `config` directory in `cfg` go package.

Generated **config/config.auto.go** file:
```go
package cfg //package name is given as cli argument

...
//name of the config file without extension.
//paths for Viper to search for the config file in.
func ReadInCfg(name string, paths ...string) { ... }

func GetAUTH_FEATURE_ON() bool { ... }
func GetAUTH_TIMEOUT_IN_HOURS() float64 { ... }

func GetCLIENT_ID() int { ... }
func GetCLIENT_PASS() string { ... }

func GetDB_URL() string { ... }
func GetLOG_LEVEL() string { ... }
func GetOTHER_FEATURES() []string { ... }
```

Configuration file **variables.yml**:
```yaml
AUTH_FEATURE_ON: false
AUTH_TIMEOUT_IN_HOURS: 2.5

CLIENT_ID: 20082
CLIENT_PASS: qyUswmix82sw

DB_URL: postgres://postgres:password@localhost:5432/db
LOG_LEVEL: DEBUG
OTHER_FEATURES: OTP,SMS
```

Use **config/config.auto.go** to load **variables.yml**:
```go
 import "cfg"
 
 func main() {
   //load config file "variables.yml" from "." directory
   cfg.ReadInCfg("variables", ".")
   if cfg.GetAUTH_FEATURE_ON() {
     middleware.authorize(user) // if the toggle is on use this feature
   }
 }
```

Generation of config file reader could be automated as follows:

```go
 //go:generate go run $GOPATH/src/github.com/sidtharthanan/go-auto-cfg/cfg.go parse schema.yml cfg config
 import "cfg"

 func main() {
   ...
 }
```
