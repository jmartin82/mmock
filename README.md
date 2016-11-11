Mmock 
=========

Mmock is a testing and fast prototyping tool for developers:

Easy and fast HTTP mock server.

* Download Mmock
* Create a mock definition.
* Configure your application endpoints to use Mmock
* Receive the expected responses
* Inspect the requests in the web UI
* Release it to a real server

Built with Go - Mmock runs without installation on multiple platforms.

### Features

* Easy mock definition via JSON or YAML
* Variables in response (fake or request data, including regex support)
* Persist request body to file and load request from file
* Glob matching ( /a/b/* )
* Match request by method, URL params, headers, cookies and bodies.
* Mock definitions hot replace (edit your mocks without restart)
* Web interface to view requests data (method,path,headers,cookies,body,etc..)
* Proxy mode
* Fine grain log info in web interface
* Real-time updates using WebSockets
* Priority matching
* Crazy mode for failure testing
* Public interface auto discover
* Lightweight and portable
* No installation required

### Example

![Video of Mmock](/docs/example.gif "Mmock example")

Mock definition file example:

```
{
	"request": {
		"method": "GET",
		"path": "/hello/*"
	},
	"response": {
		"statusCode": 200,
		"headers": {
			"Content-Type":["application/json"]
		},
		"body": "{\"hello\": \"{{request.query.name}}, my name is {{fake.FirstName}}\"}"
	}
}

```
Or


```
---
request:
  method: GET
  path: "/hello/*"
response:
  statusCode: 200
  headers:
    Content-Type:
    - application/json
  body: '{"hello": "{{request.query.name}}, my name is {{fake.FirstName}}"}'

```



### Getting started

Either:

Run it from Docker using the provided ```Dockerfile``` or [from Docker Hub](https://hub.docker.com/r/jordimartin/mmock/)

```
go get github.com/jmartin82/mmock
docker build -t mmock/mmock .
docker run -v YOUR_ABS_PATH:/config -p 8082:8082 -p 8083:8083  mmock/mmock
```


Or run mmock locally from the command line. 

```
go get github.com/jmartin82/mmock
mmock -h

```

To configure Mmock, use command line flags described in help.


```
    Usage of ./mmock:
      -cconsole-port int
          Console server Port (default 8082)
      -config-path string
          Mocks definition folder (default "execution_path/config")
      -config-persist-path
          Path to the folder where requests can be persisted (default "execution_path/data") 
      -console
          Console enabled  (true/false) (default true)
      -console-ip string
          Console Server IP (default "public_ip")
      -server-ip string
          Mock server IP (default "public_ip")
      -server-port int
          Mock Server Port (default 8083)
```

### Mock

Mock definition:

```
{
	"request": {
		"method": "GET|POST|PUT|PATCH|...",
		"path": "/your/path/*",
		"queryStringParameters": {
			"name": ["value"],
			"name": ["value", "value"]
		},
		"headers": {
			"name": ["value"]
		},
		"cookies": {
			"name": "value"
		},
		"body": "Expected Body"
	},
	"response": {
		"statusCode": "int (2xx,4xx,5xx,xxx)",
		"headers": {
			"name": ["value"]
		},
		"cookies": {
			"name": "value"
		},
		"body": "Response body",
		"persisted": {
            "name" : "/users/user-{{request.url./your/path/(?P`<value>`\\d+)}}.json",
            "notFound": {
                "statusCode": 404,
                "body": "404: Not Found"
				"bodyAppend": " File for Id : {{request.url./your/path/(?P`<value>`\\d+)}}"
            },
            "bodyAppend": "{ \"id\": {{request.url./your/path/(?P`<value>`\\d+)}} }"
        }
	},
	"persist" : {
		"name" : "/users/user-{{request.url./your/path/(?P`<value>`\\d+)}}.json",
        "delete": false
	},
	"control": {
		"proxyBaseURL": "string (original URL endpoint)
		"delay": "int (response delay in seconds)",
		"crazy": "bool (return random 5xx)",
		"priority": "int (matching priority)"
	}
}

```

#### Request

This mock definition section represents the expected input data. I the request data match with mock request section, the server will response the mock response data.  

* *method*: Request http method. **Mandatory**
* *path*: Resource identifier. It allows * pattern. **Mandatory**
* *queryStringParameters*: Array of query strings. It allows more than one value for the same key.
* *headers*: Array of headers. It allows more than one value for the same key.
* *cookies*: Array of cookies.
* *body*: Body string. It allows * pattern.

To do a match with queryStringParameters, headers, cookies. All defined keys in mock will be present with the exact value.

#### Response

* *statusCode*: Request http method.
* *headers*: Array of headers. It allows more than one value for the same key and vars.
* *cookies*: Array of cookies. It allows vars.
* *body*: Body string. It allows vars.
* *bodyAppend*: Additional text or json object to be appended to the body. It allows vars. If the body is in JSON format and the bodyAppend is JSON the two JSONs will be merged and bodyAppend fields will replace any field from both JSONS. If the body type is not JSON the strings will be concatenated.  
* *persisted*: Configuration for reading the body content from file.

##### Persisted

* *name*: The relative path from config-persist-path to the file where the response body to be loaded from. It allows vars.
* *notFound*: The status code and body which will be returned if the file does not exist. The default values are statusCode: **404** and body: **Not Found**
* *notFound.statusCode*: The status code to be returned if the file is not found. The default value is **404**
* *notFound.body*: The body to be returned if the file is not found. It allows vars. The default value is **Not Found**
* *notFound.bodyAppend*: Additional text or json object to be appended to the body if the file is not found. It allows vars.
* *bodyAppend*: Additional text or json object to be appended to the body loaded from the file. It allows vars.
* *persisted*: Configuration for reading the body content from file.

#### Persist

* *name*: The relative path from config-persist-path to the file where the ressponse body will be persisted. It allows vars.
* *delete*: True or false. This is useful for making **DELETE** verb to delete the file.

#### Control

* *proxyBaseURL*: If this parameter is present, it sends the request data to the BaseURL and resend the response to de client. Useful if you don't want mock a the whole service. NOTE: It's not necessary fill the response field in this case.
* *delay*: Delay the response in seconds. Simulate bad connection or bad server performance.
* *crazy*: Return random server errors (5xx) in some request. Simulate server problems.
* *priority*: Set the priority to avoid match in less restrictive mocks.

### Variable tags

You can use variable data (random data or request data) in response. The variables will be defined as tags like this {{nameVar}} 

Request data:

 - request.query."*key*"
 - request.cookie."*key*"
 - request.url
 - request.body
 - request.url."regex to match value"
 - request.body."regex to match value"

> Regex: The regex should contain a group named **value** which will be matched and its value will be returned. E.g. if we want to match the id from this url **/your/path/4** the regex should look like **/your/path/(?P`<value>`\\d+)**. Note that in *golang* the named regex group match need to contain a **P** symbol after the question mark. The regex should be prefixed either with **request.url.** or **request.body.** considering your input.


Fake data:

 - fake.Brand
 - fake.Character
 - fake.Characters
 - fake.City
 - fake.Color
 - fake.Company
 - fake.Continent
 - fake.Country
 - fake.CreditCardVisa
 - fake.CreditCardMasterCard
 - fake.CreditCardAmericanExpress
 - fake.Currency
 - fake.CurrencyCode
 - fake.Day
 - fake.Digits
 - fake.EmailAddress
 - fake.FirstName
 - fake.FullName
 - fake.LastName
 - fake.Gender
 - fake.IPv4
 - fake.Language
 - fake.Model
 - fake.Month
 - fake.Year
 - fake.Paragraph
 - fake.Paragraphs
 - fake.Phone
 - fake.Product
 - fake.Sentence
 - fake.Sentences
 - fake.SimplePassword
 - fake.State
 - fake.StateAbbrev
 - fake.Street
 - fake.StreetAddress
 - fake.UserName
 - fake.WeekDay
 - fake.Word
 - fake.Words
 - fake.Zip


### Benchmark

Basic benchmark with ab (Apache HTTP server benchmarking tool)

Intel(R) Pentium(R) CPU 2127U @ 1.90GHz
350 Concurrent
20000 Request

```
ab -k -c 350 -n 20000 http://YOURIP:8083/

Server Software:        
Server Hostname:        172.17.0.2
Server Port:            8083

Document Path:          /
Document Length:        0 bytes

Concurrency Level:      350
Time taken for tests:   11.348 seconds
Complete requests:      20000
Failed requests:        0
Non-2xx responses:      20000
Keep-Alive requests:    20000
Total transferred:      2940000 bytes
HTML transferred:       0 bytes
Requests per second:    1762.49 [#/sec] (mean)
Time per request:       198.583 [ms] (mean)
Time per request:       0.567 [ms] (mean, across all concurrent requests)
Transfer rate:          253.01 [Kbytes/sec] received

```

### Contributing

Clone this repository to ```$GOPATH/src/github.com/jmartin82/mmock``` and type ```go get .```.

Requires Go 1.4+ to build.

If you make any changes, run ```go fmt ./...``` before submitting a pull request.

### Licence

Copyright ©‎ 2016 - 2017, Jordi Martín (http://jordi.io)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
