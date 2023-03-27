
# getircase-study

# Run  The Application
-  Port and dsn are hardcoded into server 
- Simply run `go run ./cmd` command
## Deployment
App running on a AWS EC2 instance( ipaddr: http://18.185.42.203)


## Custom Responses
| code | description                  |
|------------|------------------------|
| 0          | Success                |
| 1          | BadRequest             |
| 2          | UnsoppertedMethod      |
| 3          | InternalServerError    |

# Endpoints
- ## Get Entry From In-memory Database
```bash
curl --location 'http://18.185.42.203/in-memory?key=test_key'
```
- ## Set Entry To In-memory Database
```bash
curl --location 'http://18.185.42.203/in-memory' \
--header 'Content-Type: application/json' \
--data '{
"key": "test_key",
"value": "test_value"
}'
 ```
- ## Get Data From Mongo DB
```bash
curl --location 'http://18.185.42.203/records' \
--header 'Content-Type: application/json' \
--data '{
"startDate": "2011-01-28",
"endDate": "2018-02-02",
"minCount": 2700,
"maxCount": 3000
}'
```