# Jinya Backup Server

The Jinya Backup Server is the server side component, installed on a server of your choice. This server must be able to
run Docker images.

## Installation

The installation is fairly simple, just get the Docker image from `jinyacms/jinya-backup:latest`. The version numbers
are simply the build numbers from Jenkins. If you don't want to use the latest tag check for the highest number in the
docker tags.

## Configuration

To configure Jinya Backup you need to set several environment variables. Please make sure, that you have set all of
them.

Variable                 | Description
-------------------------|------
`DB_HOST`                | The host of the database
`DB_PORT`                | The port of the database
`DB_USER`                | The username for the database
`DB_PASSWORD`            | The password for the database
`DB_DATABASE`            | The database to use
`DB_SECRET_KEY`          | The secret key used for encrypting FTP passwords
`DB_SECRET_NONCE`        | The secret nonce used for encrypting FTP passwords
`DB_FIRST_USER_NAME`     | The first user created will get this username
`DB_FIRST_USER_PASSWORD` | The first user created will get this password