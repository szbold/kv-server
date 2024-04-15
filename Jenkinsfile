pipeline {
  agent any
  
  environment {
    DOCKER_CREDENTIALS = credentials('docker-hub')
    DOCKER_IMAGE_NAME = 'szbold/kv-server'
    DOCKER_TAG = 'latest'
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
        sh "docker build -t builder --target builder . > build.log 2>&1"
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
        sh "docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG} --target kv-server ."
      }
    }

    stage('Smoke test') {
      steps {
        echo 'Running smoke test'
        sh "chmod +x smoke_test.sh ${DOCKER_IMAGE_NAME}:${DOCKER_TAG}"
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
    
    stage('Publish') {
      steps {
        echo 'Publishing to docker hub...'
        sh "docker login -u ${DOCKER_CREDENTIALS_USR} -p ${DOCKER_CREDENTIALS_PSW}"
        sh "docker push ${DOCKER_IMAGE_NAME}:${DOCKER_TAG}"
      }
    }

    stage('Cleanup') {
      steps {
        echo 'Cleaning this mess...'
        sh "chmod +x ./cleanup.sh"
        sh "./cleanup.sh"
      }
    }
  }
}
