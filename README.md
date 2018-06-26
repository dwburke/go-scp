# go-scp

example:

```
package main

import (
    "os"

    "github.com/dwburke/go-scp"
    "github.com/dwburke/go-tools"
)

func main() {
    private_key, err := scp.LoadPrivateKey("") // loads $HOME/.ssh/id_rsa by default
    tools.FatalError(err)

    scp_client, err := scp.New(private_key, "127.0.0.1", 22, os.GetEnv("USER"))
    tools.FatalError(err)
    defer scp_client.Close()

    err = scp_client.Get("/file/test_file.txt", "/tmp/test_file.txt")
    tools.FatalError(err)
}
```
