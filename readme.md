

To start server:


```bash

git clone https://github.com/ORESoftware/mux-problem-demo.git 
cd mux-problem-demo

export GOPATH="$PWD"

go get	"github.com/dgrijalva/jwt-go"
go get	"github.com/gorilla/context"
go get	"github.com/gorilla/mux"
go get	"github.com/mitchellh/mapstructure"
go get	"github.com/sirupsen/logrus"

export huru_api_port="3000"

go clean
go install huru
"$GOPATH/bin/huru"

```

