## flow for grant type = oauth2.AuthorizationCode

```plantuml
@startuml
!theme sketchy-outline
scale 3000*1080
actor User

note over Client
  1. oauth client configuration
  2. set client_id, client_secret
  3. Endpoint(auth and token urls)
  4. State (used to prevent CSFR)
endnote

note over Server
  1. srv.SetPasswordAuthorizationHandler(
       func(ctx, clientID, user, pass) (uID, err)
     )
  2. set user authorization handler to our function 
       srv.SetUserAuthorizationHandler(
         func(w http.ResponseWriter, r *http.Request) (uID string, err)
       )
endnote

User -> Client : / (protected endpoint)

note over Client
http.Redirect to config.AuthCodeURL
endnote

Client -> User : Redirect to AuthCodeURL (http://server/oauth/authorize)
User -> Server : Get (/oauth/authorize)

note over Server
checks if session store has user id.
if not it means the user has not login yet and redirect user to login page
endnote

Server -> User : redirect to Login Page (http://server/login)
User -> Server : Get (/login)

Server -> User : Login Page
User -> Server : Post uname and pass (/login)

note over Server
  1. set LoggedInUserID in session
      this will be used to tell if user is logged in or not
  2. set Location header to /auth
  3. set status code to 302
endnote

Server -> User : Redirect to /auth
User -> Server : Get /auth
Server -> User : Auth Page (similar to github ask for permission)
User -> Server : Clicks on Allow button. Post /oauth/authorize

note over Server
  1. set returUri if found in session
  2. forward the request to srv.HandleAuthorizeRequest(w, r)
  3. calls userAuthorizationHandler that we set
      get the userid from session.LoggedInUserID and returns use
  4. redirect to returnUri with token set in query param
     client/oauth2?code=<>&state=<>
endnote

Server -> User : Redirect
User -> Client : http://client/oauth2?code=<>&state=<>

note over Client
  config.Exchange(ctx, code)
endnote

Client -> Server : Get token using oauth2 client library (/oauth/token)

@enduml
```

![UML Diagram](https://www.plantuml.com/plantuml/svg/ZLN1RkGs3BtxAzXBCLvaFRD93qNGRdUHRji5zX0a27ei59jZ2z58BabbadvzJKbkZ3W9cai6flZu-FX4VM8n9F7Wgkzam0C2tw6qm-Fx4CLPZnMtnY6SxNQxTnzsF-ugqqeWk66agl914C8z4bmuYrugW0yD11DbW3PzWJRulUqZ6R71J_7J1XYbHF-otSdyax4bL8Ip1hxuRWpMooOX6T-1X3lq4CbnNG7yqC2r647OH4ODmaXuh-Kkhd-xgYlqdL9R4hn6clvbWaptpJNAfM4-1ke-JqK2sNyInJ-cQWvfKq7-sqVVRbfvc7b-_VK48YETm3YbrnC3_O94TSwedthKOs2Mu31aT6KS8YNmz7skLfYfkc-oAhoEC8YCpHNo63pZdsH5cH6yAz__ZiYIIG8BMTylkG8ifDAQyFvZcIISmnOs8qsXLh03B1EfNnXwgTLPmbP0GfbvenrSX0vlhhuzbSavMYWLF8TbuY83DWfxljro6jmscMexgudrp3SFLe5-Hu7D-jpQ0-s0xHs3tGCZiyhF4WXXC9ndPhkccg9J8bY10nhF80Fce9xIY0kzzV28aln9fOjyHa89ZwR7fzvp_KNljEZzMqgud1AUjvwWNcquHzVW2y1rwcLWWUZD0PLzil82wfLRWwAOFNPVlG9Ebh9-bh0OIWRBSBJEmIt2V3S5dLErixOCBe5eSY2LCj-O3D-MUuAcopgcKUQje2TOZ4I6DdIeqRFTwLhUjRLcd9SrpA5Ljfgnp8FjmJgJ-FHMXdWBXkzW7mX6f8DD2gm7fFT3ZHOyV7Ok7E4sYWJVPFMVE_KjsGabqWrPLN4Vek-MogjwokLeg4isfNpjbQtkaxmxvdM2PIbiZYT0THQsDSunn5UM3iXW18weJCgOUvJ_BeJjO4_XCDDfdbbamiZqVQgGb_VIzpcepHsj36NHCugsE36bHnWDcKEkd7Tmbkxq5tN1ppzz_5vDeJ_-XnNMQwxSjRUHNtpgofhxyj0EnlUONmdDh9-lkvNbodF6rlUGIvRMmDbRCbFJyoRJatedFqsG-ah_2m00)


## Migrations

### Installation
```bash
platform=linux
version=v4.15.2
curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$platform-amd64.tar.gz | tar xvz migrate
sudo mv migrate /usr/bin
```

### Creating new migration
```bash
migrate create -seq -ext sql -dir migrations/ oauth_client_scopes_table
```

### Run migrations
```bash
# all migrations
migrate -source file://migrations -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" up

# 1 migration
migrate -source file://migrations -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" up 1

# migrate to specific version
migrate -source file://migrations -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" goto <version>
```

### Revert migrations
```bash
migrate -source file://migrations -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" down 1
```

### Check current version
```bash
migrate -source file://migrations -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" version
```


### Dirty DB fix
> error: Dirty database version 20221121071907. Fix and force version.
we can force change the version to previous version and fix the migration and apply again 
```bash
migrate -source file://migrations -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" force <previous-version>
migrate -source file://migrations -database "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable" up 1
```
if the first migration is failing only option is to delete the entry in schema_verision table and retry.

## DB schema

https://dbdiagram.io/d/637b89c7c9abfc6111742f21
```dbdiagram.io
// 
Table oauth2_clients {
  id TEXT [pk]
  secret TEXT
  domain TEXT
  data JSONB
}

// k8s verbs i.e. get, create, edit, update per client 
// or k8s roles admin, view, developer
// role = scope
// verbs = claims 
Table oauth2_client_scopes {
  id int [pk]
  client_id text
  scope text
  indexes {
    (client_id) [name: 'oauth2_client_scopes_index_client_id', note: 'this will be used to get all the available scopes for a client']
    (client_id, scope) [name:'oauth2_client_scopes_client_id_scope', unique, note: 'this will be used get scope_id']
  }
}

Ref: oauth2_clients.id < oauth2_client_scopes.client_id

// k8s users per client
Table users {
  id int [pk]
  email text
  username text
  password text
  client_id TEXT
  indexes {
    (email, client_id) [unique, name: 'users_index_email_client_id', unique, note: 'this will be used for getting username, id and password of a user']
  }
}

Ref: oauth2_clients.id < users.client_id



// k8s rolesbindings
Table user_scope_bindings {
  id int [pk]
  user_id int
  // is this redundant as from user_id we can get the client_id?
  // client_id int
  scope_id int
  indexes {
    (user_id) [name: 'user_scope_bindings_index_user_id', note: 'this will be used to get all the scopes assigned to the user']
    (user_id, scope_id) [name: 'user_scope_bindings_index_user_id_scope_id', note: 'this will be used to verify if the user has this scope assigned']
  }
}

Ref: users.id < user_scope_bindings.user_id
Ref: oauth2_client_scopes.id < user_scope_bindings.scope_id
```


### Example client schema

```dbdiagram.io
// this should be handled by the resource server
// k8s resources pod, deployment
Table resource {
  id int [pk]
  resource_name text
  resource_type text
  owner user_id
  indexes {
    (id)
    (owner)
    (owner, resource_type)
    (owner, resource_name)
  }
}

Table resource_shared {
  id int [pk]
  resource_id int
  shared_with int
  indexes {
    (resource_id)
    (shared_with)
    (resource_id, shared_with)
  }
}
```

## Generate keys for JWT token

```bash
openssl genrsa -out private-key.pem 2048

# we actualy don't need to create a separate file for public key as public key is just a part of private, in code we can just get public key from private key.
openssl rsa -in private-key.pem -pubout -out public-key.pem
# add output of following in config jwt.privateKey
cat private-key.pem | base64 -w 0
# add output of following in config jwt.publicKey
cat public-key.pem | base64 -w 0
```

## Client CMD

### Create Client 

```bash
# todo update doc that profile and email scopes are default and will be ignored in this command
go run cmd/auth-server/main.go client create --scopes user --scopes admin --defaultScopes user --domain localhost --id test --secret aslkjfdalksdfjlksadfjlkasjlkfjasdfjlkasjd --returnUri localhost
```

### create userScopeBinding
```bash
go run cmd/auth-server/main.go userScopeBinding create --client test --email test2 --scopes user --scopes admin
```

### Useful CURLs

#### Singup

```bash
curl -d '{
  "email":"value1", 
  "username":"value2",
  "password":"value2"
}' -H "Content-Type: application/json" -X POST http://localhost:8080/signup
```

#### User Info

```bash
curl -H "Authorization: Bearer ${access_token}" http://localhost:8080/oauth/userinfo
```

#### Revoke token

```bash
curl -X POST -d "{
  \"token\": \"${access_token}\",
  \"token_type_hint\": \"access\"
}" http://localhost:8080/oauth/revoke
```

### Token Interospection

```bash
curl -H "Authorization: Bearer ${access_token}" http://localhost:8080/oauth/introspect
```