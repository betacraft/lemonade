export GOPATH="$(dirname "$(dirname "$(dirname "$(dirname "$(pwd)")")")")"
export PATH=$PATH:$GOPATH/bin
# Setting up the env variables for the web app
export ENV=dev
export LEMN_MG_URI="mongodb://localhost:27017/lemonade"
export LEMN_MG_DB_NAME="lemonade"
# --
echo Cleaning up
go clean
echo Running go generate 
go generate github.com/rainingclouds/lemonade/models
echo Building Lemonade 
go vet 
go get 
echo Installing Lemonade
go install
echo Runnning Lemonade
lemonade