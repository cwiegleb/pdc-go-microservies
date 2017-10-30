# PDC 
## How to Run 
  - `./build_standalone_apps.sh`
  - `./dockerize_apps.sh`
  - `docker-compose -f docker-stack.yml up`
  - `docker run --rm --link pdcservices_db_1:postgresContainer --network pdcservices_default cwiegleb/pdc-db:latest`

## DB Model
### Status Article
 - 1 (true) Available 
 - 0 (false) Not Available

### Status Order
 - 0 Initial
 - 1 Closed

## Article Service 
### POST request article service
```{"Text":"Test Article 1","Size":"M","DealerID":1,"Available":true,"Costs":12,"Currency":"EUR"}```
### PUT request article service
```{"ID":1,"Text":"Test Article 1","Size":"M","DealerID":1,"Available":true,"Costs":12,"Currency":"EUR"}```
## Dealer Service 
### POST request dealer service (without articles)
``` {"Text":"Test 2","Articles":null}```
### PUT request dealer service (without articles)
``` {"ID":1,"Text":"Test 2","Articles":null}```
## Cashbox Service 
### POST request cashbox service (without orders)
```{"Name":"Cashbox 1","ValidFromDate":"2017-10-01T20:00:00+02:00","ValidToDate":"2017-10-10T20:00:00+02:00","Orders":null}```
### PUT request cashbox service (without orders)
```{"ID": 1, Name":"Cashbox 1","ValidFromDate":"2017-10-01T20:00:00+02:00","ValidToDate":"2017-10-10T20:00:00+02:00","Orders":null}```
## Order Service 
### POST request order service (with orderline)
```{"CashboxID":1,"OrderStatus":"2","CreationDate":"2017-10-01T20:00:00+02:00","OrderLines":[{"ArticleID":1,"Price":12,"Currency":"EUR"}]}```
### PUT request order service (with orderline)
```{"ID": 1, CashboxID":1,"OrderStatus":"1","CreationDate":"2017-10-01T20:00:00+02:00","OrderLines":[{"ArticleID":1,"Price":12,"Currency":"EUR"}]}```