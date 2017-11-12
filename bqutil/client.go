package bqutil

import (
	"context"
	"sync"

	"cloud.google.com/go/bigquery"
	"github.com/tsaikd/KDGoLib/errutil"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// errors
var (
	ErrUnknownOption2 = errutil.NewFactory("Unknown option type: %T , value: %+v")
)

// Client is instance of bigquery client
type Client struct {
	ctx            context.Context
	client         *bigquery.Client
	msgQueue       chan uploadMessage
	ensureDataset  OptionEnsureDatasetWhenFirstUpload
	ensureTable    OptionEnsureTableWhenFirstUpload
	ensureDatasets sync.Map
	ensureTables   sync.Map
}

// OptionUploaderQueueSize queue size for background uploader, default: 100
type OptionUploaderQueueSize uint64

// OptionUploaderQueueSizeDefault default value of OptionUploaderQueueSize
const OptionUploaderQueueSizeDefault = OptionUploaderQueueSize(100)

// OptionUploaderErrorHandler error handler when upload error in background uploader,
// default: errutil.Trace(err)
type OptionUploaderErrorHandler func(error)

// OptionUploaderErrorHandlerDefault default value of OptionUploaderErrorHandler
var OptionUploaderErrorHandlerDefault = OptionUploaderErrorHandler(errutil.Trace)

// OptionEnsureDatasetWhenFirstUpload ensure dataset when first upload occur
type OptionEnsureDatasetWhenFirstUpload bool

// OptionEnsureDatasetWhenFirstUploadDefault default value of OptionEnsureDatasetWhenFirstUpload
const OptionEnsureDatasetWhenFirstUploadDefault = OptionEnsureDatasetWhenFirstUpload(true)

// OptionEnsureTableWhenFirstUpload ensure table when first upload occur
type OptionEnsureTableWhenFirstUpload bool

// OptionEnsureTableWhenFirstUploadDefault default value of OptionEnsureTableWhenFirstUpload
const OptionEnsureTableWhenFirstUploadDefault = OptionEnsureTableWhenFirstUpload(true)

// NewClient return Client with Config
func NewClient(
	ctx context.Context,
	gcKeyFile string,
	projectID string,
	options ...interface{},
) (*Client, error) {
	queueSize := OptionUploaderQueueSizeDefault
	errorHandler := OptionUploaderErrorHandlerDefault
	ensureDataset := OptionEnsureDatasetWhenFirstUploadDefault
	ensureTable := OptionEnsureTableWhenFirstUploadDefault
	for _, option := range options {
		switch opt := option.(type) {
		case OptionUploaderQueueSize:
			queueSize = opt
		case OptionUploaderErrorHandler:
			errorHandler = opt
		case OptionEnsureDatasetWhenFirstUpload:
			ensureDataset = opt
		case OptionEnsureTableWhenFirstUpload:
			ensureTable = opt
		default:
			return nil, ErrUnknownOption2.New(nil, option, opt)
		}
	}

	client, err := bigquery.NewClient(
		ctx,
		projectID,
		option.WithServiceAccountFile(gcKeyFile),
	)
	if err != nil {
		return nil, err
	}

	msgQueue := make(chan uploadMessage, queueSize)
	go startUploaderLoop(ctx, msgQueue, errorHandler)

	return &Client{
		ctx:           ctx,
		client:        client,
		msgQueue:      msgQueue,
		ensureDataset: ensureDataset,
		ensureTable:   ensureTable,
	}, nil
}

// Close bigquery client
func (t *Client) Close() error {
	return t.client.Close()
}

// EnsureDataset ensure dataset exists in project
func (t *Client) EnsureDataset(datasetName string) (*bigquery.Dataset, error) {
	ctx := t.ctx
	dataset := t.client.Dataset(datasetName)
	if _, err := dataset.Metadata(ctx); err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return dataset, dataset.Create(ctx, &bigquery.DatasetMetadata{})
		}
		return dataset, err
	}
	return dataset, nil
}

// EnsureTable create bigquery table with schema if not exists
func (t *Client) EnsureTable(datasetName string, tableName string, schemaObject interface{}) (*bigquery.Table, error) {
	ctx := t.ctx
	dataset := t.client.Dataset(datasetName)
	table := dataset.Table(tableName)

	schema, err := bigquery.InferSchema(schemaObject)
	if err != nil {
		return table, err
	}
	setSchemaOptional(&schema)

	if _, err = table.Metadata(ctx); err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			return table, table.Create(ctx, &bigquery.TableMetadata{Schema: schema})
		}
		return table, err
	}
	return table, nil
}

func setSchemaOptional(schema *bigquery.Schema) {
	if schema != nil {
		for _, field := range *schema {
			field.Required = false
			setSchemaOptional(&field.Schema)
		}
	}
}

type uploadMessage struct {
	table *bigquery.Table
	data  interface{}
}

// UploadAsync data to bigquery table async
func (t *Client) UploadAsync(datasetName string, tableName string, data interface{}) (err error) {
	if t.ensureDataset {
		if _, exists := t.ensureDatasets.LoadOrStore(datasetName, true); !exists {
			if _, err = t.EnsureDataset(datasetName); err != nil {
				return
			}
		}
	}
	if t.ensureTable {
		key := datasetName + "/" + tableName
		if _, exists := t.ensureTables.LoadOrStore(key, true); !exists {
			if _, err = t.EnsureTable(datasetName, tableName, data); err != nil {
				return
			}
		}
	}

	dataset := t.client.Dataset(datasetName)
	table := dataset.Table(tableName)
	t.msgQueue <- uploadMessage{
		table: table,
		data:  data,
	}
	return nil
}

func startUploaderLoop(ctx context.Context, msgQueue <-chan uploadMessage, errorHandler OptionUploaderErrorHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgQueue:
			if err := msg.table.Uploader().Put(ctx, msg.data); err != nil {
				errorHandler(err)
			}
		}
	}
}
