export GOPATH="$(dirname "$(dirname "$(dirname "$(dirname "$(pwd)")")")")"
export PATH=$PATH:$GOPATH/bin
# --
echo Current GOPATH : $GOPATH
echo Cleaning up
go clean
echo Running go generate 
go generate github.com/rainingclouds/lemonades/models
echo Building lemonade 
go vet 
go get 