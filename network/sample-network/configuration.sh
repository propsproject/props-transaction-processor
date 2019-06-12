
#!/bin/bash

export CONFIG_LOG_FILE_PATH="/var/log/sawtooth/validator-debug.log"
export CONFIG_PATH="/etc/sawtooth"
export KEY_DIR="/etc/sawtooth/keys"
export SAWTOOTH_DATA="/var/lib/sawtooth"

sawadm genesis
