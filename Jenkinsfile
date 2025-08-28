pipeline {
    agent any

    environment {
        REGISTRY = "your-dockerhub-username"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Go Services') {
            steps {
                script {
                    def services = ["ordersvc", "productsvc", "usersvc"]

                    for (service in services) {
                        dir("services/${service}") {
                            sh "go mod tidy"
                            sh "go build -o service ."
                        }
                    }
                }
            }
        }

        stage('Build Docker Images') {
            steps {
                script {
                    def services = ["ordersvc", "productsvc", "usersvc"]

                    for (service in services) {
                        dir("services/${service}") {
                            sh "docker build -t ${REGISTRY}/${service}:latest ."
                        }
                    }
                }
            }
        }

        stage('Push Docker Images') {
            steps {
                script {
                    def services = ["ordersvc", "productsvc", "usersvc"]

                    docker.withRegistry('https://index.docker.io/v1/', 'dockerhub-credentials') {
                        for (service in services) {
                            sh "docker push ${REGISTRY}/${service}:latest"
                        }
                    }
                }
            }
        }
    }
}
