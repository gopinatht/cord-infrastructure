.PHONY: all cord-infra-initializer cord-infra-sidecar

all_build_targets = cord-infra-sidecar cord-infra-initializer

cord_infra_initializer_dir = ./cord-infra-initializer/
cord_infra_sidecar_dir = ./cord-infra-sidecar/

all: $(all_build_targets)

cord-infra-initializer:
	$(MAKE) -C $(cord_infra_initializer_dir)

cord-infra-sidecar:
	$(MAKE) -C $(cord_infra_sidecar_dir)
