# Get Retail API Predictions

[![Workflows](https://github.com/wintermi/get-retail-predictions/workflows/Go/badge.svg)](https://github.com/wintermi/get-retail-predictions/actions)
[![Go Report](https://goreportcard.com/badge/github.com/wintermi/get-retail-predictions)](https://goreportcard.com/report/github.com/wintermi/get-retail-predictions)
[![License](https://img.shields.io/github/license/wintermi/get-retail-predictions.svg)](https://github.com/wintermi/get-retail-predictions/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/wintermi/get-retail-predictions?include_prereleases)](https://github.com/wintermi/get-retail-predictions/releases)


## Description

A command line application designed to provide a simple method of requesting predictions from a Google Cloud Retail API model for all sets of parameters contained within the input file.

```
USAGE:
    get-retail-predictions -p PROJECT_NUMBER -s SERVING_CONFIG -i INPUT_FILE

ARGS:
  -b string
    	Branch (default "0")
  -c string
    	Catalog (default "default_catalog")
  -f string
    	Filter String
  -i string
    	Parameter Input File  (Required)
  -l string
    	Location (default "global")
  -n int
    	Number of Results, 1 to 100 (default 5)
  -p string
    	Project Number  (Required)
  -s string
    	Serving Config  (Required)
  -v	Output Verbose Detail
```

## Example Parameter Input File

```
[
  {
    "event_type": "detail-page-view",
    "visitor_id": "1",
    "product_details": [{ "product": { "id": "100" } }]
  },
  {
    "event_type": "detail-page-view",
    "visitor_id": "1",
    "product_details": [{ "product": { "id": "200" } }]
  }
]
```


## License

**get-retail-predictions** is released under the [Apache License 2.0](https://github.com/wintermi/get-retail-predictions/blob/main/LICENSE) unless explicitly mentioned in the file header.
