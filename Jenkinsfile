pipeline {
  agent { docker 'golang:1.7' }

  environment {
    GOPATH = "/usr/src/go"
    PROJECT_GO_PATH = "/usr/src/go/src/prometheus-bridge"
  }

  stages {
    stage('Prepare Environment') {
      steps {
        sh 'mkdir -p $(dirname $PROJECT_GO_PATH)'
        sh 'ln -s $(pwd) $PROJECT_GO_PATH'
        sh 'go get https://github.com/tcnksm/ghr'
      }
    }

    stage('Test') {
      steps {
        sh 'cd $PROJECT_GO_PATH && go list ./... | grep -v vendor | xargs go test -v'
      }
    }

    stage('Build') {
      steps {
        sh 'cd $PROJECT_GO_PATH && GOOS=linux CGO_ENABLED=0 go build -a -ldflags \'-extldflags "-static"\''
      }
    }

    stage('Github Reelease') {
      environment {
        GITHUB_TOKEN = credentials('tmhjenkins')
      }
      steps {
        sh '$GOPATH/bin/ghr $BUILD_TAG $PROJECT_GO_PATH/prometheus-bridge'
      }
    }
  }
}
