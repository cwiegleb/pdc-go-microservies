#!/bin/bash

pushd pdc-cashbox-service
 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pdc-cashbox-service
popd

pushd pdc-dealer-service
 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pdc-dealer-service
popd

pushd pdc-csv-upload-service
 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pdc-csv-upload-service
popd

pushd pdc-order-service
 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pdc-order-service
popd

pushd pdc-db
 CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pdc-db
popd
 

