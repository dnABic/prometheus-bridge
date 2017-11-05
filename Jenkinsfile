pipeline {
  agent any

    stages {
      stage('Prepare Build') {
      }
      stage('Build') {
        steps {
          echo 'Building..'
          def build_id = "dnabic/prometheus-bridge:${env.BUILD_NUMBER}"
          sh 'ls -la'
          sh 'pwd'
          docker.image(build_id).inside("-u 1000:1000") {
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
