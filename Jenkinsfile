pipeline {
  agent any
  options { timestamps(); ansiColor('xterm') }

  environment {
    // Compose file that builds and runs db + usersvc + productsvc + ordersvc
    COMPOSE = "docker compose -f services/compose.yaml"
    DOCKER_BUILDKIT = "1"
    // On Docker Desktop, containers can reach the host via host.docker.internal
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
        sh """#!/bin/sh -e
          ${COMPOSE} build --no-cache
        """
      }
    }

    stage('Unit tests') {
      steps {
        sh """#!/bin/sh -e
          # Run go tests inside built images (no need to install Go on Jenkins)
          ${COMPOSE} run --rm usersvc go test ./...
          ${COMPOSE} run --rm productsvc go test ./...
          ${COMPOSE} run --rm ordersvc go test ./...
        """
      }
    }

    stage('Start stack') {
      steps {
        sh """#!/bin/sh -e
          ${COMPOSE} up -d
        """
      }
    }

    stage('Wait for health') {
      steps {
        sh """#!/bin/sh
          set -e

          wait_up () {
            local url="$1"
            echo "Waiting for: $url"
            for i in $(seq 1 60); do
              if curl -fsS "$url" >/dev/null 2>&1; then
                echo "OK: $url"
                return 0
              fi
              sleep 2
            done
            echo "TIMEOUT: $url"
            return 1
          }

          # hit published ports on the Docker host from inside Jenkins container
          wait_up http://$HOST:8081/healthz
          wait_up http://$HOST:8082/healthz
          wait_up http://$HOST:8083/healthz
        """
      }
    }

    stage('Integration smoke test') {
      steps {
        sh """#!/bin/sh -e

          curl -fsS -X POST "http://$HOST:8081/users" \
            -H "Content-Type: application/json" \
            -d '{ "name":"CI User", "email":"ci@demo.local" }'

          curl -fsS -X POST "http://$HOST:8082/products" \
            -H "Content-Type: application/json" \
            -d '{ "name":"Laptop", "price":75000 }'

          curl -fsS -X POST "http://$HOST:8083/orders" \
            -H "Content-Type: application/json" \
            -d '{ "userId":1, "productId":1, "quantity":2 }'
        """
      }
    }

    // OPTIONAL: push images to Docker Hub (turn on when needed)
    stage('Push images (optional)') {
      when { expression { return false } } // set to true to enable
      steps {
        withCredentials([usernamePassword(credentialsId: 'dockerhub', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
          sh """#!/bin/sh -e
            echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin

            # Tag & push
            docker tag services-usersvc      $DOCKER_USER/usersvc:latest
            docker tag services-productsvc   $DOCKER_USER/productsvc:latest
            docker tag services-ordersvc     $DOCKER_USER/ordersvc:latest

            docker push $DOCKER_USER/usersvc:latest
            docker push $DOCKER_USER/productsvc:latest
            docker push $DOCKER_USER/ordersvc:latest
          """
        }
      }
    }
  }

  post {
    always {
      // For CI, keep it clean
      sh """#!/bin/sh
        ${COMPOSE} down -v
      """
    }
    success { echo '✅ Pipeline successful' }
    failure { echo '❌ Pipeline failed' }
  }
}
