// Package elasticcontainer is responsible for creating and managing an Elasticsearch container for testing purposes.
package elasticcontainer

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/testcontainers/testcontainers-go"
	esmodule "github.com/testcontainers/testcontainers-go/modules/elasticsearch"
)

//go:embed vss_mapping.json
var vssMapping []byte

// Create function starts and elastic container then uses datagen to generate vehicle data and insert it into the elastic container.
func Create(ctx context.Context) (esmodule.Options, func(), error) {
	esContainer, err := esmodule.RunContainer(
		ctx,
		testcontainers.WithImage("docker.elastic.co/elasticsearch/elasticsearch:8.3.0"),
		esmodule.WithPassword("testpassword"),
	)
	cleanup := func() {
		if err := esContainer.Terminate(ctx); err != nil {
			log.Fatalf("Could not terminate Elasticsearch container: %v", err)
		}
	}
	if err != nil {
		return esmodule.Options{}, cleanup, fmt.Errorf("Could not start Elasticsearch container: %v", err)
	}

	return esContainer.Settings, cleanup, nil
}

// AddVSSMapping inserts a mapping into Elasticsearch using a typed client.
func AddVSSMapping(ctx context.Context, client *elasticsearch.TypedClient, index string) error {
	_, err := client.Indices.Create(index).WaitForActiveShards("1").Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	_, err = client.Indices.PutMapping(index).Raw(bytes.NewReader(vssMapping)).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert mapping: %w", err)
	}

	return nil
}
