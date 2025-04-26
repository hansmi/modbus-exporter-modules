.PHONY: download-all
download-all:

.PHONY: clean-download
clean-download:

clean: | clean-download

# Args:
# 1: Local filename.
# 2: URL.
# 3: SHA256 checksum.
define define_download
download/$1/filename := $1
download/$1/url := $2
download/$1/checksum_sha256 := $3
download/$1 := $(directory/cache)/$$(download/$1/filename)
download/$1/unverified := $$(download/$1).unverified

.SECONDARY: $$(download/$1/unverified)

$$(download/$1/unverified): private export URL = $$(download/$1/url)
$$(download/$1/unverified): | $(directory/cache/target)
	$(CURL) --fail --output $$@ "$$$$URL"

$$(download/$1): $$(download/$1/unverified)
	echo '$$(download/$1/checksum_sha256) *$$<' | sha256sum --check && \
	mv $$< $$@

download-all: | $$(download/$1)

_download_clean/$$(download/$1):
	rm -f $$(download/$1) $$(download/$1/unverified)

clean-download: | _download_clean/$$(download/$1)
endef
