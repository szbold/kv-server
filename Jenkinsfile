pipeline {
  agent any
  
  triggers {
    pollSCM('* * * * *')
  }

  stages {
    stage('Init') {
      steps {
        echo "${env.BUILD_TAG}"
      }
    }

    stage('Build') {
      steps {
        echo 'Building...'
        sh "docker build -t builder:1.0 -f ./Dockerfile.build . > build.log 2>&1"
      }
    }
      
    stage('Test') {
      steps {
        echo 'Testing...'
        sh "docker build -t tester:1.0 -f ./Dockerfile.test . > test.log 2>&1"
      }
    }
    
    stage('Upload artifacts') {
      steps {
        echo 'Uploading artifacts...'
        sh "tar -czf artifact_${env.BUILD_TAG}.tar.gz ./build.log ./test.log"
        archiveArtifacts artifacts: "artifact_*.tar.gz", fingerprint: true
      }
    }
  }
}
