GO=go
GCLOUD=gcloud

project=go-yhsif-com
region=us-central1
cloudrunname=go-yhsif-com
image=go-yhsif-com

gcloud=$(GCLOUD) --project=$(project)
docker_image=$(region)-docker.pkg.dev/$(project)/$(image)/$(image)

.PHONY: deploy

deploy:
	go mod vendor
	@echo Building $(docker_image)...
	$(gcloud) builds submit --tag $(docker_image)
	@echo Pushing to Google Cloud Run...
	$(gcloud) run deploy --image $(docker_image) \
		--platform managed \
		--allow-unauthenticated \
		--region $(region) $(cloudrunname) \
		--execution-environment gen1 \
		--concurrency 1 \
		--cpu-throttling \
		--cpu 0.1 \
		--memory 128Mi
	rm -Rf vendor
