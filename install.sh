#!/bin/sh

GIT_CHECK=`which git`
if [ $GIT_CHECK = "" ] 
then
    apt-get install -y git
fi

GCC_CHECK=`which gcc`
if [ $GCC_CHECK = "" ] 
then
    apt-get install -y gcc
fi

GO_CHECK=`which go`
if [ $GO_CHECK = "" ] 
then
    apt-get install -y golang
fi

GTK_CHECK=`dpkg-query --list | grep libgtk-3-dev*`
if [ $GO_CHECK = "" ] 
then
    apt-get install -y libgtk-3-dev
fi

MVN_CHECK=`which mvn`
if [ $MVN_CHECK = "" ] 
then
    apt-get install -y maven
fi

rm -rf gopos_golibs
GOPOS_ROOT=$PWD
export GOPATH=$PWD/gopos_golibs
go get github.com/go-sql-driver/mysql

go get -d github.com/muller95/gotk3
cd gopos_golibs/src/github.com/
mv muller95 gotk3
cd $GOPOS_ROOT

cd server 
go build

cd ../client
go build -tags gtk_3_16

cd ../monitor
go build -tags gtk_3_16

cd ../printservice/my-app/
mvn clean compile
mvn package assembly:single
