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
var applicationText = "%s 0.2.0%s"
var copyrightText = "Copyright 2022, Matthew Winter\n"
var indent = "..."

var helpText = `
A command line application designed to provide a simple method of requesting
predictions from a Google Cloud Retail API model for all sets of parameters
contained within an input file.

Use --help for more details.


USAGE:
    get-retail-predictions -p PROJECT_NUMBER -s SERVING_CONFIG -i INPUT_FILE

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
	var projectNumber = flag.String("p", "", "Project Number  (Required)")
	var location = flag.String("l", "global", "Location")
	var catalog = flag.String("c", "default_catalog", "Catalog")
	var branch = flag.String("b", "0", "Branch")
	var servingConfig = flag.String("s", "", "Serving Config  (Required)")
	var parameterInputFile = flag.String("i", "", "Parameter Input File  (Required)")
	var numberResults = flag.Int("n", 5, "Number of Results, 1 to 100")
	var filterString = flag.String("f", "", "Filter String")
	var verbose = flag.Bool("v", false, "Output Verbose Detail")

	// Parse the flags
	flag.Parse()

	// Validate the Required Flags
	if *projectNumber == "" || *location == "" || *catalog == "" || *branch == "" || *servingConfig == "" || *parameterInputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Verify Number of Results is between 1 and 100
	if *numberResults < 1 || *numberResults > 100 {
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
	logger.Info().Str("Project Number", *projectNumber).Msg(indent)
	logger.Info().Str("Location", *location).Msg(indent)
	logger.Info().Str("Catalog", *catalog).Msg(indent)
	logger.Info().Str("Branch", *branch).Msg(indent)
	logger.Info().Str("Serving Config", *servingConfig).Msg(indent)
	logger.Info().Str("Parameter Input File", *parameterInputFile).Msg(indent)
	logger.Info().Int("Number of Results", *numberResults).Msg(indent)
	logger.Info().Str("Filter String", *filterString).Msg(indent)
	logger.Info().Msg("Begin")

	// Setup the Prediction Request
	var prediction = NewPrediction(*projectNumber, *location, *catalog, *branch, *servingConfig, *numberResults, *filterString)

	//  Load the Parameter Input File
	logger.Info().Msg("Loading Parameter Input File")
	err := prediction.LoadParameters(*parameterInputFile)
	if err != nil {
		logger.Error().Err(err).Msg("Load Parameter Input File Failed")
		os.Exit(1)
	}

	// Request to get a Retail API Prediction
	err = prediction.ExecuteRequests()
	if err != nil {
		logger.Error().Err(err).Msg("Prediction Request Failed")
		os.Exit(1)
	}
	logger.Info().Msg("End")

}
