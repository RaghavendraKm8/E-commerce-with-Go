pipeline {
    agent any

    environment {
        DOCKER_HUB_REPO = "Raghavendra Km
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Services') {
            steps {
                script {
                    bat 'docker build -t %DOCKER_HUB_REPO%/service1 ./service1'
                    bat 'docker build -t %DOCKER_HUB_REPO%/service2 ./service2'
                    bat 'docker build -t %DOCKER_HUB_REPO%/service3 ./service3'
                }
            }
        }

        stage('Docker Login') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'docker-hub-password',
                                                 usernameVariable: 'DOCKER_HUB_USER',
                                                 passwordVariable: 'DOCKER_HUB_PASS')]) {
                    script {
                        bat """
                          echo %DOCKER_HUB_PASS% | docker login -u %DOCKER_HUB_USER% --password-stdin
                        """
                    }
                }
            }
        }

        stage('Push Images') {
            steps {
                script {
                    bat 'docker push %DOCKER_HUB_REPO%/service1'
                    bat 'docker push %DOCKER_HUB_REPO%/service2'
                    bat 'docker push %DOCKER_HUB_REPO%/service3'
                }
            }
        }
    }
}
