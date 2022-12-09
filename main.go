// Copyright 2022, Matthew Winter
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger
var applicationText = "%s 0.1.0%s"
var copyrightText = "Copyright 2022, Matthew Winter\n"
var indent = "..."

var helpText = `
A command line application designed to provide a simple method to request
predictions from a given Retail API model.

Use --help for more details.


USAGE:
    get-retail-predictions -p PROJECT_NUMBER -s SERVING_CONFIG

ARGS:
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, applicationText, filepath.Base(os.Args[0]), "\n")
		fmt.Fprint(os.Stderr, copyrightText)
		fmt.Fprint(os.Stderr, helpText)
		flag.PrintDefaults()
	}

	// Define the Long CLI flag names
	var targetProject = flag.String("p", "", "Google Cloud Project Number  (Required)")
	var targetLocation = flag.String("l", "global", "Location  (Required)")
	var targetCatalog = flag.String("c", "default_catalog", "Catalog  (Required)")
	var targetServingConfig = flag.String("s", "", "Serving Config  (Required)")
	var requestNumberPredictions = flag.Int("n", 10, "Number of Predictions  (Required)")
	var requestEventType = flag.String("type", "", "Event Type  (Required)")
	var requestVisitorID = flag.String("visitor", "", "Visitor ID  (Required)")
	var requestProductID = flag.String("product", "", "Product ID  (Required)")
	var requestFilter = flag.String("filter", "", "Filter String")
	var requestExperiementID = flag.String("experiment", "", "Experiment Group")
	var verbose = flag.Bool("v", false, "Output Verbose Detail")

	// Parse the flags
	flag.Parse()

	// Validate the Required Flags
	if *targetProject == "" || *targetLocation == "" || *targetCatalog == "" || *targetServingConfig == "" || *requestEventType == "" || *requestVisitorID == "" || *requestProductID == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Setup Zero Log for Consolo Output
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logger = zerolog.New(output).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.DurationFieldInteger = true
	if *verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Output Header
	logger.Info().Msgf(applicationText, filepath.Base(os.Args[0]), "")
	logger.Info().Msg("Arguments")
	logger.Info().Str("Project Number", *targetProject).Msg(indent)
	logger.Info().Str("Location", *targetLocation).Msg(indent)
	logger.Info().Str("Catalog", *targetCatalog).Msg(indent)
	logger.Info().Str("Serving Config", *targetServingConfig).Msg(indent)
	logger.Info().Int("Number of Predictions", *requestNumberPredictions).Msg(indent)
	logger.Info().Str("Event Type", *requestEventType).Msg(indent)
	logger.Info().Str("Visitor ID", *requestVisitorID).Msg(indent)
	logger.Info().Str("Product ID", *requestProductID).Msg(indent)
	logger.Info().Str("Filter String", *requestFilter).Msg(indent)
	logger.Info().Str("Experiment Group", *requestExperiementID).Msg(indent)
	logger.Info().Msg("Begin")

	// Setup the Prediction Request
	var prediction = Prediction{
		project:       *targetProject,
		location:      *targetLocation,
		catalog:       *targetCatalog,
		servingConfig: *targetServingConfig,
		numberResults: *requestNumberPredictions,
		eventType:     *requestEventType,
		visitor:       *requestVisitorID,
		product:       *requestProductID,
		filter:        *requestFilter,
		experiment:    *requestExperiementID,
	}

	// Request to get a Retail API Prediction
	err := prediction.ExecuteRequest()
	if err != nil {
		logger.Error().Err(err).Msg("Prediction Request Failed")
		os.Exit(1)
	}
	logger.Info().Msg("End")

}
