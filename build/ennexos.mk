ENNEXOS_PARAM_ARCHIVE_BASE_URL := https://files.sma.de/downloads
ENNEXOS_PARAM_EXTRACT := tools/extract-ennexos-parameters
ENNEXOS_MODULE_GENERATE := tools/generate-ennexos-module/generate-ennexos-module

$(eval $(call define_go_program,generate-ennexos-module))

# Args:
# 1: Local filename.
# 2: Upstream filename.
# 3: SHA256 checksum.
define define_ennexos_parameter_download
ennexos_parameter_download/$1/url := $(ENNEXOS_PARAM_ARCHIVE_BASE_URL)/$2

$(call define_download,$1,$$(ennexos_parameter_download/$1/url),$3)

ennexos_parameter_download/$1 := $$(download/$1)
endef

define ennexos_module_reset :=
undefine firmware_version
undefine params_archive_name
undefine params_archive_sha256
undefine generator_config
endef

# Args:
# 1: Module name.
# 2: Configuration variable assignments to evaluate.
define define_ennexos_module
$(eval $(ennexos_module_reset))

# Default configuration
$(eval
generator_config := config.yaml
)
$(eval $2)

ennexos_module/$1/dir := src/$1
ennexos_module/$1/parameters := $$(ennexos_module/$1/dir)/parameters.json
ennexos_module/$1/generator_config := $$(ennexos_module/$1/dir)/$(generator_config)
ennexos_module/$1 := $(directory/module)/$1.yaml

$(call define_ennexos_parameter_download,$(params_archive_name),$(params_archive_name),$(params_archive_sha256))

ennexos_module/$1/parameters_archive := $$(ennexos_parameter_download/$(params_archive_name))

$$(ennexos_module/$1/parameters): private export DESCRIPTION = Parameters from firmware version $(firmware_version).
$$(ennexos_module/$1/parameters): $$(ennexos_module/$1/parameters_archive) $(ENNEXOS_PARAM_EXTRACT)
	set -x && \
	$(ENNEXOS_PARAM_EXTRACT) --output=$$@ --description="$$$$DESCRIPTION" $$<

$$(ennexos_module/$1): $$(ennexos_module/$1/parameters) $$(ennexos_module/$1/generator_config) $(ENNEXOS_MODULE_GENERATE) | $(directory/module/target)
	$(ENNEXOS_MODULE_GENERATE) -module=$1 -config=$$(ennexos_module/$1/generator_config) -params=$$< > $$@

all: | $$(ennexos_module/$1)

.PHONY: _ennexos_module_clean/$1
_ennexos_module_clean/$1:
	rm -f $$(ennexos_module/$1/parameters)

clean: | _ennexos_module_clean/$1

$$(eval $$(call add_dist_file,$$(ennexos_module/$1)))

$(eval $(ennexos_module_reset))
endef
