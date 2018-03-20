#!/bin/bash

pushd pdc-cashbox-service
 docker build -t cwiegleb/pdc-cashbox-service:no-article .
 docker push cwiegleb/pdc-cashbox-service:no-article
popd

pushd pdc-dealer-service
 docker build -t cwiegleb/pdc-dealer-service:no-article .
 docker push cwiegleb/pdc-dealer-service:no-article
popd

pushd pdc-order-service
 docker build -t cwiegleb/pdc-order-service:no-article .
 docker push cwiegleb/pdc-order-service:no-article
popd 

pushd pdc-csv-upload-service
 docker build -t cwiegleb/pdc-csv-upload-service:no-article .
 docker push cwiegleb/pdc-csv-upload-service:no-article
popd 

pushd pdc-db
 docker build -t cwiegleb/pdc-db:no-article .
 docker push cwiegleb/pdc-db:no-article
popd


