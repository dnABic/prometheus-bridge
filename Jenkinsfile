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
      }
    }
    stage('Test') {
      agent { docker 'golang:1.9' }
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
