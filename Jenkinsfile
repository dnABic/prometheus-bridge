pipeline {
  agent none

  stages {
    stage('Build') {
      agent any
      steps {
        echo 'Building..'
        sh 'ls -la'
        sh 'pwd'
        sh 'whoami'
        sh 'sleep 300'
      }
    }
    stage('Test') {
      agent { docker {
          image 'golang:1.9'
          args '-v  /home/jenkins/workspace/test02:/home/jenkins/workspace/test02'
        }
      }
      steps {
        echo 'Testing..'
        sh 'sleep 200'
        sh 'ls -la'
        sh 'pwd'
        sh 'whoami'
      }
    }
    stage('Deploy') {
      agent any
      steps {
        echo 'Deploying....'
        sh 'docker ps'
      }
    }
  }
}
