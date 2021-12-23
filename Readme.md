# Sicher

Sicher is a Go implementation of the secret management system that was introduced in Ruby on Rails 6.

Sicher is a go module that allows safe storage of encrypted credentials in a version control system. The credentials can only be decrypted by a key file which is not added to the source control. The file is edited in a temp file on a local system and destroyed after each edit.

Using sicher in a project creates a set of files

- `environment.enc`
  - This is an encrypted file that stores the credentials. Since it is encrypted, it is safe to store these credentials in source control.
  - It it is encrypted using the [AES encryption](https://pkg.go.dev/crypto/aes) system.
- `environment.key`
  - This is the master key used to decrypt the credentials. This must not be committed to source control. To prevent unintentional commit, the file is appended to the `.gitignore` file.

## Installation

To use sicher in your project, you need to install the go module as a library and also install the CLI.

Installing the library,

```shell
go get github.com/dsaOx/sicher
```

Installing the command line interface,:

```shell
go install github.com/dsaOx/sicher/cmd/sicher
```

## Usage

To initialize a new sicher project, run

```shell
sicher init
```

This will create a key file `{environment}.key` and an encrypted credentials file `{environment}.enc` in the current directory. The environment name is optional and defaults to `dev`, but can be set to anything else with the `-env` flag.

To edit the credentials, run

```shell
sicher edit
```

This will create a temporary file, decrypt the credentials into it, and open it in your editor. The editor defaults to `vim`, but can be also set to `nano` or `vi` with the `-editor` flag. The temporary file is destroyed after each save, and the encrypted credentials file is updated with the new content.

Then in your app, you can use the `sicher` module to access the credentials:

```go
package main
import (
	"fmt"

	"github.com/dsa0x/sicher/sicher"
)

type Config struct {
	Port        string `required:"true" envconfig:"PORT"`
	MongoDbURI  string `required:"true" envconfig:"MONGO_DB_URI"`
	MongoDbName string `required:"true" envconfig:"MONGO_DB_NAME"`
	AppUrl   string `required:"false" envconfig:"APP_URL"`
}

func main() {
	var config Config

	s := sicher.New("prod")
	err := s.LoadEnv("", &config)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

The `LoadEnv` function will load the credentials from the encrypted file `{environment.enc}`, decrypt it with the key file `{environment.key}`, and then unmarshal the result into the given config object. The example above uses a `struct`, but the object can be of type `struct` or `map[string]string`.

All env files should be in the format like the example below:

```
PORT=8080
MONGO_DB_URI=mongodb://localhost:27017
MONGO_DB_NAME=sicher
APP_URL=http://localhost:8080
```

### Todo

- Make addition to the `.gitignore` file optional
- When a user wants to initialize sicher, add a warning if an encrypted file already exists, but there is no key file
- Add a `-force` flag to `sicher init` to overwrite the encrypted file if it already exists
- Enable support for yaml style env files
- Add support for other types of encryption
