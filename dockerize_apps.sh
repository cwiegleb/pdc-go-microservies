#!/bin/bash

pushd pdc-cashbox-service
 docker build -t cwiegleb/pdc-cashbox-service:latest .
popd

pushd pdc-dealer-service
 docker build -t cwiegleb/pdc-dealer-service:latest .
popd

pushd pdc-order-service
 docker build -t cwiegleb/pdc-order-service:latest .
popd 

pushd pdc-csv-upload-service
 docker build -t cwiegleb/pdc-csv-upload-service:latest .
popd 

pushd pdc-db
 docker build -t cwiegleb/pdc-db:latest .
popd


