// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
// --------------------------------------------------------------------------------------------
// Generated file, DO NOT EDIT
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// --------------------------------------------------------------------------------------------

package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var ResourceAreaId, _ = uuid.Parse("94ff054d-5ee1-413d-9341-3f4a7827de2e")

type Client interface {
	// [Preview API] Create new Audit Stream
	CreateStream(context.Context, CreateStreamArgs) (*AuditStream, error)
	// [Preview API] Delete Audit Stream
	DeleteStream(context.Context, DeleteStreamArgs) error
	// [Preview API] Downloads audit log entries.
	DownloadLog(context.Context, DownloadLogArgs) (io.ReadCloser, error)
	// [Preview API] Get all auditable actions filterable by area.
	GetActions(context.Context, GetActionsArgs) (*[]AuditActionInfo, error)
	// [Preview API] Return all Audit Streams scoped to an organization
	QueryAllStreams(context.Context, QueryAllStreamsArgs) (*[]AuditStream, error)
	// [Preview API] Queries audit log entries
	QueryLog(context.Context, QueryLogArgs) (*AuditLogQueryResult, error)
	// [Preview API] Return Audit Stream with id of streamId if one exists otherwise throw
	QueryStreamById(context.Context, QueryStreamByIdArgs) (*AuditStream, error)
	// [Preview API] Update existing Audit Stream status
	UpdateStatus(context.Context, UpdateStatusArgs) (*AuditStream, error)
	// [Preview API] Update existing Audit Stream
	UpdateStream(context.Context, UpdateStreamArgs) (*AuditStream, error)
}

type ClientImpl struct {
	Client azuredevops.Client
}

func NewClient(ctx context.Context, connection *azuredevops.Connection) (Client, error) {
	client, err := connection.GetClientByResourceAreaId(ctx, ResourceAreaId)
	if err != nil {
		return nil, err
	}
	return &ClientImpl{
		Client: *client,
	}, nil
}

// [Preview API] Create new Audit Stream
func (client *ClientImpl) CreateStream(ctx context.Context, args CreateStreamArgs) (*AuditStream, error) {
	if args.Stream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Stream"}
	}
	queryParams := url.Values{}
	if args.DaysToBackfill == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "daysToBackfill"}
	}
	queryParams.Add("daysToBackfill", strconv.Itoa(*args.DaysToBackfill))
	body, marshalErr := json.Marshal(*args.Stream)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("77d60bf9-1882-41c5-a90d-3a6d3c13fd3b")
	resp, err := client.Client.Send(ctx, http.MethodPost, locationId, "6.0-preview.1", nil, queryParams, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AuditStream
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the CreateStream function
type CreateStreamArgs struct {
	// (required) Stream entry
	Stream *AuditStream
	// (required) The number of days of previously recorded audit data that will be replayed into the stream. A value of zero will result in only new events being streamed.
	DaysToBackfill *int
}

// [Preview API] Delete Audit Stream
func (client *ClientImpl) DeleteStream(ctx context.Context, args DeleteStreamArgs) error {
	routeValues := make(map[string]string)
	if args.StreamId == nil {
		return &azuredevops.ArgumentNilError{ArgumentName: "args.StreamId"}
	}
	routeValues["streamId"] = strconv.Itoa(*args.StreamId)

	locationId, _ := uuid.Parse("77d60bf9-1882-41c5-a90d-3a6d3c13fd3b")
	_, err := client.Client.Send(ctx, http.MethodDelete, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return err
	}

	return nil
}

// Arguments for the DeleteStream function
type DeleteStreamArgs struct {
	// (required) Id of stream entry to delete
	StreamId *int
}

// [Preview API] Downloads audit log entries.
func (client *ClientImpl) DownloadLog(ctx context.Context, args DownloadLogArgs) (io.ReadCloser, error) {
	queryParams := url.Values{}
	if args.Format == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "format"}
	}
	queryParams.Add("format", *args.Format)
	if args.StartTime != nil {
		queryParams.Add("startTime", (*args.StartTime).String())
	}
	if args.EndTime != nil {
		queryParams.Add("endTime", (*args.EndTime).String())
	}
	locationId, _ := uuid.Parse("b7b98a76-04e8-4f4d-ac72-9d46492caaac")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", nil, queryParams, nil, "", "application/octet-stream", nil)
	if err != nil {
		return nil, err
	}

	return resp.Body, err
}

// Arguments for the DownloadLog function
type DownloadLogArgs struct {
	// (required) File format for download. Can be "json" or "csv".
	Format *string
	// (optional) Start time of download window. Optional
	StartTime *azuredevops.Time
	// (optional) End time of download window. Optional
	EndTime *azuredevops.Time
}

// [Preview API] Get all auditable actions filterable by area.
func (client *ClientImpl) GetActions(ctx context.Context, args GetActionsArgs) (*[]AuditActionInfo, error) {
	queryParams := url.Values{}
	if args.AreaName != nil {
		queryParams.Add("areaName", *args.AreaName)
	}
	locationId, _ := uuid.Parse("6fa30b9a-9558-4e3b-a95f-a12572caa6e6")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []AuditActionInfo
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the GetActions function
type GetActionsArgs struct {
	// (optional) Optional. Get actions scoped to area
	AreaName *string
}

// [Preview API] Return all Audit Streams scoped to an organization
func (client *ClientImpl) QueryAllStreams(ctx context.Context, args QueryAllStreamsArgs) (*[]AuditStream, error) {
	locationId, _ := uuid.Parse("77d60bf9-1882-41c5-a90d-3a6d3c13fd3b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", nil, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue []AuditStream
	err = client.Client.UnmarshalCollectionBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryAllStreams function
type QueryAllStreamsArgs struct {
}

// [Preview API] Queries audit log entries
func (client *ClientImpl) QueryLog(ctx context.Context, args QueryLogArgs) (*AuditLogQueryResult, error) {
	queryParams := url.Values{}
	if args.StartTime != nil {
		queryParams.Add("startTime", (*args.StartTime).String())
	}
	if args.EndTime != nil {
		queryParams.Add("endTime", (*args.EndTime).String())
	}
	if args.BatchSize != nil {
		queryParams.Add("batchSize", strconv.Itoa(*args.BatchSize))
	}
	if args.ContinuationToken != nil {
		queryParams.Add("continuationToken", *args.ContinuationToken)
	}
	if args.SkipAggregation != nil {
		queryParams.Add("skipAggregation", strconv.FormatBool(*args.SkipAggregation))
	}
	locationId, _ := uuid.Parse("4e5fa14f-7097-4b73-9c85-00abc7353c61")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", nil, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AuditLogQueryResult
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryLog function
type QueryLogArgs struct {
	// (optional) Start time of download window. Optional
	StartTime *azuredevops.Time
	// (optional) End time of download window. Optional
	EndTime *azuredevops.Time
	// (optional) Max number of results to return. Optional
	BatchSize *int
	// (optional) Token used for returning next set of results from previous query. Optional
	ContinuationToken *string
	// (optional) Skips aggregating events and leaves them as individual entries instead. By default events are aggregated. Event types that are aggregated: AuditLog.AccessLog.
	SkipAggregation *bool
}

// [Preview API] Return Audit Stream with id of streamId if one exists otherwise throw
func (client *ClientImpl) QueryStreamById(ctx context.Context, args QueryStreamByIdArgs) (*AuditStream, error) {
	routeValues := make(map[string]string)
	if args.StreamId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.StreamId"}
	}
	routeValues["streamId"] = strconv.Itoa(*args.StreamId)

	locationId, _ := uuid.Parse("77d60bf9-1882-41c5-a90d-3a6d3c13fd3b")
	resp, err := client.Client.Send(ctx, http.MethodGet, locationId, "6.0-preview.1", routeValues, nil, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AuditStream
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the QueryStreamById function
type QueryStreamByIdArgs struct {
	// (required) Id of stream entry to retrieve
	StreamId *int
}

// [Preview API] Update existing Audit Stream status
func (client *ClientImpl) UpdateStatus(ctx context.Context, args UpdateStatusArgs) (*AuditStream, error) {
	routeValues := make(map[string]string)
	if args.StreamId == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.StreamId"}
	}
	routeValues["streamId"] = strconv.Itoa(*args.StreamId)

	queryParams := url.Values{}
	if args.Status == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "status"}
	}
	queryParams.Add("status", string(*args.Status))
	locationId, _ := uuid.Parse("77d60bf9-1882-41c5-a90d-3a6d3c13fd3b")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "6.0-preview.1", routeValues, queryParams, nil, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AuditStream
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateStatus function
type UpdateStatusArgs struct {
	// (required) Id of stream entry to be updated
	StreamId *int
	// (required) Status of the stream
	Status *AuditStreamStatus
}

// [Preview API] Update existing Audit Stream
func (client *ClientImpl) UpdateStream(ctx context.Context, args UpdateStreamArgs) (*AuditStream, error) {
	if args.Stream == nil {
		return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.Stream"}
	}
	body, marshalErr := json.Marshal(*args.Stream)
	if marshalErr != nil {
		return nil, marshalErr
	}
	locationId, _ := uuid.Parse("77d60bf9-1882-41c5-a90d-3a6d3c13fd3b")
	resp, err := client.Client.Send(ctx, http.MethodPut, locationId, "6.0-preview.1", nil, nil, bytes.NewReader(body), "application/json", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var responseValue AuditStream
	err = client.Client.UnmarshalBody(resp, &responseValue)
	return &responseValue, err
}

// Arguments for the UpdateStream function
type UpdateStreamArgs struct {
	// (required) Stream entry
	Stream *AuditStream
}
