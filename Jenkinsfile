pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                git branch: 'master',
                    credentialsId: '7022835052',
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

        stage('Push Images') {
            steps {
                bat '''
                docker tag ordersvc:latest your-dockerhub-username/ordersvc:latest
                docker tag productsvc:latest your-dockerhub-username/productsvc:latest
                docker tag usersvc:latest your-dockerhub-username/usersvc:latest

                docker push your-dockerhub-username/ordersvc:latest
                docker push your-dockerhub-username/productsvc:latest
                docker push your-dockerhub-username/usersvc:latest
                '''
            }
        }
    }
}
