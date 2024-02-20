# Deployment

> This guide is intended for setting up Media Vault in a non-production environment for testing and evaluation purposes. 

## Prerequisites:

* **Docker**: Install Docker from https://docs.docker.com/get-started/.
* **Docker Compose**: Install Docker Compose from https://docs.docker.com/compose/install/.

## Steps:

* **Clone the Repository**:
    
    ```bash
    git clone https://github.com/rishabhkailey/media-vault.git
    ```
* **Navigate to the Development Environment**:
    
    ```bash
    cd deployment/development/
    ```
* **Update Secrets (Optional)**:

    The `.env` file contains sensitive information (passwords, database details). If you want to customize these, open and update the file accordingly. Please replace the default passwords in `.env` with unique, strong passwords. They are publicly accessible and pose a security risk.

* **Start the Local Environment:**
    
    ```bash
    docker-compose up -d
    ```
    This command builds and starts all necessary Docker containers in detached mode (-d), allowing them to run in the background.

* **Wait for Services**:

    ```bash
    # wait for media-vault service to become healthy
    docker-compose ps
    ```

* ***Access the Media Vault Application**
    
    Access the Media Vault UI via your web browser at [http://localhost:8181/](http://localhost:8181/). Use the credentials from the `.env` file (`INITIAL_USER` and `INITIAL_USER_PASSWORD` variables) to log in.

* **Access the Keycloak Admin UI**:
    
    Keycloak handles user management and authentication. Access the admin UI via your web browser at [http://localhost:8181/accounts/admin/](http://localhost:8181/accounts/admin/). Use the credentials from the `.env` file (`ADMIN_USER` and `ADMIN_PASSWORD` variables) to log in.



