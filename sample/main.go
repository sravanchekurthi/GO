package main

import (
    "fmt"
    "log"
    "os"
    "text/template"

    "github.com/spf13/cobra"
)

var namespace string
var imageName string

func main() {
    var rootCmd = &cobra.Command{
        Use:   "nscli",
        Short: "Namespace automation CLI",
    }

    var createCmd = &cobra.Command{
        Use:   "create",
        Short: "Create namespace, Dockerfile, and GitHub Action",
        Run: func(cmd *cobra.Command, args []string) {
            if namespace == "" || imageName == "" {
                log.Fatal("Namespace and image name are required")
            }
            createNamespace(namespace)
            createDockerfile(imageName)
            createGitHubAction(namespace, imageName)
            fmt.Println("✅ All files created successfully!")
        },
    }

    createCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace name")
    createCmd.Flags().StringVarP(&imageName, "image", "i", "", "Docker image name")

    rootCmd.AddCommand(createCmd)
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func createNamespace(ns string) {
    content := fmt.Sprintf("apiVersion: v1\nkind: Namespace\nmetadata:\n  name: %s\n", ns)
    err := os.WriteFile("namespace.yaml", []byte(content), 0644)
    if err != nil {
        log.Fatalf("Error creating namespace.yaml: %v", err)
    }
    fmt.Println("📄 namespace.yaml created")
}

func createDockerfile(image string) {
    tmpl, err := template.ParseFiles("templates/Dockerfile.tmpl")
    if err != nil {
        log.Fatalf("Error reading Dockerfile template: %v", err)
    }
    f, err := os.Create("Dockerfile")
    if err != nil {
        log.Fatalf("Error creating Dockerfile: %v", err)
    }
    defer f.Close()
    tmpl.Execute(f, map[string]string{"ImageName": image})
    fmt.Println("📄 Dockerfile created")
}

func createGitHubAction(ns, image string) {
    tmpl, err := template.ParseFiles("templates/github-action.yml.tmpl")
    if err != nil {
        log.Fatalf("Error reading GitHub Action template: %v", err)
    }
    os.MkdirAll(".github/workflows", 0755)
    f, err := os.Create(".github/workflows/deploy.yml")
    if err != nil {
        log.Fatalf("Error creating GitHub Action file: %v", err)
    }
    defer f.Close()
    tmpl.Execute(f, map[string]string{"Namespace": ns, "ImageName": image})
    fmt.Println("📄 GitHub Action workflow created")
}