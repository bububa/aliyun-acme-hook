PROJECT_NAME=aliyun-acme-hook
ENTRY_POINT = ./cmd/app
DIST_PATH=./dist
EXEC_PATH=/usr/local/bin
LDFLAGS = -s -w -extldflags "-static"

.PHONY : app

app:
ifeq (,$(wildcard $(DIST_PATH)/$(PROJECT_PATH)))
	rm -rf $(DIST_PATH)/$(PROJECT_NAME)
endif
	go build -o $(DIST_PATH)/$(PROJECT_NAME) -ldflags "$(LDFLAGS)" $(ENTRY_POINT)

clean:
	rm -rf $(DIST_PATH)/*

install:
	@echo "Installing $(PROJECT_NAME) to $(EXEC_PATH)/$(PROJECT_NAME)"
	@if [ -w $(EXEC_PATH) ]; then \
		mv $(DIST_PATH)/$(PROJECT_NAME) $(EXEC_PATH)/$(PROJECT_NAME); \
	else \
		echo "Need to use sudo to install in $(EXEC_PATH). Please run:"; \
		echo "sudo mv $(DIST_PATH)/$(PROJECT_NAME) $(EXEC_PATH)/$(PROJECT_NAME)"; \
	fi
