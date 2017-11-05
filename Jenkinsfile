pipeline {
  agent any

  stages {
    stage('Build') {
      steps {
        echo 'Building..'
        sh 'ls -la'
        sh 'pwd'
        sh 'whoami'
        sh 'find / -name "docker"'
        sh 'sleep 300'
      }
    }
    stage('Test') {
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
        sh 'docker ps'
      }
    }
  }
}
