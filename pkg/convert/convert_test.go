package convert

import (
	"context"
	"os"
	"testing"

	argocdv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/yaml"
)

func TestConvertToKustomization(t *testing.T) {
	app := readApplication(t, "testdata/application.yaml")
	fc := newFakeClient(t, app)

	objs, err := ConvertToKustomization(context.TODO(), fc, client.ObjectKeyFromObject(app))
	if err != nil {
		t.Fatal(err)
	}
	want := []runtime.Object{
		&sourcev1.GitRepository{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "source.toolkit.fluxcd.io/v1beta2",
				Kind:       "GitRepository",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "guestbook",
				Namespace: "argocd",
			},
			Spec: sourcev1.GitRepositorySpec{
				Interval:  metav1.Duration{Duration: defaultInterval},
				URL:       "https://github.com/argoproj/argocd-example-apps.git",
				Reference: &sourcev1.GitRepositoryRef{},
			},
		},
		&kustomizev1.Kustomization{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "kustomize.toolkit.fluxcd.io/v1beta2",
				Kind:       "Kustomization",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "guestbook",
				Namespace: "argocd",
			},
			Spec: kustomizev1.KustomizationSpec{
				TargetNamespace: "argocd",
				Interval:        metav1.Duration{Duration: defaultInterval},
				SourceRef: kustomizev1.CrossNamespaceSourceReference{
					Kind: "GitRepository",
					Name: "guestbook",
				},
				Path: "guestbook",
			},
		},
	}
	if diff := cmp.Diff(want, objs); diff != "" {
		t.Fatalf("failed to convert:\n%s", diff)
	}
}

func readApplication(t *testing.T, filename string) *argocdv1.Application {
	b, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	app := &argocdv1.Application{}
	if err := yaml.Unmarshal(b, app); err != nil {
		t.Fatal(err)
	}

	return app
}

func newFakeClient(t *testing.T, objs ...runtime.Object) client.Client {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := argocdv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	return fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(objs...).
		Build()
}
