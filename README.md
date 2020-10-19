# CRD-client-go

crd-client-go defines and implements the custom types used by FogAtlas in terms of
k8s Custom Resource Definition (CRD). Once defined the types, it generates the corresponding
API using k8s code generator (see below).

## FogAtlas CRDs definition and API

Types in FogAtlas are defined starting from the following domain models: the first one models
a distributed infrastructure while the second models an application as a graph of vertices (microservices)
and edges (dataflows).

![infrastructure model](./docs/images/fogatlas-datamodel-infra.png)
*infrastructure model*

![application model](./docs/images/fogatlas-datamodel-app.png)
*application model*

From these models the CRD definition are derived.

In the folder _crd_definitions_ you can find (in yaml format) the CRDs that extend the k8s resources:
* _crd_dynamicnode.yaml_: used only by the provisioning phase. Currently not used.
* _crd_externalendpoint.yaml_: defines an external endpoint (sensor, camera or external service)
* _crd_fadepl.yaml_: defines a so called FogAtlas deployment (FADepl) that models a cloud native application.
* _crd_fedfadepl.yaml_: extension of FADepl for multi-cluster. Not used at the moment
* _crd_region.yaml_: defines a region interconnected by a Link
* _crd_link.yaml_: defines a network link between two Regions

The file [_types.go_](pkg/apis/fogatlas/v1alpha1/types.go) defines programmatically the aforementioned CRDs.

## How to define or change CRDs
1. Define three files in pkg/apis/<api-group>/v1alpha1
   * _doc.go_ where global generation tags are defined
   * _types.go_ where custom types are defined
   * _register.go_ where custom types are registered to the k8s API
2. Use _update-codegen.sh_ script to generate the code. This step needs:
   * go get k8s.io/code-generator
   * go get k8s.io/apimachinery
3. Align the files inside _crd_definitions_ to the _types.go_  

## How to install CRDs

In order to install the defined CRD on a k8s cluster, just do the following:
```sh
cd crd-definitions
kubectl apply -f crd_region.yaml
kubectl apply -f crd_link.yaml
kubectl apply -f crd_externalendpoint.yaml
kubectl apply -f crd_fadepl.yaml
```

## How to test

The file _main.go_ provides an example on how to call both k8s vanilla API
(_get nodes()_ belonging to a k8s cluster) and Fogatlas API (_get regions()_ configured on a k8s cluster).
Of course in order to use it, you need a k8s cluster where the aforementioned CRDs are loaded
(see above) and where some instances of Region are defined. Moreover if you don't access the cluster
as k8s admin, you might need additional RBAC setup.

The syntax to launch the _main.go_ is as follows:

```sh
go run main.go --kubeconfig=<kube config path> --loglevel=<log level>
```
where _kube config path_ is the path where the k8s configuration file to access the cluster is stored
and _log level_ is the level of the log (use "trace").

## Schema validation.

With CRD v1beta1, we use the following section:

```sh
validation:
  openAPIV3Schema:
```
For activating the validation of the different fields, we need to create a so called _structural schema_.
If you want to impose some constraints on the schema, you need to use _allOf_, _oneOf_ etc. keywords.
An example can be found in the _crd_fadepl.yaml_ file.  

## License

Copyright 2019 FBK CREATE-NET

Licensed under the Apache License, Version 2.0 (the “License”); you may not use this
file except in compliance with the License. You may obtain a copy of the License
[here](http://www.apache.org/licenses/LICENSE-2.0).

Unless required by applicable law or agreed to in writing, software distributed under
the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
either express or implied. See the License for the specific language governing permissions
and limitations under the License.

