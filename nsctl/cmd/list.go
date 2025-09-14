package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Kubernetes namespaces",
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("❌ Failed to load kubeconfig: %v\n", err)
			os.Exit(1)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Printf("❌ Failed to create Kubernetes client: %v\n", err)
			os.Exit(1)
		}

		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("❌ Failed to list namespaces: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tSTATUS\tAGE")

		for _, ns := range namespaces.Items {
			fmt.Fprintf(w, "%s\t%s\t%s\n", ns.Name, ns.Status.Phase, ns.CreationTimestamp.Time.Format("2006-01-02 15:04:05"))
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}