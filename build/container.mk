CONTAINER_TAGS := $(PACKAGE_NAME):devel
CONTAINER_BUILD_ARGS := PACKAGE_ARCHIVE=$(PACKAGE_ARCHIVE)

.PHONY: container
container: build/Dockerfile $(PACKAGE_ARCHIVE)
	$(DOCKER) build --pull=false --progress=plain --file=$< \
		$(patsubst %,--tag='%',$(sort $(CONTAINER_TAGS))) \
		$(patsubst %,--build-arg='%',$(sort $(CONTAINER_BUILD_ARGS))) \
		.
