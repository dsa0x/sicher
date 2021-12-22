# Sicher

Sicher is a port of Rails secret management system.

Sicher is a go module that allows storage of encrypted credentials right in VCS. The can only be decrypted by a key file which is not added to the source control.
The file is edited in a temp file and destroyed after each edit.

Using sicher in a project creates a set of files

- `environment.enc`
  - This is an encrypted file that stores the credentials. Since it is encrypted, it is safe to store these credentials in source control.
- `environment.key`
  - This is the master key used to decrypt the encrypted credentials. This must not be committed to source control. To prevent unintentional commit, the file is appended to the `.gitignore` file.

## Installation

To use as a library,

```shell
go get github.com/dsaOx/sicher
```

To use in the command line, install the binary:

```shell
go install github.com/dsaOx/sicher/cmd/sicher
```

## Usage

To initialize a new sicher project,

```shell
sicher init --env environmentname
```

This will create a key file `{environment}.key` and an encrypted credentials file `{environment}.enc` in the current directory. The environment name is optional and defaults to `dev`.

To edit the credentials,

```shell
sicher edit --env environmentname
```

This will create a temporary file, decrypt the credentials into it, and open it in your editor. The editor defaults to `vim`, but can be also set to `nano` or `vi` with the `--editor` variable. The temporary file is destroyed after each edit, and the encrypted credentials file is updated with the new content.

Then in your app, you can use the `Sicher` struct to access the credentials:

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

The `LoadEnv` function will load the credentials from the encrypted file `environment.enc` and decrypt it with the key file `environment.key`, and unmarshal the result into the given struct.

An example of the temporary env file:

```
PORT=8080
MONGO_DB_URI=mongodb://localhost:27017
MONGO_DB_NAME=sicher
APP_URL=http://localhost:8080
```
