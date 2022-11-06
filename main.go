package main

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"
)

const LoopInterval = 15 * time.Second

func main() {
	ctx := context.Background()

	db, err := newStorage("badger")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *storage) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	clientset, err := getClientset()
	if err != nil {
		panic(err.Error())
	}

	statmap := NewStatMap(db)

	loop(ctx, clientset, &statmap)
}

func loop(ctx context.Context, clientset *kubernetes.Clientset, statmap *StatMap) {
	ticker := time.Tick(LoopInterval)
	for ; true; <-ticker { // interesting hack to start right away
		if err := loopStep(ctx, clientset, statmap); err != nil {
			panic(err.Error())
		}
		statmap.Print()
	}
}

func loopStep(ctx context.Context, clientset *kubernetes.Clientset, statmap *StatMap) error {
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error getting nodes. %w", err)
	}

	for _, node := range nodes.Items {
		summary, err := getNodeSummary(ctx, clientset, node.Name)
		if err != nil {
			return fmt.Errorf("error getting node summary. %w", err)
		}
		populateStatMap(statmap, summary)
	}
	return nil
}

func populateStatMap(s *StatMap, summary nodeSummary) {
	for _, pod := range summary.Pods {
		s.Add(pod.PodRef.Name, pod.Memory.WorkingSetBytes)
	}
}

func getNodeSummary(ctx context.Context, clientset *kubernetes.Clientset, node string) (nodeSummary, error) {
	path := fmt.Sprintf("/api/v1/nodes/%s/proxy/stats/summary", node)

	result := clientset.RESTClient().Get().AbsPath(path).Do(ctx)
	if result.Error() != nil {
		return nodeSummary{}, fmt.Errorf("error getting node summary. %w", result.Error())
	}

	data, err := result.Raw()
	if err != nil {
		return nodeSummary{}, fmt.Errorf("error getting raw data. %w", err)
	}

	var summary nodeSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		return nodeSummary{}, fmt.Errorf("error unmarshalling data. %w", err)
	}
	return summary, nil
}
