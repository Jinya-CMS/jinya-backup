# Jinya Backup

Jinya Backup is a push based backup system for shared hosting webspaces. It can backup directories and mysql databases.

## Jinya Backup Server

The Jinya Backup Server is the server side component, installed on a server of your choice. This server must be able to
run Docker images.

### Installation

The installation is fairly simple, just get the Docker image from `jinyacms/jinya-backup:latest`. The version numbers
are simply the build numbers from Jenkins. If you don't want to use the latest tag check for the highest number in the
docker tags.

### Configuration

To configure Jinya Backup you need to set several environment variables. Please make sure, that you have set all of
them.

| Variable                 | Description                                        |
|--------------------------|----------------------------------------------------|
| `DB_HOST`                | The host of the database                           |
| `DB_PORT`                | The port of the database                           |
| `DB_USER`                | The username for the database                      |
| `DB_PASSWORD`            | The password for the database                      |
| `DB_DATABASE`            | The database to use                                |
| `DB_SECRET_KEY`          | The secret key used for encrypting FTP passwords   |
| `DB_SECRET_NONCE`        | The secret nonce used for encrypting FTP passwords |
| `DB_FIRST_USER_NAME`     | The first user created will get this username      |
| `DB_FIRST_USER_PASSWORD` | The first user created will get this password      |

## Jinya Backup Worker

The Jinya Backup Worker is the push part of Jinya Backup. It executes either `mysqldump` or archive jobs.

### Installation

Installing the Worker is pretty easy, you only need to download the latest version from
here: https://jenkins.imanuel.dev/blue/organizations/jenkins/Jinya-CMS%2Fjinya-backup/activity Click the latest build
number and download the `jinya-backup-worker` file from the artifacts tab.

Apart from the download you need to create a config file, check the Configuration section for the YAML structure.

### Configuration

```yaml
server: Your Jinya backup server
jobs:
  - id: The ID of the Jinya Backup job, you can find it in the Jinya Backup UI
    output: The path of the created dump, the same filename must be specified in Jinya Backup
    host: The mysql server host
    port: The mysql server port
    user: The mysql server username
    pass: The mysql server password
    database: The mysql server database
    type: mysqldump

  - id: The ID of the Jinya Backup job, you can find it in the Jinya Backup UI
    input: The input directory
    output: The output file, should end with .tar.gz
    type: archive
```

## Found a bug?

If you found a bug feel free to create an issue on Github or on my personal Taiga
instance: https://taiga.imanuel.dev/project/jinya-backup/

## License

Like all projects Jinya projects, Jinya Backup is distributed under the MIT License.
