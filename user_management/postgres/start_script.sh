#!/bin/bash/

cd common-files

until PGPASSWORD="restsim" psql -h eric-oss-porter-postgres -U restsim -d restsim_portal -c \\dt ; do sleep 2; echo "Retrying to connect to server" ; done

for i in *.sql; do

PGPASSWORD="restsim" psql -h eric-oss-porter-postgres -U restsim -d restsim_portal < $i

done

cd ..

cd "$1"

for i in *.sql; do

PGPASSWORD="restsim" psql -h eric-oss-porter-postgres -U restsim -d restsim_portal < $i

done
