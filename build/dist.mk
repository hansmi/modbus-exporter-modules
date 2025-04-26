PACKAGE_NAME := modbus-exporter-modules
PACKAGE_VERSION := devel
PACKAGE_ARCHIVE := $(directory/dist)/$(PACKAGE_NAME)-$(PACKAGE_VERSION).tar.gz
PACKAGE_TAR_ARGS := \
	--sort=name --group=root:0 --owner=root:0 --mode=a=r,u+w,a+X \
	--transform='s|^/*|$(PACKAGE_NAME)-$(PACKAGE_VERSION)/|' \
	--show-transformed-names

# Args:
# 1: Path to file.
define add_dist_file
$(PACKAGE_ARCHIVE): $1
endef

$(PACKAGE_ARCHIVE): | $(directory/dist/target)
	now=$$(date +%s) && \
	$(TAR) --mtime="@$$now" $(PACKAGE_TAR_ARGS) -cvzf $@ -- $(sort $^) && \
	$(TAR) -tvzf $@

.PHONY: dist
dist: | $(PACKAGE_ARCHIVE)

.PHONY: clean-dist
clean-dist:
	rm -rf $(directory/dist)

clean: | clean-dist
