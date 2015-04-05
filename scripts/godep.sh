export GOPATH="$(dirname "$(dirname "$(dirname "$(dirname "$(pwd)")")")")"
echo $GOPATH
$GOPATH/bin/godep save