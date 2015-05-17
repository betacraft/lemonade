export GOPATH="$(dirname "$(dirname "$(dirname "$(dirname "$(pwd)")")")")"
export PATH=$PATH:$GOPATH/bin
echo Current GOPATH : $GOPATH
export ENV=dev
export LEMN_MG_URI="mongodb://localhost:27017/lemonade2"
export LEMN_MG_DB_NAME="lemonade2"
export AWS_ACCESS="AKIAINECWOX2MEE4UOSA"
export AWS_SECRET="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"

go test github.com/rainingclouds/lemonades/parsers -v
# go test github.com/rainingclouds/lemonades/mailer -v
#go test github.com/rainingclouds/lemonades/models -v