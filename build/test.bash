#!/bin/bash

usage() {
    echo "[?] Usage: $0 [ peer | sequencer | registration ]+"
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

export COMUNIGO_TEST_PORTS=$(comunigo-peer-discovery -t)
go test -v ../integration-tests/ -run Test${MOD^}Send${TOS^}

