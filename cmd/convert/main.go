package main

import (
	"context"
	"os"

	argocdv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/bigkevmcd/argo2flux/pkg/convert"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/yaml"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(argocdv1.AddToScheme(scheme))
}

func main() {
	cfg, err := config.GetConfig()
	cobra.CheckErr(err)

	cl, err := client.New(cfg, client.Options{Scheme: scheme})
	cobra.CheckErr(err)

	cobra.CheckErr(NewRootCommand(cl).Execute())
}

// Build the cobra command that handles our command line tool.
func NewRootCommand(cl client.Client) *cobra.Command {
	name := types.NamespacedName{}
	rootCmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert ArgoCD applications to Flux",
		Long:  `Convert ArgoCD applications to Flux resources automatically`,
		RunE: func(cmd *cobra.Command, args []string) error {
			objs, err := convert.ConvertToKustomization(context.Background(), cl, name)
			if err != nil {
				return err
			}
			if err := writeMultiDoc(objs); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.Flags().StringVar(&name.Name, "app-name", "", "name of the ArgoCD Application to convert")
	rootCmd.Flags().StringVar(&name.Namespace, "app-namespace", "", "namespace of the ArgoCD Application to convert")
	cobra.CheckErr(rootCmd.MarkFlagRequired("app-name"))
	cobra.CheckErr(rootCmd.MarkFlagRequired("app-namespace"))

	return rootCmd
}

func writeMultiDoc(objs []runtime.Object) error {
	for _, obj := range objs {
		if _, err := os.Stdout.Write([]byte("---\n")); err != nil {
			return err
		}
		b, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		if _, err := os.Stdout.Write(b); err != nil {
			return err
		}
	}
	return nil
}
