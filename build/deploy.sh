#! /bin/bash

RunMode="online"
ServiceName="socks_proxy"

KillName()
{
        local l_name="$1"
        ps ux | grep ${l_name} | grep -v grep | awk '{print $2}' | xargs kill
}

while getopts :r: OPTION
do
    case ${OPTION} in
        r)
            if [ "${OPTARG}" == "test" ]
            then
                RunMode="test"
            fi
            ;;
        ?)
            echo -e "WRONG USE WAY!\n";
            exit 1;
            ;;
    esac
done

git pull

if [ "${RunMode}" == "online" ]
then
    mv ${ServiceName} ${ServiceName}.bak
    sh -x ./build/build.sh
    mkdir -p log/
    KillName ${ServiceName}
    nohup ./${ServiceName} &
fi

exit 0
