build:
	@echo "Building Go Lambda function"
	@gox -os="linux" -arch="amd64" -output="araya"  

deploy-staging:build
	rsync -a araya ametory@103.172.205.9:/home/ametory/araya -v --stats --progress
	rsync -a template ametory@103.172.205.9:/home/ametory/araya -v --stats --progress

deploy-prod:build
	rsync -a araya ametory@146.190.86.62:/home/ametory/araya -v --stats --progress
	rsync -a template ametory@146.190.86.62:/home/ametory/araya -v --stats --progress



