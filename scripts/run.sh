export GOPATH="$(dirname "$(dirname "$(dirname "$(dirname "$(pwd)")")")")"
export PATH=$PATH:$GOPATH/bin
# Setting up the env variables for the web app
export ENV=dev
export LEMN_MG_URI="mongodb://localhost:27017/lemonade2"
export LEMN_MG_DB_NAME="lemonade2"
# --
echo Cleaning up
go clean
echo Running go generate 
go generate github.com/rainingclouds/lemonades/models
echo Building Lemonades 
go vet 
go get 
echo Installing Lemonades
go install
echo Runnning Lemonades
lemonades