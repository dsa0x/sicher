# Sicher

Sicher is a Go implementation of the secret management system that was introduced in Ruby on Rails 6.

Sicher is a go package that allows the secure storage of encrypted credentials in a version control system. The credentials can only be decrypted by a key file, and this key file is not added to the source control. The file is edited in a temp file on a local system and destroyed after each edit.

Using sicher in a project creates a set of files

- `environment.enc`
  - This is an encrypted file that stores the credentials. Since it is encrypted, it is safe to store these credentials in source control.
  - It it is encrypted using the [AES encryption](https://pkg.go.dev/crypto/aes) system.
- `environment.key`
  - This is the master key used to decrypt the credentials. This must not be committed to source control.

## Installation

To use sicher in your project, you need to install the go module as a library and also as a CLI tool.

Installing the library,

```shell
go get github.com/dsa0x/sicher
```

Installing the command line interface,:

```shell
go install github.com/dsa0x/sicher/cmd/sicher
```

## Usage

**_To initialize a new sicher project_**

```shell
sicher init
```

**_Optional flags:_**

| flag       | description                                                           | default | options        |
| ---------- | --------------------------------------------------------------------- | ------- | -------------- |
| -env       | set the environment name                                              | dev     |                |
| -path      | set the path to the credentials file                                  | .       |                |
| -style     | set the style of the decrypted credentials file                       | dotenv  | dotenv or yaml |
| -gitignore | path to the gitignore file. the key file will be added here, if given |         |                |

This will create a key file `{environment}.key` and an encrypted credentials file `{environment}.enc` in the current directory. The environment name is optional and defaults to `dev`, but can be set to anything else with the `-env` flag.

**_To edit the credentials:_**

```shell
sicher edit
```

OR

to use the key from environment variable:

```shell
env SICHER_MASTER_KEY=`{YOUR_KEY_HERE}` sicher edit
```

**_Optional flags:_**

| flag    | description                                     | default | options        |
| ------- | ----------------------------------------------- | ------- | -------------- |
| -env    | set the environment name                        | dev     |                |
| -path   | set the path to the credentials file            | .       |                |
| -editor | set the editor to use                           | vim     |                |
| -style  | set the style of the decrypted credentials file | dotenv  | dotenv or yaml |

This will create a temporary file, decrypt the credentials into it, and open it in your editor. The editor defaults to `vim`, but can be also set to other editors with the `-editor` flag. The temporary file is destroyed after each save, and the encrypted credentials file is updated with the new content.

Known good editors are:

- code
- emacs
- gvim
- mvim
- nano
- nvim
- subl
- vi
- vim
- vimr

Graphical editors require a flag to instruct the CLI to wait for the editor to exit. Additional graphical editors can be supported by adding the binary name and flag to the `waitFlagMap` in `sicher.go`. Most CLI editors should work out of the box, but your mileage may vary.

Then in your app, you can use the `sicher` library to load the credentials:

```go
package main
import (
	"fmt"

	"github.com/dsa0x/sicher/sicher"
)

type Config struct {
	Port        string `required:"true" env:"PORT"`
	MongoDbURI  string `required:"true" env:"MONGO_DB_URI"`
	MongoDbName string `required:"true" env:"MONGO_DB_NAME"`
	AppUrl   string `required:"false" env:"APP_URL"`
}

func main() {
	var config Config

	s := sicher.New("dev", ".")
	s.SetEnvStyle("yaml") // default is dotenv
	err := s.LoadEnv("", &cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

The `LoadEnv` function will load the credentials from the encrypted file `{environment.enc}`, decrypt it with the key file `{environment.key}` or the environment variable `SICHER_MASTER_KEY`, and then unmarshal the result into the given config object. The example above uses a `struct`, but the object can be of type `struct` or `map[string]string`.

**_LoadEnv Parameters:_**

| name   | description                             | type          |
| ------ | --------------------------------------- | ------------- |
| prefix | the prefix of the environment variables | string        |
| config | the config object                       | struct or map |

The key also be loaded from the environment variable `SICHER_MASTER_KEY`. In production, storing the key in the environment variable is recommended.

All env files should be in the format like the example below:

For `dotenv`:

```
PORT=8080
MONGO_DB_URI=mongodb://localhost:27017
MONGO_DB_NAME=sicher
APP_URL=http://localhost:8080
```

For `yaml`:

```
PORT:8080
MONGO_DB_URI:mongodb://localhost:27017
MONGO_DB_NAME:sicher
APP_URL:http://localhost:8080
```

If the object is a struct, the `env` tag must be attached to each variable. The `required` tag is optional, but if set to `true`, it will be used to check if the field is set. If the field is not set, an error will be returned.
An example of how the struct will look like:

```go
type Config struct {
	Port        string `required:"true" env:"PORT"`
	MongoDbURI  string `required:"true" env:"MONGO_DB_URI"`
	MongoDbName string `required:"true" env:"MONGO_DB_NAME"`
	AppUrl   string `required:"false" env:"APP_URL"`
}
```

If object is a map, the keys are the environment variables and the values are the values.

### Note

- Not tested with Windows.

### Todo or not todo

- Add a `-force` flag to `sicher init` to overwrite the encrypted file if it already exists
- Enable support for nested yaml env files
- Add support for other types of encryption
- Test on windows

### License

MIT License
