pipeline {
  agent any

    stages {
      stage('Build') {
        steps {
          echo 'Building..'
          sh 'ls -la'
          sh 'pwd'
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
