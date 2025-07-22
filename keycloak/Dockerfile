# keycloak/Dockerfile

FROM quay.io/keycloak/keycloak:24.0.1

# Use keycloak.conf file or pass env vars instead
ENV KEYCLOAK_ADMIN=admin
ENV KEYCLOAK_ADMIN_PASSWORD=admin

# Copy custom realm config if needed
COPY realm-export.json /opt/keycloak/data/import/

# Use the modern start command
ENTRYPOINT ["/opt/keycloak/bin/kc.sh", "start-dev", "--import-realm"]
