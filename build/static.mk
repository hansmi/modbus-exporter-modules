# Args:
# 1: Module name.
define define_static_module
static_module/$1/dir := src/$1
static_module/$1/source := src/$1/module.yaml
static_module/$1 := $(directory/module)/$1.yaml

$$(static_module/$1): $$(static_module/$1/source) | $(directory/module/target)
	cp $$< $$@

all: | $$(static_module/$1)

$$(eval $$(call add_dist_file,$$(static_module/$1)))
endef
