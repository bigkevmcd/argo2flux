package convert

import (
	"context"
	"fmt"
	"time"

	argocdv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const defaultInterval time.Duration = time.Second * 600

// ConvertToKustomiation takes an ArgoCD Application and outputs a set of
// resources that will deploy it from Flux.
func ConvertToKustomization(ctx context.Context, cl client.Client, name types.NamespacedName) ([]runtime.Object, error) {
	app := &argocdv1.Application{}
	if err := cl.Get(ctx, name, app); err != nil {
		return nil, fmt.Errorf("failed to get Application %s: %w", name, err)
	}

	result := []runtime.Object{
		&sourcev1.GitRepository{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "source.toolkit.fluxcd.io/v1beta2",
				Kind:       "GitRepository",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.GetName(),
				Namespace: app.GetNamespace(),
			},
			Spec: sourcev1.GitRepositorySpec{
				Interval: metav1.Duration{Duration: defaultInterval},
				URL:      app.Spec.Source.RepoURL,
				Reference: &sourcev1.GitRepositoryRef{
					// TODO: This should do more to convert
					Branch: translateTargetRevision(app.Spec.Source.TargetRevision),
				},
			},
		},
		&kustomizev1.Kustomization{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "kustomize.toolkit.fluxcd.io/v1beta2",
				Kind:       "Kustomization",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.GetName(),
				Namespace: app.GetNamespace(),
			},
			Spec: kustomizev1.KustomizationSpec{
				Interval:        metav1.Duration{Duration: defaultInterval},
				TargetNamespace: app.GetNamespace(),
				SourceRef: kustomizev1.CrossNamespaceSourceReference{
					Kind: "GitRepository",
					Name: app.GetName(),
				},
				Path: app.Spec.Source.Path,
			},
		},
	}

	return result, nil
}

func translateTargetRevision(s string) string {
	if s == "HEAD" {
		return ""
	}
	return s
}

func loadSecretForRepo(ctx context.Context, cl client.Client) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	if err := cl.Get(ctx, name, app); err != nil {
		return nil, fmt.Errorf("failed to get Application %s: %w", name, err)
	}
}
