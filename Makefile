SPOT_INVENTORY := ./inventory.yml
SSH_KEY := ~/.ssh/id_pi_ed25519

deploy_%:
	spot -t $* -v -i ${SPOT_INVENTORY} -k ${SSH_KEY}