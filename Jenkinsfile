pipeline {
    agent any

    environment {
        DOCKER_CREDENTIALS = credentials('7022835052')
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'master',
                    url: 'https://github.com/RaghavendraKm8/E-commerce-with-Go.git'
            }
        }

        stage('Build Services') {
            parallel {
                stage('Build Order Service') {
                    steps {
                        bat '''
                        cd services\\ordersvc
                        go build -o ordersvc.exe
                        '''
                    }
                }
                stage('Build Product Service') {
                    steps {
                        bat '''
                        cd services\\productsvc
                        go build -o productsvc.exe
                        '''
                    }
                }
                stage('Build User Service') {
                    steps {
                        bat '''
                        cd services\\usersvc
                        go build -o usersvc.exe
                        '''
                    }
                }
            }
        }

        stage('Run Tests') {
            parallel {
                stage('Test Order Service') {
                    steps {
                        bat '''
                        cd services\\ordersvc
                        go test ./...
                        '''
                    }
                }
                stage('Test Product Service') {
                    steps {
                        bat '''
                        cd services\\productsvc
                        go test ./...
                        '''
                    }
                }
                stage('Test User Service') {
                    steps {
                        bat '''
                        cd services\\usersvc
                        go test ./...
                        '''
                    }
                }
            }
        }

        stage('Build Docker Images') {
            parallel {
                stage('Order Service Image') {
                    steps {
                        bat 'docker build -t ordersvc:latest ./services/ordersvc'
                    }
                }
                stage('Product Service Image') {
                    steps {
                        bat 'docker build -t productsvc:latest ./services/productsvc'
                    }
                }
                stage('User Service Image') {
                    steps {
                        bat 'docker build -t usersvc:latest ./services/usersvc'
                    }
                }
            }
        }

        stage('Login to DockerHub') {
            steps {
                bat """
                echo ${DOCKER_CREDENTIALS_PSW} | docker login -u ${DOCKER_CREDENTIALS_USR} --password-stdin
                """
            }
        }

        stage('Push Images') {
            steps {
                bat '''
                docker tag ordersvc:latest RaghavendraKm/ordersvc:latest
                docker tag productsvc:latest RaghavendraKm/productsvc:latest
                docker tag usersvc:latest RaghavendraKm/usersvc:latest

                docker push RaghavendraKm/ordersvc:latest
                docker push RaghavendraKm/productsvc:latest
                docker push RaghavendraKm/usersvc:latest
                '''
            }
        }
    }
}
