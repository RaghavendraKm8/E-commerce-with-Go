pipeline {
    agent any

    environment {
        REGISTRY = "your-dockerhub-username"    // change this to your DockerHub/registry username
        IMAGE_TAG = "latest"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/RaghavendraKm8/E-commerce-with-Go.git'
            }
        }

        stage('Build Docker Images') {
            steps {
                script {
                    // Build user service
                    sh 'docker build -t $REGISTRY/usersvc:$IMAGE_TAG ./services/usersvc'

                    // Build product service
                    sh 'docker build -t $REGISTRY/productsvc:$IMAGE_TAG ./services/productsvc'

                    // Build order service
                    sh 'docker build -t $REGISTRY/ordersvc:$IMAGE_TAG ./services/ordersvc'
                }
            }
        }

        stage('Push Images') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub-creds',
                                                 usernameVariable: 'DOCKER_USER',
                                                 passwordVariable: 'DOCKER_PASS')]) {
                    script {
                        sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'

                        sh 'docker push $REGISTRY/usersvc:$IMAGE_TAG'
                        sh 'docker push $REGISTRY/productsvc:$IMAGE_TAG'
                        sh 'docker push $REGISTRY/ordersvc:$IMAGE_TAG'
                    }
                }
            }
        }

        stage('Deploy with Docker Compose') {
            steps {
                script {
                    sh 'docker compose down || true'
                    sh 'docker compose up -d --build'
                }
            }
        }
    }
}
