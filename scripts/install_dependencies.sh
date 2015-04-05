export GOPATH="$(dirname "$(dirname "$(dirname "$(dirname "$(pwd)")")")")"
export PATH=$PATH:$GOPATH/bin
echo "Current go path : $GOPATH"
echo "Adding bone-mux"
go get github.com/go-zoo/bone
echo "Adding godep"
go get github.com/tools/godep
echo "Adding mgo"
go get gopkg.in/mgo.v2
echo "Adding logrus"
go get github.com/sirupsen/logrus
echo "Adding pq"
go get github.com/lib/pq
echo "Adding gorm"
go get github.com/jinzhu/gorm
echo "Installing easytags"
go get github.com/rainingclouds/easytags
go install github.com/rainingclouds/easytags
echo "Adding gouuid"
go get "github.com/nu7hatch/gouuid"
echo "Adding binding"
go get github.com/rainingclouds/binding
echo "Adding Go Metrics"
go get github.com/rcrowley/go-metrics
echo "Installing autobindings"
go get github.com/rainingclouds/autobindings
go install github.com/rainingclouds/autobindings
echo "Adding s3"
go get github.com/AdRoll/goamz/s3
echo "Adding aws"
go get github.com/AdRoll/goamz/aws
echo "Adding logrus sentry hook"
go get github.com/sirupsen/logrus/hooks/sentry
echo "Adding mailer"
go get github.com/jordan-wright/email