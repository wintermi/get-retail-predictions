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
	Project       string
	Location      string
	Catalog       string
	Branch        string
	ServingConfig string
	PageSize      int32
	Filter        string
	UserEvents    []*retailpb.UserEvent
}

type Results struct {
	Id    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

//---------------------------------------------------------------------------------------

// Create a new Prediction struct populated
func NewPrediction(project string, location string, catalog string, branch string, servingConfig string, numberResults int, filter string) Prediction {
	return Prediction{
		Project:       project,
		Location:      location,
		Catalog:       catalog,
		Branch:        branch,
		ServingConfig: servingConfig,
		PageSize:      int32(numberResults),
		Filter:        filter,
	}
}

//---------------------------------------------------------------------------------------

// Execute the Requests to get a Retail API Prediction for each of the UserEvent objects
func (prediction *Prediction) ExecuteRequests() error {

	// Establish a Retail API Prediction Client
	logger.Info().Msg("Establishing a Retail Prediction Client")
	ctx := context.Background()
	client, err := retail.NewPredictionClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed Establishing a Retail Prediction Client: %w", err)
	}
	defer client.Close()

	placement := fmt.Sprintf("projects/%s/locations/%s/catalogs/%s/servingConfigs/%s",
		prediction.Project, prediction.Location, prediction.Catalog, prediction.ServingConfig)

	// Iterate through the User Events
	for i := 0; i < len(prediction.UserEvents); i++ {

		// Populate Request Parameters
		request := retailpb.PredictRequest{
			Placement:    placement,
			UserEvent:    prediction.UserEvents[i],
			PageSize:     prediction.PageSize,
			Filter:       prediction.Filter,
			ValidateOnly: false,
		}

		// Encode the User Event before sending to the log
		rawUserEvent, err := json.Marshal(&prediction.UserEvents[i])
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

		// Iterate through the results and get the Product Title
		for r := 0; r < len(response.Results); r++ {

			// Get the Product Title for the Product Id from the response
			title, err := prediction.GetProductTitle(response.Results[r].Id)
			if err != nil {
				return fmt.Errorf("Failed to Get Product Title: %w", err)
			}

			results := Results{
				Id:    response.Results[r].Id,
				Title: *title,
			}

			// Encode the Response Result before sending to the log
			rawResults, err := json.Marshal(&results)
			if err != nil {
				return fmt.Errorf("Encoding the Response Results as JSON Failed: %w", err)
			}
			logger.Info().RawJSON("Results", rawResults).Msg(indent)

		}

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

	err = json.Unmarshal(buf, &prediction.UserEvents)
	if err != nil {
		return fmt.Errorf("Parsing the Parameter Input File Failed: %w", err)
	}

	return nil
}

//---------------------------------------------------------------------------------------

// Execute a Requests to get a Product Title from the Retail API
func (prediction *Prediction) GetProductTitle(id string) (*string, error) {

	// Establish a Retail API Product Client
	ctx := context.Background()
	client, err := retail.NewProductClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed Establishing a Retail Product Client: %w", err)
	}
	defer client.Close()

	// Populate Request Parameters
	request := retailpb.GetProductRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/catalogs/%s/branches/%s/products/%s",
			prediction.Project, prediction.Location, prediction.Catalog, prediction.Branch, id),
	}

	// Raise the Prediction Request
	response, err := client.GetProduct(ctx, &request)
	if err != nil {
		return nil, fmt.Errorf("Prediction Request Failed: %w", err)
	}
	title := response.Title

	return &title, nil
}
