pipeline {
    agent any

    options {
        timestamps()
        ansiColor('xterm')
    }

    environment {
        DOCKER_BUILDKIT = "1"
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
                    sh '''
                      echo "🐳 Building Docker images..."
                      docker compose -f docker-compose.yml build --no-cache
                    '''
                }
            }
        }

        stage('Unit tests') {
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "🔍 Running unit tests in each service..."

                      echo "➡️ Testing usersvc ..."
                      docker run --rm -v $PWD:/go-ecommerce -w /go-ecommerce/usersvc golang:1.22 go test ./...

                      echo "➡️ Testing productsvc ..."
                      docker run --rm -v $PWD:/go-ecommerce -w /go-ecommerce/productsvc golang:1.22 go test ./...

                      echo "➡️ Testing ordersvc ..."
                      docker run --rm -v $PWD:/go-ecommerce -w /go-ecommerce/ordersvc golang:1.22 go test ./...
                    '''
                }
            }
        }

        stage('Start stack') {
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "🚀 Starting full stack with docker-compose..."
                      docker compose -f docker-compose.yml up -d
                    '''
                }
            }
        }

        stage('Wait for health') {
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "⏳ Waiting for services to become healthy..."
                      sleep 20
                    '''
                }
            }
        }

        stage('Integration smoke test') {
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "🧪 Running integration tests..."

                      curl -f http://localhost:8081/healthz || exit 1
                      curl -f http://localhost:8082/healthz || exit 1
                      curl -f http://localhost:8083/healthz || exit 1

                      echo "✅ All services passed smoke test"
                    '''
                }
            }
        }

        stage('Push images (optional)') {
            when {
                expression { return false } // enable later if pushing to registry
            }
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "📦 Pushing Docker images to registry..."
                      docker tag go-ecommerce-usersvc myrepo/go-ecommerce-usersvc:latest
                      docker tag go-ecommerce-productsvc myrepo/go-ecommerce-productsvc:latest
                      docker tag go-ecommerce-ordersvc myrepo/go-ecommerce-ordersvc:latest
                      docker push myrepo/go-ecommerce-usersvc:latest
                      docker push myrepo/go-ecommerce-productsvc:latest
                      docker push myrepo/go-ecommerce-ordersvc:latest
                    '''
                }
            }
        }
    }

    post {
        always {
            ansiColor('xterm') {
                sh '''
                  docker compose -f docker-compose.yml down -v || true
                '''
            }
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
