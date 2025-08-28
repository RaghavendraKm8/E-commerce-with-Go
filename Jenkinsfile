pipeline {
    agent any

    tools {
        go 'Go'    // configure Go in Jenkins tools
    }

    environment {
        DOCKER_CLI_EXPERIMENTAL = "enabled"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'master', url: 'https://github.com/RaghavendraKm8/E-commerce-with-Go.git'
            }
        }

        stage('Build images') {
            steps {
                bat 'echo Building Docker images...'
                bat 'docker compose -f docker-compose.yml build --no-cache'
            }
        }

        stage('Unit tests') {
            steps {
                bat 'go test ./...'
            }
        }

        stage('Start stack') {
            steps {
                bat 'docker compose -f docker-compose.yml up -d'
            }
        }

        stage('Wait for health') {
            steps {
                script {
                    sleep(time:30, unit:"SECONDS")
                }
            }
        }

        stage('Integration smoke test') {
            steps {
                bat 'curl -f http://localhost:8081 || exit 1'
                bat 'curl -f http://localhost:8082 || exit 1'
                bat 'curl -f http://localhost:8083 || exit 1'
            }
        }
    }

    post {
        always {
            bat 'docker compose -f docker-compose.yml down -v || echo Nothing to clean'
            echo "✅ Cleanup finished"
        }
        failure {
            echo "❌ Pipeline failed"
        }
    }
}
