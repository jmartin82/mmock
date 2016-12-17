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
* Glob matching ( /a/b/* )
* Match request by method, URL params, query string, headers, cookies and bodies.
* Mock definitions hot replace (edit your mocks without restart)
* Web interface to view requests data (method,path,headers,cookies,body,etc..)
* Stateful behaviour with scenarios
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
		"proxyBaseURL": "string (original URL endpoint)
		"delay": "int (response delay in seconds)",
		"crazy": "bool (return random 5xx)",
		"priority": "int (matching priority)"
	}
}

```

#### Request

A core feature of Mmock is the ability to return canned HTTP responses for requests matching criteria. 

* *method*: Request http method. **Mandatory**
* *path*: Resource identifier. It allows * pattern. **Mandatory**
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

When mmock recieve a new request and one scenario is defined in the matching mock it checks if the mock is valid for the current state.

When mmock return a new mock we can set the state value for the scenario.

By default all scenarios has the state "not_started" until some mock triggers a new one.

Example of REST services using scenarios:

```
+----------------------------------------------------------------------------------------+
|                                                                                        |
|   GET /user                  POST /user                     GET /user                  |
|   StatusCode: 404            StatusCode: 201                StatusCode: 200            |
|                                                                                        |
|                                                                                        |
|  +----------------------+   +--------------------------+   +------------------------+  |
|  |                      |   |                          |   |                        |  |
|  | matchStatus: created +-> | matchStatus: not_started +-> |  matchStatus: created  |  |
|  |                      |   | nextStatus: created      |   |                        |  |
|  +----------------------+   +--------------------------+   +------------------------+  |
|                                                                                        |
+----------------------------------------------------------------------------------------+
```

Working examples [here](/config/crud) 

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
- [@vtrifonov](https://github.com/vtrifonov)

### Contributing

Clone this repository to ```$GOPATH/src/github.com/jmartin82/mmock``` and type ```go get .```.

Requires Go 1.6+ to build.

If you make any changes, run ```go fmt ./...``` before submitting a pull request.

### Licence

Copyright ©‎ 2016 - 2018, Jordi Martín (http://jordi.io)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
