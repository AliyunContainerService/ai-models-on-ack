/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler/common"
	"github.com/KunWuLuan/hostlist"
	"github.com/spf13/cobra"
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("resume called, args %v\n", args)
		resume(args[0])
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resumeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resumeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func resume(name string) {
	cls := os.Getenv("SLURM_CLUSTER")
	ns := os.Getenv("SLURM_NAMESPACE")
	fmt.Printf("slurm cluster metadata: %v %v\n", ns, cls)

	client := common.InitClientSet()
	cluster, err := client.KaiV1().SlurmClusters(ns).Get(context.Background(), cls, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("failed to get SlurmCluster CR: %v\n", err)
		os.Stdout.Sync()
		os.Exit(-1)
	}
	fmt.Printf("get SlurmCluster CR %v succeed\n", cls)

	hostlists, err := hostlist.ExpandHostlist(name, false, true)
	if err != nil {
		fmt.Printf("failed to expand hostlist: %v\n", name)
		os.Stdout.Sync()
		os.Exit(-1)
	}
	fmt.Printf("hostlists: %v\n", hostlists)
	for _, nodename := range hostlists {
		fmt.Printf("resume node %v\n", nodename)
		if !strings.HasPrefix(nodename, cls) {
			fmt.Printf("node name must be the format {$cluster_name}-{$group_name}-{$index}: %v\n", nodename)
			os.Stdout.Sync()
			return
			// os.Exit(-1)
		}
		nodename = strings.TrimPrefix(nodename, cls+"-worker")

		infos := strings.Split(nodename, "-")
		idxStr := infos[len(infos)-1]
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			fmt.Printf("failed to get node index %v: %v\n", idxStr, err)
			os.Stdout.Sync()
			os.Exit(-1)
		}
		groupName := strings.Trim(strings.TrimSuffix(nodename, idxStr), "-")
		found := false
		for i, workerSpec := range cluster.Spec.WorkerGroupSpecs {
			if workerSpec.GroupName == groupName {
				for _, w := range workerSpec.PodIndexesToBeCreated {
					if int32(w) == int32(idx) {
						fmt.Printf("the worker " + nodename + " is already in use\n")
						break
					}
				}
				fmt.Printf("resume worker %v\n", nodename)
				cluster.Spec.WorkerGroupSpecs[i].PodIndexesToBeCreated = append(cluster.Spec.WorkerGroupSpecs[i].PodIndexesToBeCreated, int32(idx))
				cluster.Spec.WorkerGroupSpecs[i].Replicas++
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("failed to find worker group %v in SlurmCluster CR\n", groupName)
			os.Stdout.Sync()
			os.Exit(-1)
		}
		fmt.Printf("resume node %v end\n", nodename)
	}
	if _, err := client.KaiV1().SlurmClusters(ns).Update(context.Background(), cluster, metav1.UpdateOptions{}); err != nil {
		fmt.Printf("failed to update SlurmCluster CR: %v\n", err)
		os.Stdout.Sync()
		os.Exit(-1)
	}
}
