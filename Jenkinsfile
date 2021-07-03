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
  volumes:
    - name: docker-sock
      hostPath:
        path: /var/run/docker.sock
  containers:
  - name: golang
    image: quay.imanuel.dev/dockerhub/library---golang:latest
    command:
    - cat
    tty: true
  - name: docker
    image: quay.imanuel.dev/dockerhub/library---docker:stable
    command:
    - cat
    tty: true
    volumeMounts:
    - mountPath: /var/run/docker.sock
      name: docker-sock
'''
            defaultContainer 'docker'
        }
    }
    stages {
        stage('Build production') {
            parallel {
                stage('Build CGO worker') {
                    steps {
                        container('golang') {
                            dir('./worker') {
                                sh "wget https://www.musl-libc.org/releases/musl-1.2.2.tar.gz"
                                sh "tar -xvf musl-1.2.2.tar.gz"
                                dir ('./musl-1.2.2') {
                                    sh "./configure"
                                    sh "make"
                                    sh "make install"
                                }
                                sh "CC=/usr/local/musl/bin/musl-gcc go build --ldflags '-linkmode external -extldflags \"-static\"' -o jinya-backup-worker ."
                                archiveArtifacts artifacts: 'jinya-backup-worker', followSymlinks: false
                            }
                        }
                    }
                }
                stage('Push docker image') {
                    steps {
                        container('docker') {
                            dir('./server') {
                                sh "docker build -t quay.imanuel.dev/jinya/jinya-backup:v1.$BUILD_NUMBER -f ./Dockerfile ."
                                sh "docker tag quay.imanuel.dev/jinya/jinya-backup:v1.$BUILD_NUMBER quay.imanuel.dev/jinya/jinya-backup:latest"
                                sh "docker tag quay.imanuel.dev/jinya/jinya-backup:v1.$BUILD_NUMBER jinyacms/jinya-backup:v1.$BUILD_NUMBER"
                                sh "docker tag quay.imanuel.dev/jinya/jinya-backup:v1.$BUILD_NUMBER jinyacms/jinya-backup:latest"

                                withDockerRegistry(credentialsId: 'quay.imanuel.dev', url: 'https://quay.imanuel.dev') {
                                    sh "docker push quay.imanuel.dev/jinya/jinya-backup:v1.$BUILD_NUMBER"
                                    sh "docker push quay.imanuel.dev/jinya/jinya-backup:latest"
                                }
                                withDockerRegistry(credentialsId: 'hub.docker.com', url: '') {
                                    sh "docker push jinyacms/jinya-backup:v1.$BUILD_NUMBER"
                                    sh "docker push jinyacms/jinya-backup:latest"
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
