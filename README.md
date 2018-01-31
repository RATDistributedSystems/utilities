# Utilities
Utilities that some might find useful for the project

# Config file 

Given a JSON file called `config.json`

```
{
    "serverConfigs" : {
        "webserver" : {
            "address": "0.0.0.1",
            "port": "44440",
            "protocol": "http"
        },

        "transaction" : {
            "address": "localhost",
            "port": "44441",
            "protocol": "tcp"
        }
    },

    "misc" : {
        "htmlLoc" : "../frontend/"
    } 
}
```

It can be imported like this

```
import (
    "github.com/RATDistributedSystems/utilities"
)

func main() {
    config := utilities.LoadConfigs("config.json")

    //to get transaction server details
    addr, protocol := config.GetServerDetails("transaction")

    // get something from the misc map
    loc := config.GetValue("htmlLoc")
}
```
