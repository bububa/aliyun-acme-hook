PROJECT_NAME=aliyun-acme-hook
ENTRY_POINT = ./cmd/app
DIST_PATH=./dist
LDFLAGS = -s -w -extldflags "-static"

.PHONY : all

all: app

app:
ifeq (,$(wildcard $(DIST_PATH)/$(PROJECT_PATH)))
	rm -rf $(DIST_PATH)/$(PROJECT_NAME)
endif
	go build -o $(DIST_PATH)/$(PROJECT_NAME) -ldflags "$(LDFLAGS)" $(ENTRY_POINT)

clean:
	rm -rf $(DIST_PATH)/*

