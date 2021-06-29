# Jinya Backup

Jinya Backup is a push based backup system for shared hosting webspaces. It can backup directories and mysql databases.

## Server

The server side component is available as docker image and can be pulled from docker hub. For more details check the server README.md.

## Worker

The worker component is a simple utility to execute archive and mysql dump jobs. The jobs then trigger the server component, so that the server downloads the artifacts. For more details check the worker README.md.

## Found a bug?
If you found a bug feel free to create an issue on Github or on my personal Taiga instance: https://taiga.imanuel.dev/project/jinya-backup/

## License
Like all projects Jinya projects, Jinya Backup is distributed under the MIT License.
