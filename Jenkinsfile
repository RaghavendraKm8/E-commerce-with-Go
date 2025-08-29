pipeline {
    agent any

    environment {
        REGISTRY = "docker.io/raghavendrakm"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build Services') {
            parallel {
                stage('Build Order Service') {
                    steps {
                        bat '''
                            cd ordersvc
                            go build -o ordersvc.exe
                        '''
                    }
                }
                stage('Build Product Service') {
                    steps {
                        bat '''
                            cd productsvc
                            go build -o productsvc.exe
                        '''
                    }
                }
                stage('Build User Service') {
                    steps {
                        bat '''
                            cd usersvc
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
                            cd ordersvc
                            go test ./...
                        '''
                    }
                }
                stage('Test Product Service') {
                    steps {
                        bat '''
                            cd productsvc
                            go test ./...
                        '''
                    }
                }
                stage('Test User Service') {
                    steps {
                        bat '''
                            cd usersvc
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
                        bat '''
                            docker build -t %REGISTRY%/ordersvc:latest ./ordersvc
                        '''
                    }
                }
                stage('Product Service Image') {
                    steps {
                        bat '''
                            docker build -t %REGISTRY%/productsvc:latest ./productsvc
                        '''
                    }
                }
                stage('User Service Image') {
                    steps {
                        bat '''
                            docker build -t %REGISTRY%/usersvc:latest ./usersvc
                        '''
                    }
                }
            }
        }

        stage('Push Images') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: '7022835052',   // ✅ Your actual Jenkins credential ID
                    usernameVariable: 'DOCKER_HUB_USER',
                    passwordVariable: 'DOCKER_HUB_PASS'
                )]) {
                    bat '''
                        echo %DOCKER_HUB_PASS% | docker login -u %DOCKER_HUB_USER% --password-stdin
                        docker push %REGISTRY%/ordersvc:latest
                        docker push %REGISTRY%/productsvc:latest
                        docker push %REGISTRY%/usersvc:latest
                    '''
                }
            }
        }
    }
}
