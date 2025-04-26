include build/defaults.mk
include build/dirs.mk
include build/download.mk
include build/go.mk
include build/ennexos.mk
include build/dist.mk
include build/container.mk

all:

.PHONY: clean
clean:

DIST_EXTRA := \
	LICENSE \
	README.md

$(foreach i,$(DIST_EXTRA),$(eval $(call add_dist_file,$(i))))

define MODULES
sma_stpxx_50
sma_sbsxx_50
endef

$(foreach i,$(sort $(MODULES)),$(eval include src/$(i)/module.mk))
