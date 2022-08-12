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
	@echo Building $(docker_image)...
	$(gcloud) builds submit --tag $(docker_image)
	@echo Pushing to Google Cloud Run...
	$(gcloud) run deploy --image $(docker_image) \
		--platform managed \
		--allow-unauthenticated \
		--region $(region) $(cloudrunname)
