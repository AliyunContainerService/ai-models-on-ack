/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/AliyunContainerService/ai-models-on-ack/framework/slurm/slurmcluster-scaler/common"
	"github.com/KunWuLuan/hostlist"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// suspendCmd represents the suspend command
var suspendCmd = &cobra.Command{
	Use:   "suspend",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("suspend called")
		suspend(args[0])
	},
}

func init() {
	rootCmd.AddCommand(suspendCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// suspendCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// suspendCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func suspend(name string) {
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
		fmt.Printf("suspend node %v", nodename)
		if !strings.HasPrefix(nodename, cls) {
			fmt.Printf("node name must be the format {$cluster_name}-{$group_name}-{$index}: %v\n", nodename)
			os.Stdout.Sync()
			os.Exit(-1)
		}

		nodenameWithoutClsName := strings.TrimPrefix(nodename, cls+"-worker")
		infos := strings.Split(nodenameWithoutClsName, "-")
		idxStr := infos[len(infos)-1]
		groupName := strings.Trim(strings.TrimSuffix(nodenameWithoutClsName, idxStr), "-")
		found := false
		for i, workerSpec := range cluster.Spec.WorkerGroupSpecs {
			if workerSpec.GroupName == groupName {
				for _, w := range workerSpec.PodsToBeDeleted {
					if w == nodename {
						fmt.Printf("the worker " + nodename + " is already deleted\n")
						break
					}
				}
				fmt.Printf("suspend worker %v\n", nodename)
				cluster.Spec.WorkerGroupSpecs[i].PodsToBeDeleted = append(cluster.Spec.WorkerGroupSpecs[i].PodsToBeDeleted, nodename)
				cluster.Spec.WorkerGroupSpecs[i].Replicas--
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("failed to find worker group %v in SlurmCluster CR\n", groupName)
			os.Stdout.Sync()
			os.Exit(-1)
		}
		fmt.Printf("suspend node %v end", nodename)
	}
	if _, err := client.KaiV1().SlurmClusters(ns).Update(context.Background(), cluster, metav1.UpdateOptions{}); err != nil {
		fmt.Printf("failed to update SlurmCluster CR: %v\n", err)
		os.Stdout.Sync()
		os.Exit(-1)
	}
}
