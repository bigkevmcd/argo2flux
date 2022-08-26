# argo2flux

This is a very simple convertor loads an ArgoCD Application from a resource, and spits out Flux resources that would deploy the same Application.

This does not automatically configure auth, and it's really only a 30 minute hack to see what's possible.

```shell
$ go build ./cmd/convert
$ ./convert  --app-name guestbook --app-namespace argocd
---
apiVersion: source.toolkit.fluxcd.io/v1beta2
kind: GitRepository
metadata:
  creationTimestamp: null
  name: guestbook
  namespace: argocd
spec:
  interval: 10m0s
  ref: {}
  url: https://github.com/argoproj/argocd-example-apps.git
status: {}
---
apiVersion: kustomize.toolkit.fluxcd.io/v1beta2
kind: Kustomization
metadata:
  creationTimestamp: null
  name: guestbook
  namespace: argocd
spec:
  interval: 10m0s
  path: guestbook
  prune: false
  sourceRef:
    kind: GitRepository
    name: guestbook
  targetNamespace: argocd
status: {}
```

## Further work

 - [ ] Convert to HelmReleases
 - [ ] Allow specifying destination namespace
 - [ ] Convert more of the source reference
 - [ ] Extract the credentials for the repository (if configured) to a Secret
