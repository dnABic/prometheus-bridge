pipeline {
  agent { docker 'golang:1.7' }

  environment {
    PROJECT_GO_PATH = "/usr/src/go/src/prometheus-bridge"
  }

  stages {
    stage('Prepare Environment') {
      steps {
        sh 'mkdir -p $(dirname $PROJECT_GO_PATH)'
        sh 'ln -s $(pwd) $PROJECT_GO_PATH'
      }
    }

    stage('Test') {
      steps {
        dir("${env.PROJECT_GO_PATH}") {
          sh 'go list ./... | grep -v vendor | xargs go test -v'
        }
      }
    }

    stage('Build') {
      steps {
        dir("${env.PROJECT_GO_PATH}") {
          sh 'GOOS=linux CGO_ENABLED=0 go build -a -ldflags \'-extflags "-static"\''
        }
      }
    }
  }
}
