export GOPATH="$(dirname "$(dirname "$(dirname "$(dirname "$(pwd)")")")")"
export PATH=$PATH:$GOPATH/bin
echo Current GOPATH : $GOPATH
# go test github.com/rainingclouds/lemonades/mailer -v
go test github.com/rainingclouds/lemonades/models -v