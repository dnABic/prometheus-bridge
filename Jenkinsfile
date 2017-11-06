pipeline {
  agent {
    docker {
      image 'golang:1.9'
      args '-u 0:0'
    }
  }

  environment {
    GOPATH = "/usr/src/go"
    PROJECT_NAME = "prometheus-bridge"
    PROJECT_GO_PATH = "/usr/src/go/src/prometheus-bridge"
  }

  stages {
    stage('Prepare Environment') {
      steps {
        sh 'mkdir -p $(dirname $PROJECT_GO_PATH)'
        sh 'ln -s $(pwd) $PROJECT_GO_PATH'
        sh 'go get github.com/tcnksm/ghr'
      }
    }
    stage('Test') {
      steps {
        sh 'cd $PROJECT_GO_PATH && go list ./... | grep -v vendor | xargs go test -v -cover'
      }
    }

    stage('Build') {
      steps {
        sh 'cd $PROJECT_GO_PATH && GOOS=linux CGO_ENABLED=0 go build -a -ldflags \'-extldflags "-static"\''
      }
    }

    stage('Docker build') {
      environment {
        DOCKER_REPO = 'dnabic'
      }
      steps {
        script {
          GIT_TAG = sh(returnStdout: true, script: "git describe --tags --always").trim()
        }

        sh 'curl -qo docker https://master.dockerproject.org/linux/amd64/docker && chmod u+x docker'
      }
    }

    stage('Build v2') {
      steps {
        echo 'Building..'
        sh 'ls -la'
        sh 'pwd'
      }
    }
    stage('Test v2') {
      steps {
        echo 'Testing..'
        sh 'ls -la'
        sh 'pwd'
        sh 'whoami'
      }
    }
    stage('Deploy') {
      steps {
        echo 'Deploying....'
      }
    }
  }
}
