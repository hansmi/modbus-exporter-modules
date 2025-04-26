.PHONY: _mkdir_%
_mkdir_%:
	test -e $* || mkdir $*

# Args:
# 1: Name.
# 2: Directory path.
define define_directory
directory/$1 := $(if $2,$2,$1)
directory/$1/target := _mkdir_$$(directory/$1)
endef

$(eval $(call define_directory,cache,))
$(eval $(call define_directory,module,))
$(eval $(call define_directory,dist,))
