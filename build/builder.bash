#!/bin/bash

cleanup() {
    echo "[+] Removing builder images..."
    docker image prune --filter label=stage=builder
}

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options"
    echo -e "\tFLAG    DESCRIPTION"
    echo -e "\t-p      Build peer image"
    echo -e "\t-r      Build registration service image"
    echo -e "\t-s      Build sequencer image"
}

DF_FOLDER="./dockerfiles/"
PEER_DF="${DF_FOLDER}peer/Dockerfile"
SEQ_DF="${DF_FOLDER}sequencer/Dockerfile"
REG_DF="${DF_FOLDER}registration/Dockerfile"
CONTEXT="../"

[ "$#" -eq 1 ] && [ "$1" = "-h" ] && usage && exit 0

while getopts ":psr" opt; do
    case ${opt} in
        p ) 
            echo "[+] Building comunigo/peer:latest image"
            docker build -f ${PEER_DF} -t comunigo/peer:latest ${CONTEXT}
            ;;
        s )
            echo "[+] Building comunigo/sequencer:latest image"
            docker build -f ${SEQ_DF} -t comunigo/sequencer:latest ${CONTEXT}
            ;;
        r )
            echo "[+] Building comunigo/registration:latest image"
            docker build -f ${REG_DF} -t comunigo/registration:latest ${CONTEXT}
            ;;
        ? )
            echo "[!] Invalid image requested"
            usage
            cleanup
            exit 1 
    esac
done
shift $((OPTIND -1))

cleanup
exit 0