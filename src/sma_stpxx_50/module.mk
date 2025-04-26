define sma_stpxx_50_config :=
firmware_version := 03.10.09.R
params_archive_name := PARAMETER-STPxx-50_$$(subst .,-,$$(firmware_version)).zip
params_archive_sha256 := 913c75f5d94539d52d5594839477ff29b45372e4df4634241454bf611dbc6666
endef

$(eval $(call define_ennexos_module,sma_stpxx_50,$(sma_stpxx_50_config)))
