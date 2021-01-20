#! /bin/bash

OUTPUT="socks_proxy"

while getopts :o: OPTION
do
    case ${OPTION} in
        o)
            if [ -z "${OPTARG}" ]
            then
                echo -e "Param -o OUTPUT error. OUTPUT file name is empty";
                exit 1;
            fi
            OUTPUT="${OPTARG}"
            ;;
        ?)
            echo -e "WRONG USE WAY!\n";
            exit 1;
            ;;
    esac
done

export GO111MODULE=on

VERSION=`git describe --all`
go build -o ${OUTPUT} -ldflags "-X 'main.Version=${VERSION}' -X 'main.BuildTime=`date`' -X 'main.GoVersion=`go version`' -X 'main.Build=Release'" main.go
if [ $? -ne 0 ]
then
    exit 1
fi
