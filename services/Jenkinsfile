pipeline {
    agent any

    environment {
        COMPOSE = "docker compose -f docker-compose.yml"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build images') {
            steps {
                ansiColor('xterm') {
                    sh "${COMPOSE} build"
                }
            }
        }

        stage('Unit tests') {
            steps {
                ansiColor('xterm') {
                    sh 'go test ./...'
                }
            }
        }

        stage('Start stack') {
            steps {
                ansiColor('xterm') {
                    sh "${COMPOSE} up -d"
                }
            }
        }

        stage('Wait for health') {
            steps {
                // adjust sleep time if services need more boot time
                sh 'sleep 20'
            }
        }

        stage('Integration smoke test') {
            steps {
                ansiColor('xterm') {
                    // Replace with actual service health endpoint
                    sh 'curl -f http://localhost:8081/healthz'
                }
            }
        }

        stage('Push images (optional)') {
            when {
                branch 'master'
            }
            steps {
                ansiColor('xterm') {
                    sh "${COMPOSE} push"
                }
            }
        }
    }

    post {
        always {
            ansiColor('xterm') {
                sh "${COMPOSE} down || true"
            }
            echo "✅ Cleanup finished"
        }
        success {
            echo "🎉 Pipeline succeeded"
        }
        failure {
            echo "❌ Pipeline failed"
        }
    }
}
