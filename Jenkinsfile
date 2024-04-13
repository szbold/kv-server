pipeline {
  agent any
  
  triggers {
    pollSCM('* * * * *')
  }

  stages {
    stage('Init') {
      steps {
        sh "export DOCKER_BUILDKIT=1"
      }
    }

    stage('Build executable') {
      steps {
        echo 'Building executable...'
        sh "docker build -t builder --targer builder . > build.log 2>&1"
      }
    }
      
    stage('Test') {
      steps {
        echo 'Testing...'
        sh "docker build -t tester --target tester . > test.log 2>&1"
      }
    }

    stage('Build deploy image') {
      steps {
        echo 'Building deploy container' 
        sh "docker build -t kv-server --target kv-server ."
      }
    }

    stage('Smoke test') {
      steps {
        echo 'Running smoke test'
        sh "chmod +x smoke_test.sh"
        sh "./smoke_test.sh"
      }
    }
    
    stage('Create artifacts') {
      steps {
        echo 'Creating artifacts...'
        sh "tar -czf artifact_${env.BUILD_TAG}.tar.gz ./build.log ./test.log"
        archiveArtifacts artifacts: "artifact_*.tar.gz", fingerprint: true
      }
    }

    stage('Cleanup') {
      steps {
        echo 'Cleaning this mess...'
        sh "chmox +x ./cleanup.sh"
        sh "./cleanup.sh"
      }
    }
  }
}
