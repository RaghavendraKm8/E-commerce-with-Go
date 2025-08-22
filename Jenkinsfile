pipeline {
  agent any
  options { timestamps() }

  environment {
    COMPOSE = "docker compose -f docker-compose.yml"
    DOCKER_BUILDKIT = "1"
    HOST = "host.docker.internal"
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
          sh """
            ${COMPOSE} build --no-cache
          """
        }
      }
    }

    stage('Unit tests') {
      steps {
        ansiColor('xterm') {
          sh '''
            echo "🔍 Running unit tests inside golang:1.22 container..."
            docker run --rm -v $PWD:/app -w /app golang:1.22 go test ./...
          '''
        }
      }
    }

    stage('Start stack') {
      steps {
        ansiColor('xterm') {
          sh """
            ${COMPOSE} up -d
          """
        }
      }
    }

    stage('Wait for health') {
      steps {
        ansiColor('xterm') {
          sh '''
            set -e
            wait_up () {
              local url="$1"
              echo "Waiting for: $url"
              for i in $(seq 1 60); do
                if curl -fsS "$url" >/dev/null 2>&1; then
                  echo "✅ OK: $url"
                  return 0
                fi
                sleep 2
              done
              echo "❌ TIMEOUT: $url"
              return 1
            }

            wait_up http://$HOST:8081/healthz
            wait_up http://$HOST:8082/healthz
            wait_up http://$HOST:8083/healthz
          '''
        }
      }
    }

    stage('Integration smoke test') {
      steps {
        ansiColor('xterm') {
          sh '''
            curl -fsS -X POST "http://$HOST:8081/users" \
              -H "Content-Type: application/json" \
              -d '{ "name":"CI User", "email":"ci@demo.local" }'

            curl -fsS -X POST "http://$HOST:8082/products" \
              -H "Content-Type: application/json" \
              -d '{ "name":"Laptop", "price":75000 }'

            curl -fsS -X POST "http://$HOST:8083/orders" \
              -H "Content-Type: application/json" \
              -d '{ "userId":1, "productId":1, "quantity":2 }'
          '''
        }
      }
    }

    stage('Push images (optional)') {
      when { expression { return false } } // set to true when needed
      steps {
        withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
          ansiColor('xterm') {
            sh '''
              echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin

              docker tag go-ecommerce-usersvc    $DOCKER_USER/usersvc:latest
              docker tag go-ecommerce-productsvc $DOCKER_USER/productsvc:latest
              docker tag go-ecommerce-ordersvc   $DOCKER_USER/ordersvc:latest

              docker push $DOCKER_USER/usersvc:latest
              docker push $DOCKER_USER/productsvc:latest
              docker push $DOCKER_USER/ordersvc:latest
            '''
          }
        }
      }
    }
  }

  post {
    always {
      ansiColor('xterm') {
        sh """
          ${COMPOSE} down -v || true
        """
      }
      echo "✅ Cleanup finished"
    }
    success { echo '🎉 Pipeline successful' }
    failure { echo '❌ Pipeline failed' }
  }
}
