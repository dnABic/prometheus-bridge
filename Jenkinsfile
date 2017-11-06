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

    stage('Pre Release') {
      environment {
        GITHUB_TOKEN = credentials('release-token')
      }
      steps {
        sh 'echo $GOPATH/bin/ghr -u dnabic -t "$GITHUB_TOKEN" -r "$PROJECT_NAME" -prerelease --replace "latest" "$PROJECT_GO_PATH/$PROJECT_NAME"'
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
        sh 'echo $GOPATH/bin/ghr -u mobilityhouse -t "$GITHUB_TOKEN" -r "$PROJECT_NAME" ' +  GIT_TAG + ' "$PROJECT_GO_PATH/$PROJECT_NAME"'
      }
    }

    stage('Docker build') {
      environment {
        DOCKER_HUB = credentials('hub.dnabic')
        DOCKER_REPO = 'dnabic'
      }
      steps {
        script {
          GIT_TAG = sh(returnStdout: true, script: "git describe --tags --always").trim()
        }

        sh 'curl -qo /usr/bin/docker https://master.dockerproject.org/linux/x86_64/docker && chmod u+x /usr/bin/docker'
        sh 'echo ./docker build -t "$DOCKER_REPO/$PROJECT_NAME:' + GIT_TAG + '" .'
        sh 'echo ./docker tag "$DOCKER_REPO/$PROJECT_NAME:' + GIT_TAG + '" $DOCKER_REPO/$PROJECT_NAME:edge'
        sh 'echo ./docker login -u $DOCKER_HUB_USR -p $DOCKER_HUB_PSW'
        sh 'echo ./docker push "$DOCKER_REPO/$PROJECT_NAME:' + GIT_TAG + '"'
      }
    }

    stage("Integration Test") {
      steps {
        sh 'curl -L  -qo /usr/bin/docker-compose https://github.com/docker/compose/releases/download/1.12.0/docker-compose-`uname -s`-`uname -m`'
          sh 'chmod +x /usr/bin/docker-compose'
      }
    }
  }

  post {
    always {
      sh 'echo ./docker-compose -f integration/docker-compose.yaml down --remove-orphans || true'
    }
  }
}
