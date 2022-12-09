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
	"fmt"

	retail "cloud.google.com/go/retail/apiv2"
	"cloud.google.com/go/retail/apiv2/retailpb"
)

type Prediction struct {
	project       string
	location      string
	catalog       string
	servingConfig string
	numberResults int
	eventType     string
	visitor       string
	product       string
	filter        string
	experiment    string
}

//---------------------------------------------------------------------------------------

// Execute the Request to get a Retail API Prediction
func (prediction *Prediction) ExecuteRequest() error {

	// Establish a Retail API Client
	logger.Info().Msg("Establishing a Retail Prediction Client")
	ctx := context.Background()
	client, err := retail.NewPredictionClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed Establishing a Retail Prediction Client: %w", err)
	}
	defer client.Close()

	// Populate Request Parameters
	request := retailpb.PredictRequest{
		Placement: fmt.Sprintf("projects/%s/locations/%s/catalogs/%s/servingConfigs/%s", prediction.project, prediction.location, prediction.catalog, prediction.servingConfig),
		UserEvent: &retailpb.UserEvent{
			EventType:      prediction.eventType,
			VisitorId:      prediction.visitor,
			ExperimentIds:  prediction.GetExperiments(),
			ProductDetails: prediction.GetProductDetails(),
		},
		PageSize:     int32(prediction.numberResults),
		Filter:       prediction.filter,
		ValidateOnly: false,
	}

	logger.Info().Msg("Requesting for Predictions")
	response, err := client.Predict(ctx, &request)
	if err != nil {
		return fmt.Errorf("Prediction Request Failed: %w", err)
	}

	for i, r := range response.Results {
		logger.Info().Int("Number", i).Msg("Prediction Result")
		logger.Info().Str("Product Id", r.Id).Msg(indent)
	}

	return nil
}

//---------------------------------------------------------------------------------------

// Return ExperimentsId String Array if populated
func (prediction *Prediction) GetExperiments() []string {
	if prediction.experiment == "" {
		return nil
	}

	return []string{prediction.experiment}
}

//---------------------------------------------------------------------------------------

// Return ExperimentsId String Array if populated
func (prediction *Prediction) GetProductDetails() []*retailpb.ProductDetail {
	if prediction.product == "" {
		return nil
	}

	pd := &retailpb.ProductDetail{
		Product: &retailpb.Product{
			Id: prediction.product,
		},
	}

	return []*retailpb.ProductDetail{pd}
}
