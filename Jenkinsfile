pipeline {
    agent any

    environment {
        REGISTRY = "your-dockerhub-username"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'master',
                    url: 'https://github.com/RaghavendraKm8/E-commerce-with-Go.git'
            }
        }

        stage('Build Docker Images') {
            steps {
                script {
                    def services = ["products-api", "orders-api", "payments-api"]

                    for (service in services) {
                        dir("services/${service}") {
                            bat "docker build -t %REGISTRY%/${service}:latest ."
                        }
                    }
                }
            }
        }

        stage('Push Images') {
            steps {
                script {
                    def services = ["products-api", "orders-api", "payments-api"]

                    for (service in services) {
                        bat "docker push %REGISTRY%/${service}:latest"
                    }
                }
            }
        }

        stage('Deploy with Docker Compose') {
            steps {
                script {
                    // Stop previous containers (ignore errors if none running)
                    bat "docker-compose down || exit 0"

                    // Start new containers
                    bat "docker-compose up -d --build"
                }
            }
        }
    }

    post {
        always {
            echo "Cleaning up workspace..."
            cleanWs()
        }
    }
}
