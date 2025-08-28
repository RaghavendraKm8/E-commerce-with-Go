pipeline {
    agent any

    environment {
        DOCKER_HUB_USER = "your-dockerhub-username"
        DOCKER_HUB_PASS = credentials('docker-hub-password') // <-- ID must exist in Jenkins
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Go Services') {
            steps {
                script {
                    dir("services/ordersvc") {
                        bat "go mod tidy"
                        bat "go build -o service.exe ."
                    }
                    dir("services/productsvc") {
                        bat "go mod tidy"
                        bat "go build -o service.exe ."
                    }
                    dir("services/usersvc") {
                        bat "go mod tidy"
                        bat "go build -o service.exe ."
                    }
                }
            }
        }

        stage('Build Docker Images') {
            steps {
                script {
                    dir("services/ordersvc") {
                        bat "docker build -t ${DOCKER_HUB_USER}/ordersvc:latest ."
                    }
                    dir("services/productsvc") {
                        bat "docker build -t ${DOCKER_HUB_USER}/productsvc:latest ."
                    }
                    dir("services/usersvc") {
                        bat "docker build -t ${DOCKER_HUB_USER}/usersvc:latest ."
                    }
                }
            }
        }

        stage('Push Docker Images') {
            steps {
                script {
                    bat "docker login -u ${DOCKER_HUB_USER} -p ${DOCKER_HUB_PASS}"
                    bat "docker push ${DOCKER_HUB_USER}/ordersvc:latest"
                    bat "docker push ${DOCKER_HUB_USER}/productsvc:latest"
                    bat "docker push ${DOCKER_HUB_USER}/usersvc:latest"
                }
            }
        }

        stage('Deploy with Docker Compose') {
            steps {
                script {
                    bat "docker-compose down || echo 'No containers to stop'"
                    bat "docker-compose up -d"
                }
            }
        }
    }
}
