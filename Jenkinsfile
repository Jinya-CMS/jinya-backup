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
    image: golang:latest
    command:
    - cat
    tty: true
  - name: docker
    image: docker:latest
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
                                sh "cd musl-1.2.2"
                                sh "./configure"
                                sh "make"
                                sh "sudo make install"
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
                                sh "docker build -t registry-hosted.imanuel.dev/jinya/jinya-backup:shBUILD_NUMBER -f ./Dockerfile ."
                                sh "docker tag registry-hosted.imanuel.dev/jinya/jinya-backup:$BUILD_NUMBER jinyacms/jinya-backup:$BUILD_NUMBER"
                                sh "docker tag registry-hosted.imanuel.dev/jinya/jinya-backup:$BUILD_NUMBER jinyacms/jinya-backup:latest"

                                withDockerRegistry(credentialsId: 'nexus.imanuel.dev', url: 'https://registry-hosted.imanuel.dev') {
                                    sh "docker push registry-hosted.imanuel.dev/jinya/jinya-backup:$BUILD_NUMBER"
                                }
                                withDockerRegistry(credentialsId: 'hub.docker.com', url: '') {
                                    sh "docker push jinyacms/jinya-backup:$BUILD_NUMBER"
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
