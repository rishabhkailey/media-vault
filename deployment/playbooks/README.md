## Configure keycloak clients

using vars file
```bash
ansible-playbook configure-keycloak-clients.yaml --extra-vars "@vars/media-vault-init.yaml" -vv
```

using inline variables
```bash
ansible-playbook configure-keycloak-clients.yaml \
  --extra-vars="oidc_media_vault_service_client_id=media-vault-service" \
  --extra-vars="oidc_media_vault_service_client_secret=media-vault-service" \
  --extra-vars="oidc_media_vault_spa_client_id=media-vault-spa" \
  --extra-vars="keycloak_server=http://localhost:8081" \
  --extra-vars="keycloak_admin=admin" \
  --extra-vars="keycloak_admin_password=admin" \
  --extra-vars="keycloak_application_realm=media-vault-7" \
  --extra-vars="keycloak_intial_user_name=test" \
  --extra-vars="keycloak_intial_user_password=password" \
  --extra-vars="{\"oidc_media_vault_spa_web_origins\": [\"http://localhost:5173/\",\"http://localhost:8090/\"]}" \
  --extra-vars="{\"oidc_media_vault_spa_redirect_uris\": [\"http://localhost:5173/*\",\"http://localhost:8090/*\"]}"
```
