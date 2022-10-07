pusk:
	go run cmd/main.go
pusk-docker:
	docker image build -f Dockerfile . -t gr-img
	docker rmi $$(docker images -f "dangling=true" -q)
	docker run -p 8181:8181 --rm gr-img:latest

docker-build:
	docker image build -f Dockerfile . -t gr-img
	docker rmi $$(docker images -f "dangling=true" -q)

docker-run:
	docker run -p 8181:8181 --rm gr-img:latest