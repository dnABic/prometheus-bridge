pipeline {
  agent { docker 'golang:1.7' }

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
        sh 'cd $PROJECT_GO_PATH && go list ./... | grep -v vendor | xargs go test -v'
      }
    }

    stage('Build') {
      steps {
        sh 'cd $PROJECT_GO_PATH && GOOS=linux CGO_ENABLED=0 go build -a -ldflags \'-extldflags "-static"\''
      }
    }

    stage('Pre Release') {
      when {
        expression {
          GIT_TAG = sh(returnStdout: true, script: "git describe --tags HEAD --always").trim()
          return GIT_TAG =~ /-/
        }
      }
      environment {
        GITHUB_TOKEN = credentials('release-token')
      }
      steps {
        sh '$GOPATH/bin/ghr -u mobilityhouse -t "$GITHUB_TOKEN" -r "$PROJECT_NAME" -prerelease "' + GIT_TAG + '" "$PROJECT_GO_PATH/$PROJECT_NAME"'
      }
    }

    stage('Release') {
      when {
        expression {
          GIT_TAG = sh(returnStdout: true, script: "git describe --tags --always").trim()
          return GIT_TAG =~ /^v[^\-]+$/
        }
      }
      environment {
        GITHUB_TOKEN = credentials('release-token')
      }
      steps {
        sh '$GOPATH/bin/ghr -u mobilityhouse -t "$GITHUB_TOKEN" -r "$PROJECT_NAME" ' +  GIT_TAG + ' "$PROJECT_GO_PATH/$PROJECT_NAME"'
      }
    }
  }

  stage('Docker build') {
    environment {
      DOCKER_HUB = credentials('tmhitadmin')
      DOCKER_REPO = 'mobilityhouse'
    }
    steps {
      GIT_TAG = sh(returnStdout: true, script: "git describe --tags --always").trim()
      sh 'docker build -t $DOCKER_REPO/$PROJECT_NAME:GIT_TAG .'
      sh 'docker login -u $DOCKER_HUB_USR -p $DOCKER_HUB_PSW'
      sh 'docker push $DOCKER_REPO/$PROJECT_NAME:GIT_TAG'
    }
  }
}
