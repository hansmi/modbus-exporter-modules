define sma_sbsxx_50_config :=
firmware_version := 03.08.12.R
params_archive_name := PARAMETER-SBSExx-50_$$(subst .,-,$$(firmware_version)).zip
params_archive_sha256 := a447bb54f11c2563cfa062b3fe55f30850118f19365d5f35513efacebfe41711
endef

$(eval $(call define_ennexos_module,sma_sbsxx_50,$(sma_sbsxx_50_config)))
