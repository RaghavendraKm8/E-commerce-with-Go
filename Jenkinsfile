pipeline {
    agent any

    environment {
        REGISTRY = "docker.io"
        IMAGE_PREFIX = "Raghavendra Km"
    }

    stages {
        stage('Checkout') {
            steps {
                git url: 'https://github.com/RaghavendraKm8/E-commerce-with-Go.git', branch: 'master'
            }
        }

        stage('Build Docker Images') {
            steps {
                script {
                    def services = ["orders", "payments", "users"]
                    services.each { svc ->
                        sh """
                          echo "Building image for ${svc}..."
                          docker build -t ${IMAGE_PREFIX}/${svc}:latest ./services/${svc}
                        """
                    }
                }
            }
        }

        stage('Push Images to Docker Hub') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'docker-hub-password',
                                                 usernameVariable: 'DOCKER_HUB_USER',
                                                 passwordVariable: 'DOCKER_HUB_PASS')]) {
                    script {
                        sh '''
                          echo "$DOCKER_HUB_PASS" | docker login -u "$DOCKER_HUB_USER" --password-stdin
                        '''
                        def services = ["orders", "payments", "users"]
                        services.each { svc ->
                            sh """
                              echo "Pushing image for ${svc}..."
                              docker push ${IMAGE_PREFIX}/${svc}:latest
                            """
                        }
                    }
                }
            }
        }
    }

    post {
        always {
            echo 'Cleaning up Docker images...'
            sh 'docker image prune -af || true'
        }
        success {
            echo '✅ Build and push completed successfully!'
        }
        failure {
            echo '❌ Pipeline failed. Check logs above.'
        }
    }
}
