# Ozone Operator Theory

So you want to know how Ozone works? This document is for you.

## The Reason

With OpenShift you have the Operator Framework which provides a lot of capabilities but also incurs a large computational requirement to run.  With something like MicroShift you have a much more lightweight footprint, mostly due to the lack of Operators.

Ozone is a way to bridge the gap between the two in a multi-cluster hub-and-spoke model.  It provides a way to run Operators on managed clusters in a lightweight fashion, in a way that is more efficient than running the full Operator Framework on each managed cluster.  

This means you can run Operators on managed clusters without the computational overhead of running the full Operator Framework on other non-OpenShift clusters such as MicroShift, K3s, xKS, etc.

## The Architecture

Ozone integrates with the Operator Framework and Red Hat Advanced Cluster Management on the Hub cluster.  It is also deployed as an Operator that then is able to read the `ManagedCluster` CR objects to enumerate the remote clusters available in RHACM.  With that, you can create the `RemoteSubscription` CR object to tell Ozone which Operator you want to run on which remote cluster.  The `RemoteSubscription` CR has the following fields:

- `.metadata.name` - The name of the `RemoteSubscription` CR
- `.spec.remoteCluster` - The struct that defines the remote cluster to target and how
- `.spec.remoteCluster.name` - The name of the `ManagedCluster` CR object that represents the remote cluster
- `.spec.operator` - The struct that defines the Operator to deploy - this Operator must be available on the Hub clusters.
- `.spec.operator.packageName` - The name of the local Operator to deploy on the remote cluster
- `.spec.operator.channel` - The channel of the local Operator to deploy on the remote cluster - if this is not specified, the default channel will be used
- `.spec.operator.startingCSV` - The name of the CSV to deploy on the remote cluster - if this is not specified, the latest CSV in the channel will be used
- `.spec.operator.installPlanApproval` - The approval strategy for the install plan - if this is not specified, the default approval strategy will be used
- `.spec.operator.source` - The name of the `CatalogSource` CR object that contains the local Operator
- `.spec.operator.sourceNamespace` - The namespace of the `CatalogSource` CR object that contains the local Operator

With this applied to the Hub cluster, the Ozone Operator Controller will then look for a local Operator by looking up the `PackageManifest` that corresponds to what is defined in the `RemoteSubscription` CR.  The `.status.remoteCluster` field is updated with the results of the local Operator lookup.

If the local Operator is found, then a connection test is performed - the default credentials that RHACM uses to manage the cluster are used, in the future this could be overriden with a specific Secret.

If the connection test is successful, then the `RemoteSubscription` CR is updated with the `status.remoteCluster` field to indicate the connection test was successful.

## Conditions

The `RemoteSubscription` CR has the following conditions:

