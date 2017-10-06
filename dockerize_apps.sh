#!/bin/bash

pushd pdc-article-service
 docker build -t cwiegleb/pdc-article-service:latest .
popd

pushd pdc-cashbox-service
 docker build -t cwiegleb/pdc-cashbox-service:latest .
popd

pushd pdc-dealer-service
 docker build -t cwiegleb/pdc-dealer-service:latest .
popd

pushd pdc-order-service
 docker build -t cwiegleb/pdc-order-service:latest .
popd 

pushd pdc-db
 docker build -t cwiegleb/pdc-db:latest .
popd


