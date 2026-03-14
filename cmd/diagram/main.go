package main

import (
	"context"
	"fmt"
	"os"

	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrel"
	"oss.terrastruct.com/d2/lib/log"
)

// GenerateClusterDiagram takes GKE metadata and produces a D2 script.
func GenerateClusterDiagram(nodes []string, services []string, outputPath string) error {
	ctx := context.Background()
	graph := d2graph.NewGraph()

	// 1. Define the Global Cloud Container
	gcp, _ := graph.AddNode(ctx, "GCP_Project", nil)
	gcp.Label = "Google Cloud Platform (Stash-House)"

	// 2. Define the VPC and Cluster subgraphs
	vpc, _ := graph.AddNode(ctx, "GCP_Project.VPC", nil)
	vpc.Label = "stash-house-vpc"

	cluster, _ := graph.AddNode(ctx, "GCP_Project.VPC.GKE", nil)
	cluster.Label = "GKE Cluster: stash-house-poc"

	// 3. Programmatically add Nodes from the provided list
	for _, nodeName := range nodes {
		node, _ := graph.AddNode(ctx, fmt.Sprintf("GCP_Project.VPC.GKE.%s", nodeName), nil)
		node.Label = fmt.Sprintf("Compute Node: %s", nodeName)
	}

	// 4. Programmatically add Services/Pods
	for _, svcName := range services {
		svc, _ := graph.AddNode(ctx, fmt.Sprintf("GCP_Project.VPC.GKE.%s_Pod", svcName), nil)
		svc.Label = fmt.Sprintf("Service: %s", svcName)
		
		// Style services differently (Blue stroke for GKE objects)
		svc.Values.Stroke = "#4285F4"
	}

	// 5. Render the AST to a D2 DSL string
	out := d2format.Format(graph.AST)

	// 6. Write to disk
	err := os.WriteFile(outputPath, []byte(out), 0644)
	if err != nil {
		return fmt.Errorf("failed to write D2 file: %w", err)
	}

	fmt.Printf("Successfully generated diagram definition at %s\n", outputPath)
	return nil
}

func main() {
	nodes := []string{"gke-node-01", "gke-node-02"}
	services := []string{"Hornet", "FOKS-Server", "Stash-Daemon"}
	
	err := GenerateClusterDiagram(nodes, services, "cluster_arch.d2")
	if err != nil {
		panic(err)
	}
}
