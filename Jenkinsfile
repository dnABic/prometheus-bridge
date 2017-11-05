pipeline {
  agent any

  stages {
    stage('Build') {
      steps {
        echo 'Building..'
        sh 'ls -la'
        sh 'pwd'
        sh 'whoami'
        sh 'curl -qo docker https://master.dockerproject.org/linux/x86_64/docker && chmod u+x docker'
        sh 'sleep 300'
        sh 'find / -name "docker"'
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
