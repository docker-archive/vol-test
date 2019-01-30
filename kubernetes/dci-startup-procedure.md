## Get started

Install dci locally, follow these instructions:

https://github.com/docker/certified-infrastructure

## Configure DCI for your environment:

* Add your Docker EE bundle IDs:
** `dci secret set docker/subscriptions_centos sub-super-sekrit-codez`
** Repeat for all subscriptions_sles
** `dci secret set docker/ucp_admin_username admin`
** `dci secret set docker/ucp_admin_password password`
** `dci secrets set aws/access_key yourkey`
** `dci secret set aws/secret_key`
* Start up your cluster:
** `dci cluster config set enable_kubernetes_aws_ebs true`
** `dci cluster apply`
** eval "$(dci cluster env)"
** `dci cluster env > kube.yml`
** `kubectl get pods` - confirm no pods running and we can talk to kube
** `kubectl apply -f aws-example-storage-class.yaml`
** `kubectl apply -f testapp.yaml`
** Check UCP for the exposed port for your test app. Kubernetes -> Load Balancers -> voltestservice. You are looking for "Nodeport"
** copy ~/.dci/cluster/aws/docker/ucp-bundle/kube.yml kube.yml (for convenience)
** ./voltestkube --config="kube.yml" --podurl="http://3.17.176.248:34511/" (adjust podurl for your environment)
