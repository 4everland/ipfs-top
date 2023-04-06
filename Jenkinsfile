pipeline {
    agent any

    environment {
        DOCKER_IMAGE_NAME = 'ipfs-servers'
        DOCKER_IMAGE = ''
        DEPLOY_IMAGE_NAME = ''
    }

    stages {
        stage("Chose Environment") {
            parallel {
                stage('Dev') {
                    when {
                        environment name: 'APP_ENV', value: 'DEV'
                    }

                    stages {
                        stage('build image') {
                            steps {
                                script {
                                    DOCKER_IMAGE = docker.build(env.DOCKER_NAMESPACE + "/" + DOCKER_IMAGE_NAME + "-" + env.APP_NAME, "--build-arg APP_NAME="+env.APP_NAME + " .")
                                    DEPLOY_IMAGE_NAME = DOCKER_NAMESPACE + "/" + DOCKER_IMAGE_NAME + "-" + env.APP_NAME + ":" + env.GIT_COMMIT.substring(0, 8)
                                }
                            }
                        }

                        stage('push image') {
                            steps {
                                script {
                                    docker.withRegistry(env.DOCKER_HUB_DEV_URL, env.DOCKER_HUB_DEV_CREDENTIAL) {
                                        DOCKER_IMAGE.push(env.GIT_COMMIT.substring(0, 8))
                                    }
                                }
                            }
                        }

                        stage('remove image') {
                            steps {
                                script {
                                    DOCKER_HUB_HOST = env.DOCKER_HUB_DEV_URL.replaceAll("https://", "").replaceAll("http://", "")
                                    sh "docker rmi ${DOCKER_NAMESPACE}/${DOCKER_IMAGE_NAME}-${APP_NAME}"
                                    sh 'docker rmi ' + DOCKER_HUB_HOST + DEPLOY_IMAGE_NAME
                                    sh 'docker image prune -f --filter label=stage=ipfs-servers'
                                }
                            }
                        }
                    }
                }

                stage('Prod') {
                    when {
                        environment name: 'APP_ENV', value: 'PROD'
                    }
                    stages {
                        stage('build image') {
                            steps {
                                script {
                                    DOCKER_IMAGE = docker.build(env.DOCKER_NAMESPACE + "/" + DOCKER_IMAGE_NAME + "-" + env.APP_NAME, "--build-arg APP_NAME="+env.APP_NAME + " .")
                                    DEPLOY_IMAGE_NAME = env.DOCKER_NAMESPACE + "/" + DOCKER_IMAGE_NAME + "-" + env.APP_NAME + ":" + env.branch
                                }
                            }
                        }

                        stage('push image') {
                            steps {
                                script {
                                    docker.withRegistry(env.DOCKER_HUB_PROD_URL, env.DOCKER_HUB_PROD_CREDENTIAL) {
                                        DOCKER_IMAGE.push(env.branch)
                                    }
                                }
                            }
                        }

                        stage('remove image') {
                            steps {
                                script {
                                    DOCKER_HUB_HOST = env.DOCKER_HUB_PROD_URL.replaceAll("https://", "").replaceAll("http://", "")
                                    sh "docker rmi ${DOCKER_NAMESPACE}/${DOCKER_IMAGE_NAME}-${APP_NAME}"
                                    sh 'docker rmi ' + DOCKER_HUB_HOST + DEPLOY_IMAGE_NAME
                                    sh 'docker image prune -f --filter label=stage=ipfs-servers'
                                }
                            }
                        }
                    }
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    def body = '''
                    {
                        "spec": {
                            "template": {
                                "spec": {
                                    "containers": [
                                        {
                                            "name": "''' + DOCKER_IMAGE_NAME + "-" + env.APP_NAME + '''",
                                            "image": "''' + DOCKER_HUB_HOST + DEPLOY_IMAGE_NAME + '''"
                                        }
                                    ]
                                }
                            }
                        }
                    }
                   '''

                    response = httpRequest consoleLogResponseBody: true, contentType: 'NOT_SET',
                    httpMode: 'PATCH', requestBody: body, url: env.K8S_DEPLOY_API_URL, validResponseCodes: '200',
                    customHeaders: [[name: 'Content-Type', value: "application/strategic-merge-patch+json"]]
                }
            }
        }
    }
}
