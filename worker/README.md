# Jinya Backup Worker
The Jinya Backup Worker is the push part of Jinya Backup. It executes either `mysqldump` or archive jobs.

## Installation

Installing the Worker is pretty easy, you only need to download the latest version from here: https://jenkins.imanuel.dev/blue/organizations/jenkins/Jinya-CMS%2Fjinya-backup/activity Click the latest build number and download the `jinya-backup-worker` file from the artifacts tab.

Apart from the download you need to create a config file, check the Configuration section for the YAML structure. 

## Configuration
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