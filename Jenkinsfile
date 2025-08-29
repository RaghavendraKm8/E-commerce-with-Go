pipeline {
  agent any

  environment {
    DOCKER_HUB_REPO = "RaghavendraKm"
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
        stage('Build ordersvc') {
          steps {
            dir('services/ordersvc') {
              bat 'go mod tidy'
              bat 'go build -o service.exe .'
            }
          }
        }
        stage('Build productsvc') {
          steps {
            dir('services/productsvc') {
              bat 'go mod tidy'
              bat 'go build -o service.exe .'
            }
          }
        }
        stage('Build usersvc') {
          steps {
            dir('services/usersvc') {
              bat 'go mod tidy'
              bat 'go build -o service.exe .'
            }
          }
        }
      }
    }

    stage('Build Docker Images') {
      parallel {
        stage('Docker build ordersvc') {
          steps {
            dir('services/ordersvc') {
              bat "docker build -t %DOCKER_HUB_REPO%/ordersvc:%IMAGE_TAG% ."
            }
          }
        }
        stage('Docker build productsvc') {
          steps {
            dir('services/productsvc') {
              bat "docker build -t %DOCKER_HUB_REPO%/productsvc:%IMAGE_TAG% ."
            }
          }
        }
        stage('Docker build usersvc') {
          steps {
            dir('services/usersvc') {
              bat "docker build -t %DOCKER_HUB_REPO%/usersvc:%IMAGE_TAG% ."
            }
          }
        }
      }
    }

    stage('Docker Login') {
      steps {
        withCredentials([usernamePassword(credentialsId: '7022835052'
                                         usernameVariable: 'DOCKER_HUB_USER',
                                         passwordVariable: 'DOCKER_HUB_PASS')]) {
          bat 'echo %DOCKER_HUB_PASS% | docker login -u %DOCKER_HUB_USER% --password-stdin'
        }
      }
    }

    stage('Push Images') {
      parallel {
        stage('Push ordersvc') {
          steps {
            bat "docker push %DOCKER_HUB_REPO%/ordersvc:%IMAGE_TAG%"
          }
        }
        stage('Push productsvc') {
          steps {
            bat "docker push %DOCKER_HUB_REPO%/productsvc:%IMAGE_TAG%"
          }
        }
        stage('Push usersvc') {
          steps {
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
