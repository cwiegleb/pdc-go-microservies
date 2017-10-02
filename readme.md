# PDC 

## Stack
### DB Stack
    image: postgres
    environment:
        POSTGRES_PASSWORD: xxxxxxxx
        POSTGRES_USER: pdc
        POSTGRES_DB: pdcDB
    ports:
      - 5432:5432

### Service Stack
tbd


## DB Model
### Status Article
 - 1 (true) Available 
 - 0 (false) Not Available

### Status Order
 - 1 Open
 - 2 Closed

## Examples
### Request dealer service (without articles)
    {"ID":1,"CreatedAt":"2017-10-01T20:00:00+02:00","UpdatedAt":"2017-10-01T20:00:00+02:00","DeletedAt":null,"Text":"Test 2","Articles":null}

### Request cashbox service (without orders)
    {"ID":1,"CreatedAt":"2017-10-01T20:00:00+02:00","UpdatedAt":"2017-10-01T20:00:00+02:00","DeletedAt":null,"Name":"Cashbox 1","ValidFromDate":"2017-10-01T20:00:00+02:00","ValidToDate":"2017-10-10T20:00:00+02:00","Orders":null}

### Request article service
    {"ID":1,"CreatedAt":"2017-10-01T20:00:00+02:00","UpdatedAt":"2017-10-01T20:00:00+02:00","DeletedAt":null,"Text":"Test Article 1","Size":"M","DealerID":1,"Available":true,"Costs":12,"Currency":"EUR"}

### Request order service (with orderline)
    {"ID":1,"CreatedAt":"2017-10-01T20:00:00+02:00","UpdatedAt":"2017-10-01T20:00:00+02:00","DeletedAt":null,"CashboxID":1,"OrderStatus":"2","CreationDate":"2017-10-01T20:00:00+02:00","OrderLines":[{"ID":1,"CreatedAt":"2017-10-01T20:00:00+02:00","UpdatedAt":"2017-10-01T20:00:00+02:00","DeletedAt":null,"OrderID":1,"ArticleID":1,"Price":12,"Currency":"EUR"}]}