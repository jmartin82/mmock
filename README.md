Mmock
=========
[![Build Status](https://travis-ci.org/jmartin82/mmock.svg?branch=master)](https://travis-ci.org/jmartin82/mmock)

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
* Variables in response (fake or request data)
* Route patterns may include named parameters (/hello/:name)
* Glob matching ( * hello * )
* Match request by method, URL params, query string, headers, cookies and bodies.
* Mock definitions hot replace (edit your mocks without restart)
* Web interface to view requests data (method,path,headers,cookies,body,etc..)
* Stateful behaviour with scenarios
* Verifying
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

You can see more complex examples [here](/config).

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
	"description": "Some text that describes the intended usage of the current configuration",
	"request": {
		"host": "example.com",
		"method": "GET|POST|PUT|PATCH|...",
		"path": "/your/path/:variable",
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
		"body": "Response body"
	},
	"control": {
		"scenario": {
			"name": "string (scenario name)",
            "requiredState": [
				"not_started (default state)",
                "another_state_name"
            ],
            "newState": "new_stat_neme"
        },
		"proxyBaseURL": "string (original URL endpoint)",
		"delay": "int (response delay in seconds)",
		"crazy": "bool (return random 5xx)",
		"priority": "int (matching priority)"
	}
}

```

#### Request

A core feature of Mmock is the ability to return canned HTTP responses for requests matching criteria.

* *host*: Request http host. (without port)
* *method*: Request http method. It allows more than one separated by pipes "|" **Mandatory**
* *path*: Resource identifier. It allows :value matching. **Mandatory**
* *queryStringParameters*: Array of query strings. It allows more than one value for the same key.
* *headers*: Array of headers. It allows more than one value for the same key.
* *cookies*: Array of cookies.
* *body*: Body string. It allows * pattern.

In case of queryStringParameters, headers and cookies, the request can be matched only if all defined keys in mock will be present with the exact value.

#### Response (Optional on proxy call)

* *statusCode*: Request http method.
* *headers*: Array of headers. It allows more than one value for the same key and vars.
* *cookies*: Array of cookies. It allows vars.
* *body*: Body string. It allows vars.

#### Control (Optional)

* *scenario*; A scenario is essentially a state machine whose states can be arbitrarily assigned.
* *proxyBaseURL*: If this parameter is present, it sends the request data to the BaseURL and resend the response to de client. Useful if you don't want mock a the whole service. NOTE: It's not necessary fill the response field in this case.
* *delay*: Delay the response in seconds. Simulate bad connection or bad server performance.
* *crazy*: Return random server errors (5xx) in some request. Simulate server problems.
* *priority*: Set the priority to avoid match in less restrictive mocks.

### Scenarios

With the scenarios you can simulate a stateful service. It's useful to create test doubles.

A scenario is a state machine and you can assign an arbitrarily state.

When mmock recieve a new request and there is an scenario defined in the matching mock, mmock checks if the mock is valid for the current state. Also a new scenario state can be set after the mock match.

By default all scenarios has the state "not_started" until some mock triggers a new one.

Example of REST services using scenarios:

```
+---------------------------------------------------------------------------------------------+
|                                                                                             |
|   GET /user                     POST /user                     GET /user                    |
|   StatusCode: 404               StatusCode: 201                StatusCode: 200              |
|                                                                                             |
|  +-------------------------+   +---------------------------+   +-------------------------+  |
|  |                         |   |                           |   |                         |  |
|  | requiredState:created   +-> | requiredState:not_started +-> |  requiredState: created |  |
|  |                         |   | newState: created         |   |                         |  |
|  |                         |   |                           |   |                         |  |
|  +-------------------------+   +---------------------------+   +-------------------------+  |
|                                                                                             |
+---------------------------------------------------------------------------------------------+
```

Working examples [here](/config/crud)

### Verify

The Mmock rrecords all requests it receives in memory (at least until it is reset). 
This makes it possible to verify that a request matching a specific pattern was received, and also to fetch the requests details.

You can get that data throught the web console server (by default it has the same ip of mock sever and port 8082).

#### REST Endpoints

**Title** : Get all requests.<br>
**URL** : /request/all<br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>


**Title** : Get all matched requests with any mock.<br>
**URL** : /request/matched<br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Get all non matched requests.<br>
**URL** : /request/unmatched<br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Get all requests that match with an specific pattern.<br>
**URL** : /request/verify<br>
**Method** : POST<br>
**Data Params**:  <br>
Like stubbing this call also uses the same DSL to filter and query requests.

```
{
	"host": "example.com",
	"method": "GET|POST|PUT|PATCH|... (Mandatory)", 
	"path": "/your/path/:variable (Mandatory)",
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
}
```
**Response Codes**: Success (200 OK)<br>

#### Realtime Endpoints

Also there is a real time endpoint available through WebSockets that broadcast all requests.

**Title** : Get all requests.<br>
**URL** : /echo <br>


#### Output

All enpoints have the same output format:

```
[
  {
    "time": 1486563983,
    "request": {
      "host": "192.168.20.209",
      "method": "GET",
      "path": "/hello",
      "queryStringParameters": {},
      "headers": {
        "Accept": [
          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
        ]
      },
      "cookies": {},
      "body": ""
    },
    "response": {
      "statusCode": 200,
      "headers": null,
      "cookies": {
        "visit": "true"
      },
      "body": "Hello world!"
    },
    "result": {
      "match": true,
      "errors": null
    }
  }
]
```


### Variable tags

You can use variable data (random data or request data) in response. The variables will be defined as tags like this {{nameVar}}

Request data:

 - request.query."*key*"
 - request.path."*key*"
 - request.cookie."*key*"
 - request.url
 - request.body


[Fake](https://godoc.org/github.com/icrowley/fake) data:

 - fake.Brand
 - fake.Character
 - fake.Characters
 - fake.CharactersN(n)
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
 - fake.DigitsN(n)
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
 - fake.ParagraphsN(n)
 - fake.Phone
 - fake.Product
 - fake.Sentence
 - fake.Sentences
 - fake.SentencesN(n)
 - fake.SimplePassword
 - fake.State
 - fake.StateAbbrev
 - fake.Street
 - fake.StreetAddress
 - fake.UserName
 - fake.WeekDay
 - fake.Word
 - fake.Words
 - fake.WordsN(n)
 - fake.Zip
 - fake.Int(n) - random positive integer less than or equal to n
 - fake.Float(n) - random positive floating point number less than n
 - fake.UUID - generates a unique id  



### Contributors
- Awesome use statistics from [@alfonsfoubert](https://github.com/alfonsfoubert)
- Is this not enough? [@vtrifonov](https://github.com/vtrifonov) is working in a fork with a really advanced features. [HTTP API Mock](https://github.com/vtrifonov/http-api-mock)

### Contributing

Clone this repository to ```$GOPATH/src/github.com/jmartin82/mmock``` and type ```go get .```.

Requires Go 1.6+ to build.

If you make any changes, run ```go fmt ./...``` before submitting a pull request.

### Licence

Copyright ©‎ 2016 - 2018, Jordi Martín (http://jordi.io)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
