build_name = docker-deployx
build_path = build
build_target = $(build_path)/$(build_name)

source_target = cmd/deployx/*

install_path = $(HOME)/.docker/cli-plugins
install_target = $(install_path)/$(build_name)

.PHONY: build
build:
	go mod download
	go build -o $(build_target) $(source_target)

.PHONY: install
install: build
	mkdir -p $(HOME)/.docker/cli-plugins
	cp $(build_target) $(install_path)

.PHONY: uninstall
uninstall:
	rm $(install_target)
