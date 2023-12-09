BINARY_NAME=crossweaver-binary
CONFIG="./config.json"
LOG_LEVEL=info
METRICS=true
METRICS_PORT=8001
RESET=false
all: build
#  make all && ./crossweaver-binay start --reset --config ./config.json
build:
	cd ./crossweaver && go build -o $(BINARY_NAME) . && mv $(BINARY_NAME) ../

install:
	cd crossweaver && go install .

run:
	./$(BINARY_NAME) start --config $(CONFIG) --verbosity $(LOG_LEVEL) --metrics $(METRICS) --metricsPort $(METRICS_PORT) --reset $(RESET)

clean:
	cd crossweaver && go clean && cd ../.. && rm $(BINARY_NAME)
