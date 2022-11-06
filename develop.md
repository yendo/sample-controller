# sample-controller development
Kubernetes sample controller

## Preparation

Install asdf

``` bash
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.10.2
```

Install Kubebuilder

``` bash
asdf plugin add kubebuilder
asdf install kubebuilder 3.7.0
asdf global kubebuilder 3.7.0
```

Install Kustomize

``` sh
asdf plugin add kustomize
asdf install kustomize 4.5.7
asdf global kustomize 4.5.7
```

## Init Kubebilder

Init Kubebuilder

``` bash
kubebuilder init --domain yendo.github.io --repo github.com/yendo/sample-controller --project-name foo
```

Create API

``` bash
kubebuilder create api --group samplecontroller --version v1 --kind Foo
make manifests
```

Create Webhook (optional)

``` bash
kubebuilder create webhook --group samplecontroller --version v1 --kind Foo --programmatic-validation --defaulting
make manifests
```

Edit config/manager/manager.yaml

``` diff
  image: controller:latest
+ imagePullPolicy: IfNotPresent
  name: manager
```

## Develop

Make

``` bash
make docker-build
kind load docker-image controller:latest
make install
make deploy
```

Check controller-manager pod

``` bash
kubectl get pod -n foo-system
```
