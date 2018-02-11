# go-auto-cfg
Access environmental variables as go-funcs. It is a wrapper over viper. It takes a schema file as input and generates a go file with functions named after the env variables.

[![Build Status](https://travis-ci.org/sidtharthanan/go-auto-cfg.svg?branch=master)](https://travis-ci.org/sidtharthanan/go-auto-cfg)

Installation:
```bash
go get -u github.com/sidtharthanan/go-auto-cfg
go-auto-cfg              --help  # depending on your PATH configuration use either of
$GOPATH/bin/go-auto-cfg  --help  # these commands
```

Sample **schema.yml** file:
```yaml
AUTH_FEATURE_ON: bool
AUTH_TIMEOUT_IN_HOURS: float
CLIENT_ID: integer
CLIENT_PASS: string
ENVS: strings
```

Run the following command to generate the go config file:
```bash
#go-auto-cfg <schema-file> <output-file>
go-auto-cfg schema.yml config/config.auto.go
```
The above command will parse `schema.yml` file and generate `config.auto.go` file, under `config` package.

Generated **config/config.auto.go** file:
```go
package config

...
//name of the config file without extension.
//paths for Viper to search for the config file in.
func Load(name string, paths ...string) { ... }

func AUTH_FEATURE_ON() bool { ... }

func AUTH_TIMEOUT_IN_HOURS() float64 { ... }

func CLIENT_ID() int { ... }

func CLIENT_PASS() string { ... }

func ENVS() []string { ... }
```
Functions are generated with `upper case name`s regardless of their case in `schema.yml`.

Configuration file **variables.yml**:
```yaml
AUTH_FEATURE_ON: false
AUTH_TIMEOUT_IN_HOURS: 2.5

CLIENT_ID: 20082
CLIENT_PASS: qyUswmix82sw

ENVS: uat prod
```

Use **config/config.auto.go** to load **variables.yml**:
```go
 import cfg "config"
 
 func main() {
   //load config file "variables.yml" from "." directory
   cfg.Load("variables", ".")
   if cfg.AUTH_FEATURE_ON() {
     middleware.authorize(user) // if the toggle is on use this feature
   }
 }
```

Generation of config file reader could be automated as follows:

```go
 //go:generate go-auto-cfg schema.yml config/config.auto.go
 import cfg "config"

 func main() {
   ...
 }
```

Multiple instances of configuration can be loaded as follows:
```go
 import cfg "config"
 
 func main() {
   cfg.Load("global", ".")
   cfg.SOME_GLOBAL_CONFIG()

   cfgA := cfg.New()
   cfgA.Load("moduleA", ".")
   cfgA.MODULE_A_SPECIFIC_CONFIG()
   
   cfgB := cfg.New()
   cfgB.Load("moduleB", ".")
   cfgB.MODULE_B_SPECIFIC_CONFIG()

 }
```

Optional config:
The following would default to 100, If not configured.
```yaml
AUTH_FEATURE_ON: integer,@optional(100)
```
The following would default to 0, golang zero value, If not configured.
```yaml
AUTH_FEATURE_ON: integer,@optional()
```

FAQ:
1. **Q:** Why functions`cfg.SOME_CONFIG()` not simple struct members`cfg.SOME_CONFIG`?

   **A:** We want read only configuration and it is not possible to achieve without methods.
