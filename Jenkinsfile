// Uses Declarative syntax to run commands inside a container.
pipeline {
    triggers {
        pollSCM("*/5 * * * *")
    }
    agent {
        kubernetes {
            yaml '''
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: dotnet
    image: mcr.microsoft.com/dotnet/sdk:5.0
    command:
    - sleep
    args:
    - infinity
'''
            defaultContainer 'dotnet'
        }
    }
    stages {
        stage('Lint code') {
            steps {
                sh "mkdir -p /usr/share/man/man1"
                sh "apt-get update"
                sh "apt-get install -y apt-utils"
                sh "apt-get install -y openjdk-11-jre-headless libzip-dev git wget unzip zip"
                sh 'java -version'
                sh 'dotnet tool install --global dotnet-sonarscanner --version 5.0.4'
                sh 'export PATH="$PATH:/root/.dotnet/tools"'
                sh 'dotnet sonarscanner begin /key:"jinya:backup" /name:"Jinya Backup" /d:sonar.host.url=https://sonarqube.imanuel.dev'
                sh 'dotnet build jinya-backup.sln'
                sh 'dotnet sonarscanner end'
            }
        }
    }
}
