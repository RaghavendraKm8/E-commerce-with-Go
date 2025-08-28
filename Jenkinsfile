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
            parallel {
                stage('usersvc') {
                    steps {
                        bat '''
                          echo === Running usersvc tests ===
                          docker run --rm -v "%WORKSPACE%\\services\\usersvc:/app" -w /app golang:1.22 go test ./...
                        '''
                    }
                }
                stage('productsvc') {
                    steps {
                        bat '''
                          echo === Running productsvc tests ===
                          docker run --rm -v "%WORKSPACE%\\services\\productsvc:/app" -w /app golang:1.22 go test ./...
                        '''
                    }
                }
                stage('ordersvc') {
                    steps {
                        bat '''
                          echo === Running ordersvc tests ===
                          docker run --rm -v "%WORKSPACE%\\services\\ordersvc:/app" -w /app golang:1.22 go test ./...
                        '''
                    }
                }
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

                  for /l %%i in (1,1,30) do (
                    curl -s http://localhost:8081/healthz && curl -s http://localhost:8082/healthz && curl -s http://localhost:8083/healthz && exit 0
                    echo Still waiting...
                    timeout /t 2 >nul
                  )
                  echo One or more services failed to become healthy
                  exit 1
                '''
            }
        }

        stage('Integration smoke test') {
            steps {
                bat '''
                  echo Running smoke tests...
                  curl -f http://localhost:8081/healthz || exit 1
                  curl -f http://localhost:8082/healthz || exit 1
                  curl -f http://localhost:8083/healthz || exit 1
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
