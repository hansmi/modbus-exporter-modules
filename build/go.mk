# Args:
# 1: Directory.
# 1: Go subcommand and arguments, e.g. "build".
define _run_go
$(GO) -C $1 $2
endef

.PHONY: clean-go
clean-go:

clean: | clean-go

.PHONY: _go_clean/%
_go_clean/%:
	$(call _run_go,$*,clean -x)

# Args:
# 1: Program name. Must exist as a directory in "tools".
define define_go_program
go/$1/dir := tools/$1

.PHONY: $$(go/$1/dir)/$1.tmp
$$(go/$1/dir)/$1.tmp:
	$(call _run_go,$$(@D),build -o $$(@F))

$$(go/$1/dir)/$1: $$(go/$1/dir)/$1.tmp
	if test -e $$@ && cmp --silent $$@ $$@.tmp; then \
		rm -f $$@.tmp; \
	else \
		mv $$@.tmp $$@; \
	fi

all: | $$(go/$1/dir)/$1

clean-go: | _go_clean/$$(go/$1/dir)
endef
