pipeline {
  agent any

    stages {
      stage('Prepare Build') {
      }
      stage('Build') {
        steps {
          echo 'Building..'
          sh 'ls -la'
          sh 'pwd'
          docker.image('dnabic/prometheus-bridge:001').inside("-u 1000:1000") {
            sh 'whoami'
          }
        }
      }
      stage('Test') {
        steps {
          echo 'Testing..'
        }
      }
      stage('Deploy') {
        steps {
          echo 'Deploying....'
        }
      }
    }
}
