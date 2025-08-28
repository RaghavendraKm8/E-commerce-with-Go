pipeline {
    agent any

    options {
        timestamps()
        ansiColor('xterm')
    }

    environment {
        DOCKER_BUILDKIT     = '1'
        DOCKER_COMPOSE_FILE = 'docker-compose.yml'
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build images') {
            steps {
                bat """
                echo Building Docker images...
                docker compose -f %DOCKER_COMPOSE_FILE% build --no-cache
                """
            }
        }

        stage('Unit tests') {
            steps {
                bat """
                echo Running Go unit tests (via Docker golang:1.22)...

                echo === usersvc ===
                docker run --rm -v "%cd%\\services\\usersvc:/app" -w /app golang:1.22 go test ./...

                echo === productsvc ===
                docker run --rm -v "%cd%\\services\\productsvc:/app" -w /app golang:1.22 go test ./...

                echo === ordersvc ===
                docker run --rm -v "%cd%\\services\\ordersvc:/app" -w /app golang:1.22 go test ./...
                """
            }
        }

        stage('Start stack') {
            steps {
                bat """
                echo Starting full stack...
                docker compose -f %DOCKER_COMPOSE_FILE% up -d --remove-orphans
                """
            }
        }

        stage('Wait for health') {
            steps {
                bat """
                echo Waiting for services to become healthy...
                timeout /t 20 >NUL
                """
            }
        }

        stage('Integration smoke test') {
            steps {
                bat """
                echo Running smoke test...
                curl -fsS http://localhost:8081/healthz || exit /b 1
                curl -fsS http://localhost:8082/healthz || exit /b 1
                curl -fsS http://localhost:8083/healthz || exit /b 1
                echo All services passed smoke test
                """
            }
        }
    }

    post {
        always {
            bat """
            docker compose -f %DOCKER_COMPOSE_FILE% down -v || echo Nothing to clean
            """
            echo "✅ Cleanup finished"
        }
        failure {
            echo "❌ Pipeline failed"
        }
        success {
            echo "🎉 Pipeline succeeded"
        }
    }
}
