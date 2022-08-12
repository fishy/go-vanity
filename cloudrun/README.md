# Google Cloud Run deploy

This directory deploys a go-vanity service to
[Google Cloud Run](https://cloud.google.com/run).

## Prerequisite

To deploy this, you need to first:

* [Install gcloud](https://cloud.google.com/sdk/docs/install)
* [Enable Google Cloud Run API](https://console.cloud.google.com/run)
  (you do not need to create a project there)
* [Enable and create an Artifactory Registry docker repository](https://console.cloud.google.com/artifacts)

## Config and `Makefile` overrides

Copy [`config.yaml.example`](config.yaml.example) into `config.yaml` and
configure it as you needed, then use `make deploy` with the following overrides:

* `--project`: Your Google Cloud project id
* `--region`: The region of your Artifactory Registry docker repository
* `--cloudrunname`: The cloud run service name you want to use (does not need to
  be created before hand)
* `--image`: The docker image name you want to use (can be anything)
