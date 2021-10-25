# Snapper

![snapper](snapper.png)

### A Web microservice for capturing a website's OpenGraph data built in Golang

## Building Snapper

```shell
git clone https://github.com/owl-93/snapper
cd snapper
go build .
```

Optionally give the executable a name
```shell
go build -o <executable name> .
```

## Running Snapper (Default Options)
By default, snapper will run on port 8888 and will try to set up caching with redis running on the localhost at the default redis port of 6379,
and a cache TTL of 24 hours. See the below section on configuring snapper for more options to fit your use case.
```shell
./snapper
```

## Running Snapper with Arguments & Options

### Specifying a different port
By default, snapper runs on port `8888`. You can tell snapper to use a different port by passing the port as a command line argument
```shell
./snapper --port 8081
```

### Specifying a Redis instance to back caching
You can pass the `--cache` flag followed by a redis connection address to point to a redis instance to use. Note that an invalid 
connection URI does not exit the application it will run as if started with the `--no-cache` flag. 

```shell
    ./snapper --cache "some-redis-instance:6379"
```

### Setting the cache TTL
You can set the cache TTL for caching page metadata. The default cache TTL is 24 hours, and the cache TTL is specified in number of hours
using the `--cache-ttl` option. The TTL applies to each fetched page, not the entire cache

running snapper with a cache life of 12 hours
```shell
    /.snapper --cache-ttl 12
```


### Globally disabling the cache
You can disable caching entirely for the application by passing the --no-cache flag to snapper. Note that this is equivalent
to passing the `forceRefresh` option in every request to snapper. However, this option also prevents snapper from storing data in
the cache as well. With the `forceRefresh` option, the fetched data is simply not read from the cache, but it is still stored in the cache.
This means that even after a request that specifies the `forceRefresh` option, subsequent requests to snapper for that page that don't specify
the `forceRefresh` option will be read from the cache if there is a cache hit (the cache entry hasn't expired)
**note that because `--no-cache` is an option and not an argument, it must come after any named arguments you specify
```shell
    ./snapper --no-cache
```



## Using Snapper
To snap a webpage's Opengraph metadata, just make a http POST request to `/` with
the target website specified in the request body using the key `page`. You can optionally
pass the `forceRefresh` option in the request body to force snapper to fetch the latest metadata
and not use any cached values if present, and the optional `raw` key to specify your desired response format.

### Request Body Format
```typescript
{
  page: string // the url of the page you wish to fetch metadata for,
  forceRefresh: boolean //(optional) - optionally tell snapper to ignore any cached data and fetch the latest page data (cache will be updated),
  raw: boolean //(optional) - optionally tell snapper that you want a response type with array of MetaTag objects containing the property names and content values
}
```


### Response Body Formats

#### Default Format
The default format contains the 6 main Opengraph property types
```typescript
{
    url: string // og:url
    title: string // og:title
    description: string // og:description
    image: string // og:image
    type: string //og:type
    locale: string //og:locale
}
```

#### Raw Format
The raw format contains the full array of Opengraph property tag names & values
```typescript
[
    {
        name: string, // opengraph key
        value: string // value for that key
    },
    {
        name: string,
        value: string
    },
    //...
]
```



## Examples

### Default format request & response
```shell
curl --location --request POST 'http://localhost:8888/' 
--header 'Content-Type: application/json' 
--data-raw '{
    "page": "https://github.com/owl-93/snapper"
}'
```

Response Code 200\
Response Body:

```json
{
  "url": "https://github.com/owl-93/snapper",
  "title": "GitHub - owl-93/snapper: Golang based web site opengraph data scraper with caching",
  "description": "Golang based web site opengraph data scraper with caching - GitHub - owl-93/snapper: Golang based web site opengraph data scraper with caching",
  "image": "https://opengraph.githubassets.com/b63c65ebc5492a24715bae27d7efa53e333686a06cce9ab11ecc0c9ec64615ab/owl-93/snapper",
  "type": "object",
  "locale": ""
}
```

### Raw format request & response
```shell
curl --location --request POST 'http://localhost:8888/' 
--header 'Content-Type: application/json' 
--data-raw '{
    "page": "https://github.com/owl-93/snapper",
    "raw" : true
}'
```

Response Code 200\
Response Body:

**Note that the Raw response type contains more tags and data than the default response type.**
```json
[
    {
        "name": "fb:app_id",
        "value": "1401488693436528"
    },
    {
        "name": "og:image",
        "value": "https://opengraph.githubassets.com/b63c65ebc5492a24715bae27d7efa53e333686a06cce9ab11ecc0c9ec64615ab/owl-93/snapper"
    },
    {
        "name": "og:image:alt",
        "value": "Golang based web site opengraph data scraper with caching - GitHub - owl-93/snapper: Golang based web site opengraph data scraper with caching"
    },
    {
        "name": "og:image:width",
        "value": "1200"
    },
    {
        "name": "og:image:height",
        "value": "600"
    },
    {
        "name": "og:site_name",
        "value": "GitHub"
    },
    {
        "name": "og:type",
        "value": "object"
    },
    {
        "name": "og:title",
        "value": "GitHub - owl-93/snapper: Golang based web site opengraph data scraper with caching"
    },
    {
        "name": "og:url",
        "value": "https://github.com/owl-93/snapper"
    },
    {
        "name": "og:description",
        "value": "Golang based web site opengraph data scraper with caching - GitHub - owl-93/snapper: Golang based web site opengraph data scraper with caching"
    }
]
```
