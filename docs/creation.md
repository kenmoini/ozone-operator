
```bash=
operator-sdk init --domain operator.o3 --repo github.com/kenmoini/ozone-operator --project-name ozone-operator --owner "Ken Moini"

operator-sdk create api --group config --version v1alpha1 --kind RemoteSubscription --resource --controller

# Edited the API

# Always after editing types
make generate

# Make the manifests
make manifests
```