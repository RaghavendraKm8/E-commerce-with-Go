pipeline {
    agent any

    environment {
        DOCKER_BUILDKIT = "1"
        COMPOSE_DOCKER_CLI_BUILD = "1"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build images') {
            steps {
                bat '''
                  echo Building Docker images...
                  docker compose -f docker-compose.yml build --no-cache
                '''
            }
        }

        stage('Unit tests') {
            steps {
                bat '''
                  echo Running Go unit tests (via Docker golang:1.22)...

                  echo === usersvc ===
                  docker run --rm -v "%WORKSPACE%\\services\\usersvc:/app" -w /app golang:1.22 go test ./...

                  echo === productsvc ===
                  docker run --rm -v "%WORKSPACE%\\services\\productsvc:/app" -w /app golang:1.22 go test ./...

                  echo === ordersvc ===
                  docker run --rm -v "%WORKSPACE%\\services\\ordersvc:/app" -w /app golang:1.22 go test ./...
                '''
            }
        }

        stage('Start stack') {
            steps {
                bat '''
                  echo Cleaning up old containers...
                  docker compose -f docker-compose.yml down -v || echo Nothing to stop

                  echo Starting full stack...
                  docker compose -f docker-compose.yml up -d --remove-orphans
                '''
            }
        }

        stage('Wait for health') {
            steps {
                bat '''
                  echo Waiting for services to be healthy...
                  for /L %%i in (1,1,30) do (
                    docker inspect --format "{{.State.Health.Status}}" usersvc 2>nul | findstr "healthy" >nul && ^
                    docker inspect --format "{{.State.Health.Status}}" productsvc 2>nul | findstr "healthy" >nul && ^
                    docker inspect --format "{{.State.Health.Status}}" ordersvc 2>nul | findstr "healthy" >nul && exit /b 0
                    timeout /t 5 >nul
                  )
                  echo ❌ Services did not become healthy in time
                  exit /b 1
                '''
            }
        }

        stage('Integration smoke test') {
            steps {
                bat '''
                  echo Running smoke test...
                  curl -f http://localhost:8081/healthz || exit 1
                '''
            }
        }
    }

    post {
        always {
            bat '''
              docker compose -f docker-compose.yml down -v || echo Nothing to clean
            '''
            echo "✅ Cleanup finished"
        }
        failure {
            echo "❌ Pipeline failed"
        }
    }
}
