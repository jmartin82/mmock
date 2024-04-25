![Mmock](/docs/logo.png "Mmock logo")
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

```json
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


```yaml
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
docker image pull jordimartin/mmock
docker run -v YOUR_ABS_PATH:/config -p 8082:8082 -p 8083:8083 jordimartin/mmock
```

Or run mmock locally from the command line. (Requires Go 1.8 at least)

```
go get github.com/jmartin82/mmock/...
mmock -h

```

To configure Mmock, use command line flags described in help.


```
    Usage of ./mmock:
	  -config-path string
		Mocks definition folder (default "execution_pathconfig")
	  -console
		Console enabled  (true/false) (default true)
	  -console-ip string
		Console server IP (default "public_ip")
	  -console-port int
		Console server Port (default 8082)
 	  -request-storage-capacity int
		Request storage capacity (0 = infinite) (default 100)
	  -results-per-page uint
		Number of results per page (default 25)
	  -server-ip string
		Mock server IP (default "public_ip")
	  -server-port int
		Mock server Port (default 8083)
	  -server-statistics
		Mock server sends anonymous statistics (default true)
	  -server-tls-port int
		Mock server TLS Port (default 8084)
	  -tls-path string
		TLS config folder (server.crt and server.key should be inside) (default "execution_path/tls")
```

The default logging level is INFO, but you can change it by setting the
environment variable LOG_LEVEL to one of the following:

  * CRITICAL
  * ERROR
  * WARNING
  * NOTICE
  * INFO
  * DEBUG

### Mock

Mock definition:

```json
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
	"callback": {
		"method": "GET|POST|PUT|PATCH|...",
		"url": "http://your-callback/",
		"delay": "string (response delay in s,ms)",
		"headers": {
			"name": ["value"]
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
		"delay": "string (response delay in s,ms)",
		"crazy": "bool (return random 5xx)",
		"priority": "int (matching priority)",
		"webHookURL" : "string (URL endpoint)"
	}
}

```

#### Request

A core feature of Mmock is the ability to return canned HTTP responses for requests matching criteria.

* *host*: Request http host. (without port)
* *method*: Request http method. It allows more than one separated by pipes "|" **Mandatory**
* *path*: Resource identifier. It allows :value matching. **Mandatory**
* *queryStringParameters*: Array of query strings. It allows more than one value for the same key.
* *headers*: Array of headers. It allows more than one value for the same key. **Case sensitive!**
* *cookies*: Array of cookies.
* *body*: Body string. It allows * pattern. It also supports regular expressions for field values within JSON request bodies.

In case of queryStringParameters, headers and cookies, the request can be matched only if all defined keys in mock will be present with the exact or glob value.

Glob matching is available for:
* host
* path
* headers
* cookies
* query strings
* body

Query strings and headers support also global matches (*) in the header/parameter name. For example:
```json
		"headers": {
			"Custom-Header-*": [
				"partial val*"
			]
		}
```

Regexp matching is available for:
- body
- query strings

See https://pkg.go.dev/regexp/syntax for regexp syntax

#### Response (Optional on proxy call)

* *statusCode*: Response status code
* *headers*: Array of headers. It allows more than one value for the same key and vars.
* *cookies*: Array of cookies. It allows vars.
* *body*: Body string. It allows vars.

#### Callback (Optional)

This is used to have mmock make an API request after receiving the mocked request.

* *delay*: Delay before making the callback
* *url*: URL to make a request to
* *method*: Request http method
* *headers*: Array of headers. It allows more than one value for the same key.
* *body*: Body string. It allows vars.
* *timeout*: Duration to allow the callback to respond (default 10s)

#### Control (Optional)

* *scenario*; A scenario is essentially a state machine whose states can be arbitrarily assigned.
* *proxyBaseURL*: If this parameter is present, it sends the request data to the BaseURL and resend the response to de client. Useful if you don't want mock a the whole service. NOTE: It's not necessary fill the response field in this case.
* *delay*: Delay the response. Simulate bad connection or bad server performance.
* *crazy*: Return random server errors (5xx) in some request. Simulate server problems.
* *priority*: Set the priority to avoid match in less restrictive mocks. Higher, more priority.
* *webHookURL*: After any match if this option is defined it will notify the match to the desired endpoint.

### Variable tags

You can use variable data in response. The variables will be defined as tags like this {{nameVar}}

 - URI
 - description

**Request data:** Use them if you want to add request data in your response.

 - request.scheme
 - request.hostname
 - request.port
 - request.path (full path)
 - request.path."*key*"
 - request.query."*key*"
 - request.cookie."*key*"
 - request.fragment
 - request.url (full url with scheme, hostname, port, path and query parameters)
 - request.autority (return scheme, hostname and port (optional))
 - request.body

You can extract information from the request body too, using a dot notation path:
 
 - request.body."*key*" (support for `application/json`, `application/xml` and `application/x-www-form-urlencoded` requests)
 - request.body."*deep*"."*key*" (support for `application/json`, `application/xml` requests)

Quick overview of the path syntax available to extract values form the request: [https://github.com/tidwall/gjson#path-syntax](https://github.com/tidwall/gjson#path-syntax)

You can also use "regex" and "concat" commands to complement GJson query:

- request.body."*deep*"."*key*".regex() (support for `application/json`, `application/xml` requests)
- request.body."*deep*"."*key*".concat() (support for `application/json`, `application/xml` requests)
- request.body."*deep*"."*key*".regex().concat() (support for `application/json`, `application/xml` requests)

**Example request data:**
```json
{
  "email": "hilari@sapo.pt",
  "age": 4,
  "uuid":"^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
  "discarded": "do not return"
}
```
**Example config mock data:**
```json
{
  "email": "{{request.body.email.regex((\@gmail.com))}}",
  "age": {{request.body.age}},
  "uuid": "{{request.body.uuid.regex(\b([0-9a-zA-Z]{4})\b).concat(-878787)}}",
  "discarded": "{{request.body.discarded.concat(, Please!)}}"
}
```
**Example response data:**
```json
{
  "email": "",
  "age": 4,
  "uuid": "2307-878787",
  "discarded": "do not return, Please!"
}
```

**External streams:** Perfect for embedding big payloads or getting data from another service.

 - file.contents(FILE_PATH)
 - http.contents(URL)

**[Fake](https://godoc.org/github.com/icrowley/fake) data:** Useful to provide a more rich and random response.

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
 - fake.Hex(n) - random hexidecimal string n characters in length
 - fake.IPv4
 - fake.Language
 - fake.Model
 - fake.Month
 - fake.MonthShort
 - fake.MonthNum
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
 - fake.IntMinMax(min, max) - random positive number greater or equal to min and less than max
 - fake.Float(n) - random positive floating point number less than n
 - fake.UUID - generates a unique id  

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

## API Specification

MMock uses the Open API Specification (OAI, formerly known as Swagger) to describe its APIs. Our OAI specification schema is hosted at /swagger.json and serves as the canonical definition and comprehensive declaration of all available endpoints.

The OAI specification makes writing client applications easier by: auto-generating boilerplate code (like data object classes) and dealing with authentication and error handling.

You can find a comprehensive set of open tools for the OAI specification at: https://github.com/swagger-api.

#### REST Endpoints

### Verify

The Mmock records the incoming requests in memory (last 100 by default). 
This makes it possible to verify that a request matching a specific pattern was received, and also to fetch the requests details.

**Title** : Get all requests.<br>
**URL** : /api/request/all<br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Get all matched requests with any mock.<br>
**URL** : /api/request/matched<br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Clean all recorded request.<br>
**URL** : /api/request/reset<br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Reset all requests that match with an specific pattern.<br>
**URL** : /api/request/reset_match<br>
**Method** : POST<br>
**Data Params**:  <br>

Like stubbing this call also uses the same DSL to filter and query requests.

```json
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

**Title** : Get all non matched requests.<br>
**URL** : /api/request/unmatched<br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Get all requests that match with an specific pattern.<br>
**URL** : /api/request/verify<br>
**Method** : POST<br>
**Data Params**:  <br>

Like stubbing this call also uses the same DSL to filter and query requests.

```json
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

### Scenario

**Title** : Get all active scenarios.<br>
**URL** : /api/scenarios<br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Clean all scenarios status and pause state.<br>
**URL** : /api/scenarios/reset_all<br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Manually progress a scenario state machine to the given state.<br>
**URL** : /api/scenarios/set/:scenario/:state<br>
**Method** : PUT<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Pause prevents all scenarios state machines from progressing to a new state.<br>
**URL** : /api/scenarios/pause<br>
**Method** : PUT<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Allow scenarios state machines to continue.<br>
**URL** : /api/scenarios/unpause<br>
**Method** : PUT<br>
**Response Codes**: Success (200 OK)<br>

### Mapping

You can manage remotely your stub mappings whenever you need with this simple API:

**Title** : Get all mock definitions.<br>
**URL** : /api/mapping <br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Create new mock definition.<br>
**URL** : /api/mapping/:uri <br>
**Method** : POST<br>
**Response Codes**: Success (201 OK)<br>

**Title** : Get mock definition.<br>
**URL** : /api/mapping/:uri <br>
**Method** : GET<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Update mock definition.<br>
**URL** : /api/mapping/:uri <br>
**Method** : PUT<br>
**Response Codes**: Success (200 OK)<br>

**Title** : Delete mock definition.<br>
**URL** : /api/mapping/:uri <br>
**Method** : DELETE<br>
**Response Codes**: Success (200 OK)<br>)

#### Realtime Endpoints

Also there is a real time endpoint available through WebSockets that broadcast all requests.

**Title** : Get all requests.<br>
**URL** : /echo <br>


#### Output

All endpoints have the same output format:

```json
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

### Statistics

Mmock is collecting anonymous statistics about the usage of the following actions:

Source code: [/statistics/statistics.go](https://github.com/jmartin82/mmock/blob/master/internal/statistics/statistics.go#L30)

- `requests.mock`: Mocks served (number) 
- `requests.console`: Web console usage (number)
- `requests.verify`: Verify requests (number)
- `feature.scenario`: Mocks with scenario feature served (number)
- `feature.proxy`: Mocks with proxy feature served (number)

You can always disable this behavior adding the following flag `-server-statistics=false`


### Contributors
- Amazing request body parsing form [@hmoragrega](https://github.com/hmoragrega)
- Awesome use statistics from [@alfonsfoubert](https://github.com/alfonsfoubert)
- More request variables available thanks to [@Bartek-CMP](https://github.com/Bartek-CMP)
- Scenario pause feature thanks to [@Nekroze](https://github.com/Nekroze)
- Storage reset feature thanks to [@rubencougil](https://github.com/rubencougil)
- Improved docker image thanks to [@daroot](https://github.com/daroot)
- Added the possibility of access to an array index in dynamic responses [@jaimelopez](https://github.com/jaimelopez)
- Create mapping via console thanks to [@inabajunmr](https://github.com/inabajunmr)
- Thanks to [@joel-44](https://github.com/joel-44) for bug fixing 
- Enviroment variables as mock variables thanks to [@marcoreni](https://github.com/marcoreni)
- Support Regular Expressions for Field Values in JSON Request Body thanks to [@rosspatil](https://github.com/rosspatil)
- Improved logging with levels thanks to [@jcdietrich](https://github.com/jcdietrich) [@jdietrich-tc](https://github.com/jdietrich-tc)
- Support for Regular Expressions for QueryStringParameters [@jcdietrich](https://github.com/jcdietrich) [@jdietrich-tc](https://github.com/jdietrich-tc)
- Support for URI and Description tags [@jcdietrich](https://github.com/jcdietrich) [@jdietrich-tc](https://github.com/jdietrich-tc)

### Contributing

As of version 3.0.0, Mmock is available as a Go module. Therefore a Go version capable of understanding /vN suffixed imports is required:

1.9.7+
1.10.3+
1.11+

If you make any changes, please run ```go fmt ./...``` before submitting a pull request.

### Licence

Copyright © 2016 - 2020, Jordi Martín (http://jordi.io)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
