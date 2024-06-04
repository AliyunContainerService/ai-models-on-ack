/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler/common"
	"github.com/KunWuLuan/hostlist"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("watch called")
		watch()
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func watch() {
	cls := os.Getenv("SLURM_CLUSTER")
	ns := os.Getenv("SLURM_NAMESPACE")

	client := common.InitClientSet()
	nodes := sets.NewString()
	nodesToBeDeleted := map[string]int{}

	for {
		cluster, err := client.KaiV1().SlurmClusters(ns).Get(context.Background(), cls, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("failed to get cluster, err: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		nodesToDelete := sets.NewString()
		lines := strings.Split(string(NodesData()), "\n")
		for _, line := range lines {
			nodesInLine, err := hostlist.ExpandHostlist(line, false, true)
			if err != nil {
				fmt.Printf("failed to expand hostlist %v, err: %v\n", line, err)
				continue
			}
			nodes.Insert(nodesInLine...)
		}
		availableNodes := sets.NewString(cluster.Status.AvailableWorkers...)
		if z := availableNodes.Difference(nodes); z.Len() > 0 {
			for n := range z {
				if _, ok := nodesToBeDeleted[n]; !ok {
					nodesToBeDeleted[n] = 5
				} else {
					nodesToBeDeleted[n]--
					if nodesToBeDeleted[n] <= 0 {
						nodesToDelete.Insert(n)
					}
				}
			}
		}

		if nodesToDelete.Len() > 0 {
			nodesPerGroup := map[string][]string{}
			for nodename := range nodesToDelete {
				if !strings.HasPrefix(nodename, cls) {
					fmt.Printf("node name must be the format {$cluster_name}-{$group_name}-{$index}, found %v\n", nodename)
					continue
				}
				nodename = strings.TrimPrefix(nodename, cls)
				infos := strings.Split(nodename, "-")
				idxStr := infos[len(infos)-1]
				groupName := strings.Trim(strings.TrimSuffix(nodename, idxStr), "-")
				if _, ok := nodesPerGroup[groupName]; ok {
					nodesPerGroup[groupName] = append(nodesPerGroup[groupName], nodename)
				} else {
					nodesPerGroup[groupName] = []string{nodename}
				}
			}
			needUpdate := false
			for i, spec := range cluster.Spec.WorkerGroupSpecs {
				if nodesInGroup, ok := nodesPerGroup[spec.GroupName]; ok {
					for _, node := range nodesInGroup {
						fmt.Printf("try to delete node %v\n", node)
						updated := false
						updated, cluster.Spec.WorkerGroupSpecs[i].PodsToBeDeleted = common.AddIfNotPresent(cluster.Spec.WorkerGroupSpecs[i].PodsToBeDeleted, node)
						needUpdate = needUpdate || updated
					}
				}
			}
			if needUpdate {
				_, err = client.KaiV1().SlurmClusters(ns).Update(context.Background(), cluster, metav1.UpdateOptions{})
				if err != nil {
					fmt.Printf("failed to update cluster, err: %v\n", err)
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

// Execute the sinfo command and return its output
func NodesData() []byte {
	cmd := exec.Command("sinfo", "-h", "-o %N")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	out, _ := ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return out
}
