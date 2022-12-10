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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	retail "cloud.google.com/go/retail/apiv2"
	"cloud.google.com/go/retail/apiv2/retailpb"
)

type Prediction struct {
	placement     string
	numberResults int
	filter        string
	userEvents    []retailpb.UserEvent
}

//---------------------------------------------------------------------------------------

// Create a new Prediction struct populated
func NewPrediction(project string, location string, catalog string, servingConfig string, numberResults int, filter string) Prediction {
	return Prediction{
		placement:     fmt.Sprintf("projects/%s/locations/%s/catalogs/%s/servingConfigs/%s", project, location, catalog, servingConfig),
		numberResults: numberResults,
		filter:        filter,
	}
}

//---------------------------------------------------------------------------------------

// Execute the Requests to get a Retail API Prediction for each of the UserEvent objects
func (prediction *Prediction) ExecuteRequests() error {

	// Establish a Retail API Client
	logger.Info().Msg("Establishing a Retail Prediction Client")
	ctx := context.Background()
	client, err := retail.NewPredictionClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed Establishing a Retail Prediction Client: %w", err)
	}
	defer client.Close()

	for i := 0; i < len(prediction.userEvents); i++ {

		// Populate Request Parameters
		request := retailpb.PredictRequest{
			Placement:    prediction.placement,
			UserEvent:    &prediction.userEvents[i],
			PageSize:     int32(prediction.numberResults),
			Filter:       prediction.filter,
			ValidateOnly: false,
		}

		// Encode the User Event to send to the log
		rawUserEvent, err := json.Marshal(&prediction.userEvents[i])
		if err != nil {
			return fmt.Errorf("Encoding the User Event as JSON Failed: %w", err)
		}

		logger.Info().Int("Number", i+1).Msg("Initiating Prediction Request")
		logger.Info().RawJSON("Parameters", rawUserEvent).Msg(indent)

		// Raise the Prediction Request
		response, err := client.Predict(ctx, &request)
		if err != nil {
			return fmt.Errorf("Prediction Request Failed: %w", err)
		}

		// Encode the Results to send to the log
		rawResults, err := json.Marshal(&response.Results)
		if err != nil {
			return fmt.Errorf("Encoding the Response Results as JSON Failed: %w", err)
		}
		logger.Info().RawJSON("Results", rawResults).Msg(indent)

	}

	return nil
}

//---------------------------------------------------------------------------------------

// Walk the provided Parameter Input File GLOB and load all parameters
func (prediction *Prediction) LoadParameters(inputFile string) error {

	inputFile, _ = filepath.Abs(inputFile)
	buf, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("Reading the Parameter Input File Failed: %w", err)
	}

	err = json.Unmarshal(buf, &prediction.userEvents)
	if err != nil {
		return fmt.Errorf("Parsing the Parameter Input File Failed: %w", err)
	}

	return nil
}
