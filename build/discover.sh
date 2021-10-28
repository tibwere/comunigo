#!/bin/sh

[ $# -gt 1 ] && echo "Usage: $0 [test]" && exit 1
[ $# -eq 1 ] && [ "$1" = "-h" ] && echo "Usage: $0 [-t]" && exit 0
[ $# -eq 1 ] && [ "$1" != "-t" ] && echo "Usage: $0 [-t]" && exit 1

if [ "$1" = "-t" ]; then
    docker ps | grep comunigo/peer:latest | cut -d ":" -f 3 | cut -d "-" -f 1 | sort -u | tr "\n" "," | sed "s/,$/\n/"
else
    echo "List of active peers:"
    docker ps | grep comunigo/peer:latest | cut -d ":" -f 3 | cut -d "-" -f 1 | sort -u | xargs -L1 -I {} -n 1 echo -e "\t- http://localhost:{}/"
fi