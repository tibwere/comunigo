#!/bin/bash

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options"
    echo -e "\tFLAG   VALUES"
    echo -e "\t-t     [sequencer | scalar | vectorial]   Algorithm to test"
    echo -e "\t-m     [single | multiple]                Test sending a single message or multiple simultaneously"
}

MOD=""
TOS=""

[ "$#" -eq 1 ] && [ "$1" = "-h" ] && usage && exit 0

# Parse command line options
while getopts ":t:m:" opt; do
    case ${opt} in
        t )
            TOS=${OPTARG}
            ;;
        m )
            MOD=${OPTARG}
            ;;
        ? )
            usage
            exit 1 
    esac
done
shift $((OPTIND -1))

if [ "${TOS}" != "sequencer" ] && [ "${TOS}" != "scalar" ] && [ "${TOS}" != "vectorial" ]; then
    echo "[!] Select a valid type of service: sequencer | scalar | vectorial"
    exit 1
fi

if [ "${MOD}" != "single" ] && [ "${MOD}" != "multiple" ]; then
    echo "[!] Select a valid modality for test: single | multiple"
    exit 1
fi

export COMUNIGO_TEST_PORTS=$(sh discover.sh -t)
go test -v ../integration-tests/ -run Test${MOD^}Send${TOS^}


