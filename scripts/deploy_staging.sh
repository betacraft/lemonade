export GOPATH="$(dirname "$(dirname "$(dirname "$(dirname "$(pwd)")")")")"
export PATH=$PATH:$GOPATH/bin
echo Current GOPATH = $GOPATH
git push origin master
echo Pushing current master to heroku
git push heroku master