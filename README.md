# go-auto-cfg&nbsp;&nbsp;&nbsp;&nbsp;[![Build Status](https://travis-ci.org/sidtharthanan/go-auto-cfg.svg?branch=master)](https://travis-ci.org/sidtharthanan/go-auto-cfg)

go-auto-cfg provides a way to access environmental variables as go functions.

go-auto-cfg is a wrapper over [viper](https://github.com/spf13/viper). It takes a schema file as input and generates a go file, the config loader, with functions named after the env variables.

It is available as a CLI tool. Also, can be used with [go generate](https://blog.golang.org/generate) tool. [Go generate example](#use-go-generate).

## Installation
```bash
go get -u github.com/sidtharthanan/go-auto-cfg
go-auto-cfg              --help  # depending on your PATH configuration use either of
$GOPATH/bin/go-auto-cfg  --help  # these commands
```
## Usage
```bash
go-auto-cfg <schema-file> <output-file>
```
### Example
```
go-auto-cfg schema.yml config/config.auto.go
```
The above command will parse `schema.yml` file and generate `config.auto.go` file under `config` package.
## Schema file
The schema is defined in a yaml file. It has to be in the following format.
```yaml
<ENV-variable-name>: <data-type>[,<constraint>]
```

### Schema example
```yaml
AUTH_FEATURE_ON: bool
CLIENT_ID: string,@optional()
CLIENT_PASS: string,@optional()
POOL_SIZE: integer,@optional(5)
ENVS: strings
```

### Data types
| Type    | Description |
|:--------|:------------|
| integer | To int. |
| float   | To float64. |
| string  | To string. |
| bool    | To bool. |
| strings | To []string. Value is either a list of strings or space separated values. [example](#sample-schemayml-file)|

### Constraints
| Constraint    | Description |
|:------------|:------------|
| @optional() | Makes the variable optional. By default all variables are required. |
| @optional(default-value)   | Makes the variable optional with a default value. |

## Full example
### Sample **schema.yml** file
```yaml
AUTH_FEATURE_ON: bool
AUTH_TIMEOUT_IN_HOURS: float
CLIENT_ID: integer
CLIENT_PASS: string
ENVS: strings
SOME_PREFIX: string,@optional(abc)
SOME_SUFFIX: string,@optional()
POOL_SIZE: integer,@optional(20)
```

### Configuration file **variables.yml**
```yaml
AUTH_FEATURE_ON: false
AUTH_TIMEOUT_IN_HOURS: 2.5

CLIENT_ID: 20082
CLIENT_PASS: qyUswmix82sw

ENVS: uat prod
```

### Generate the go config loader file
Alternatively this can be [automated](#use-go-generate).
```bash
go-auto-cfg schema.yml config/config.auto.go
```

### Use **config/config.auto.go** to load **variables.yml**
```go
 import cfg "config"
 
 func main() {
   //load config file "variables.yml" from "." directory
   cfg.Load("variables", ".")
   if cfg.AUTH_FEATURE_ON() {  // false
     middleware.authorize(user)
   }
   poolSize := cfg.POOL_SIZE() // 20
   prefix := cfg.SOME_PREFIX() // "abc"
   suffix := cfg.SOME_SUFFIX() // ""
 }
```

### Use go generate

```go
 //go:generate go-auto-cfg schema.yml config/config.auto.go
 import cfg "config"

 func main() {
   ...
 }
```

### Load multiple configs
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

## FAQ

1. **Q:** Why functions`cfg.SOME_CONFIG()` not simple struct members`cfg.SOME_CONFIG`?

   **A:** We want read only configuration and it is not possible to achieve without methods.
