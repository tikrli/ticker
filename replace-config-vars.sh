#!/usr/bin/env bash

if [ ! -f ./config.yml ]; then
  echo -e "\nGenerating a config.yml file"
  cat > ./config.yml <<EOL
  listen: ${HOST_AND_PORT}
  log_level: ${LOG_LEVEL}
  initiator: ${ADMIN_EMAIL}
  database: ${DATABASE_URL}
  secret: ${SECRET}
  twitter:
    consumer_key: ${TWITTER_CONSUMER_KEY}
    consumer_secret: ${TWITTER_CONSUMER_SECRET}
EOL
fi