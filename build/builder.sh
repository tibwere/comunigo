#!/bin/sh

cleanup() {
    echo "[+] Removing builder images..."
    docker image prune --filter label=stage=builder
}

DF_FOLDER="./dockerfiles/"
PEER_DF="${DF_FOLDER}peer/Dockerfile"
SEQ_DF="${DF_FOLDER}sequencer/Dockerfile"
REG_DF="${DF_FOLDER}registration/Dockerfile"
CONTEXT="../"

[ "$#" -eq 1 ] && [ "$1" = "-h" ] && echo "[?] Usage: $0 [ peer | sequencer | registration ]+" && exit 0

for image in $@; do
    case ${image} in
        peer )
            echo "[+] Building comunigo/peer:latest image"
            docker build -f ${PEER_DF} -t comunigo/peer:latest ${CONTEXT}
            ;;
        sequencer )
            echo "[+] Building comunigo/sequencer:latest image"
            docker build -f ${SEQ_DF} -t comunigo/sequencer:latest ${CONTEXT}
            ;;
        registration)
            echo "[+] Building comunigo/registration:latest image"
            docker build -f ${REG_DF} -t comunigo/registration:latest ${CONTEXT}
            ;;
        * )
            echo "[!] Invalid image requested"
            cleanup
            exit 1
            ;;
    esac
done

cleanup
exit 0