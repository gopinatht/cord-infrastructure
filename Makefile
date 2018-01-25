all_build_targets = cord-infra-initializer

cord_infra_dir = ./cord-infra-initializer

all: $(all_build_targets)

cord-infra-initializer:
	$(MAKE) -C $(cord_infra_dir)
