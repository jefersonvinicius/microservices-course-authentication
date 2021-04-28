#!/bin/sh
docker run --volume "$(pwd)/keycloak-data":/opt/jboss/keycloak/standalone/data/ -p 8080:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin quay.io/keycloak/keycloak:12.0.4