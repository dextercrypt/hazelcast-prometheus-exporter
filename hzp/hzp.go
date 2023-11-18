package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hazelcast/hazelcast-go-client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	hazelcastMapSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hazelcast_map_size",
			Help: "Current size of the Hazelcast map",
		},
		[]string{"map_name"},
	)
)

func main() {
	// Hazelcast Configuration
	config := hazelcast.NewConfig()
	config.Cluster.Network.SetAddresses("hazelcast-service:5701") // Update with your Hazelcast service name and port

	// Create Hazelcast Client
	client, err := hazelcast.StartNewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Register Prometheus Metrics
	prometheus.MustRegister(hazelcastMapSize)

	// Start Hazelcast Metrics Collection
	go collectHazelcastMetrics(client)

	// Expose Prometheus metrics on /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	port := 8080 // Update with your desired port
	log.Printf("Server listening on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func collectHazelcastMetrics(client *hazelcast.Client) {
	for {
		// Get Hazelcast Map
		mapName := "exampleMap" // Update with your Hazelcast map name
		mapStats, err := client.GetMap(context.Background(), mapName)
		if err != nil {
			log.Println("Error getting Hazelcast map:", err)
		} else {
			// Get Map Size and update Prometheus metric
			size, err := mapStats.Size(context.Background())
			if err != nil {
				log.Println("Error getting Hazelcast map size:", err)
			} else {
				hazelcastMapSize.WithLabelValues(mapName).Set(float64(size))
			}
		}

		// Sleep for 1 minute before collecting metrics again
		time.Sleep(time.Minute)
	}
}

