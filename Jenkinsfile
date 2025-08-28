pipeline {
    agent any

    environment {
        REGISTRY = "your-dockerhub-username"   // 🔹 change to your DockerHub username
        APP_NAME = "go-ecommerce"
    }

    stages {
        stage('Checkout') {
            steps {
                // Checkout default branch (master)
                git branch: 'master', url: 'https://github.com/RaghavendraKm8/E-commerce-with-Go.git'
            }
        }

        stage('Build Docker Images') {
            steps {
                script {
                    def services = ["products-api", "orders-api", "payments-api"]

                    for (service in services) {
                        dir("services/${service}") {
                            sh "docker build -t ${REGISTRY}/${service}:latest ."
                        }
                    }
                }
            }
        }

        stage('Push Images') {
            steps {
                script {
                    withCredentials([usernamePassword(
                        credentialsId: 'docker-hub-credentials', // 🔹 Add in Jenkins Credentials
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )]) {
                        sh "echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin"

                        def services = ["products-api", "orders-api", "payments-api"]
                        for (service in services) {
                            sh "docker push ${REGISTRY}/${service}:latest"
                        }
                    }
                }
            }
        }

        stage('Deploy with Docker Compose') {
            steps {
                script {
                    sh 'docker-compose down || true'
                    sh 'docker-compose up -d'
                }
            }
        }
    }

    post {
        success {
            echo "✅ Deployment successful!"
        }
        failure {
            echo "❌ Pipeline failed. Check logs!"
        }
    }
}
