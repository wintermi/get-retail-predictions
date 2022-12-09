# Get Retail API Predictions
[![Go Workflow Status](https://github.com/winterlabs-dev/get-retail-predictions/workflows/Go/badge.svg)](https://github.com/winterlabs-dev/get-retail-predictions/actions/workflows/go.yml)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/winterlabs-dev/get-retail-predictions)](https://goreportcard.com/report/github.com/winterlabs-dev/get-retail-predictions)&nbsp;[![license](https://img.shields.io/github/license/winterlabs-dev/get-retail-predictions.svg)](https://github.com/winterlabs-dev/get-retail-predictions/blob/main/LICENSE)&nbsp;[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/winterlabs-dev/get-retail-predictions?include_prereleases)](https://github.com/winterlabs-dev/get-retail-predictions/releases)


## Description
A command line application designed to provide a simple method to request predictions from a given Retail API model.

```
USAGE:
    get-retail-predictions -p PROJECT_NUMBER -s SERVING_CONFIG

ARGS:
  -c string
    	Catalog  (Required) (default "default_catalog")
  -experiment string
    	Experiment Group
  -filter string
    	Filter String
  -l string
    	Location  (Required) (default "global")
  -n int
    	Number of Predictions  (Required) (default 10)
  -p string
    	Google Cloud Project Number  (Required)
  -product string
    	Product ID  (Required)
  -s string
    	Serving Config  (Required)
  -type string
    	Event Type  (Required)
  -v	Output Verbose Detail
  -visitor string
    	Visitor ID  (Required)
```
