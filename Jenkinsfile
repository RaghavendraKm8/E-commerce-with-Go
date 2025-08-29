pipeline {
  agent any

  environment {
    DOCKER_HUB_REPO = "RaghavendraKm"  // note: no space
    IMAGE_TAG = "latest"
  }

  stages {
    stage('Checkout') {
      steps {
        checkout scm
      }
    }

    stage('Build Go Binaries') {
      parallel {
        ordersvc {
          dir('services/ordersvc') {
            bat 'go mod tidy'
            bat 'go build -o service.exe .'
          }
        }
        productsvc {
          dir('services/productsvc') {
            bat 'go mod tidy'
            bat 'go build -o service.exe .'
          }
        }
        usersvc {
          dir('services/usersvc') {
            bat 'go mod tidy'
            bat 'go build -o service.exe .'
          }
        }
      }
    }

    stage('Build Docker Images') {
      parallel {
        ordersvc {
          dir('services/ordersvc') {
            bat "docker build -t %DOCKER_HUB_REPO%/ordersvc:%IMAGE_TAG% ."
          }
        }
        productsvc {
          dir('services/productsvc') {
            bat "docker build -t %DOCKER_HUB_REPO%/productsvc:%IMAGE_TAG% ."
          }
        }
        usersvc {
          dir('services/usersvc') {
            bat "docker build -t %DOCKER_HUB_REPO%/usersvc:%IMAGE_TAG% ."
          }
        }
      }
    }

    stage('Push Images') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'docker-hub-password',
                                         usernameVariable: 'DOCKER_HUB_USER',
                                         passwordVariable: 'DOCKER_HUB_PASS')]) {
          bat 'echo %DOCKER_HUB_PASS% | docker login -u %DOCKER_HUB_USER% --password-stdin'
        }
        parallel {
          ordersvc {
            bat "docker push %DOCKER_HUB_REPO%/ordersvc:%IMAGE_TAG%"
          }
          productsvc {
            bat "docker push %DOCKER_HUB_REPO%/productsvc:%IMAGE_TAG%"
          }
          usersvc {
            bat "docker push %DOCKER_HUB_REPO%/usersvc:%IMAGE_TAG%"
          }
        }
      }
    }

    stage('Deploy with Docker Compose') {
      steps {
        bat 'docker-compose up -d'
      }
    }
  }

  post {
    always {
      echo 'Cleaning up workspace...'
      cleanWs()
    }
    failure {
      echo 'Pipeline failed! Check logs.'
    }
  }
}
