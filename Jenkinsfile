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

        stage('Login to DockerHub') {
            steps {
                withCredentials([usernamePassword(credentialsId: '7022835052', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    bat """
                    echo %DOCKER_PASS% | docker login -u %DOCKER_USER% --password-stdin
                    """
                }
            }
        }

        stage('Push Images') {
            steps {
                bat '''
                docker tag ordersvc:latest raghavendrakm/ordersvc:latest
                docker tag productsvc:latest raghavendrakm/productsvc:latest
                docker tag usersvc:latest raghavendrakm/usersvc:latest

                docker push raghavendrakm/ordersvc:latest
                docker push raghavendrakm/productsvc:latest
                docker push raghavendrakm/usersvc:latest
                '''
            }
        }
    }
}
