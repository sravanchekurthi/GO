package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var deleteName string

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Kubernetes namespace",
	Run: func(cmd *cobra.Command, args []string) {
		if deleteName == "" {
			fmt.Println("❌ Error: --name is required")
			os.Exit(1)
		}

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

		err = clientset.CoreV1().Namespaces().Delete(context.TODO(), deleteName, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("❌ Failed to delete namespace '%s': %v\n", deleteName, err)
			os.Exit(1)
		}

		fmt.Printf("🗑️  Namespace '%s' deleted successfully!\n", deleteName)
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&deleteName, "name", "n", "", "Name of the namespace to delete")
	rootCmd.AddCommand(deleteCmd)
}