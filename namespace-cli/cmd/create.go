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

var namespace string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Kubernetes namespace",
	Long:  `This command creates a new Kubernetes namespace using the default kubeconfig.`,
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == "" {
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

		ns := &metav1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}

		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("❌ Failed to create namespace '%s': %v\n", namespace, err)
			os.Exit(1)
		}

		fmt.Printf("✅ Namespace '%s' created successfully!\n", namespace)
	},
}

func init() {
	createCmd.Flags().StringVarP(&namespace, "name", "n", "", "Name of the namespace to create")
	rootCmd.AddCommand(createCmd)
}