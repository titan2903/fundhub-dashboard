pipeline {
    agent {
        label 'fire'
    }

    environment {
        WEBHOOK_URL_DISCORD = credentials('WEBHOOK_URL_DISCORD')
        DOCKER_USERNAME = credentials('DOCKER_USERNAME')
        DOCKER_PASSWORD = credentials('DOCKER_PASSWORD')
    }

    stages {
        stage('Unit Test'){
            agent {
                docker {
                    image 'golang:1.21.4-alpine3.18'
                    label 'fire'
                }
            }

            steps {
                sh 'GOCACHE=/tmp/ go test -v ./...'
            }   
        }

        stage('Build') {
            steps {
                echo "Building Apps"
                sh "docker build -t $DOCKER_USERNAME/fundhub-dashboard:${BUILD_NUMBER} ."
            }
        }

        stage('Deploy') {
             steps {
                echo "Pushing to DockerHub"
                sh "docker login -u $DOCKER_USERNAME -p ${DOCKER_PASSWORD}"
                sh "docker push $DOCKER_USERNAME/fundhub-dashboard:${BUILD_NUMBER}"
                sh 'docker rmi $(docker images -q)'
            }

            post {
                aborted {
                    echo "Post Aborted"
                    discordSend description: "Push Docker Image", footer: "Push image fundhub-dashboard:${BUILD_NUMBER} to DockerHub Status: Aborted", link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME, webhookURL: "$WEBHOOK_URL_DISCORD"
                }
                success {
                    echo "Post Success"
                    discordSend description: "Push Docker Image", footer: "Push image fundhub-dashboard:${BUILD_NUMBER} to DockerHub Status: Success", link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME, webhookURL: "$WEBHOOK_URL_DISCORD"
                }

                failure {
                    echo "Post Failure"
                    discordSend description: "Push Docker Image", footer: "Push image fundhub-dashboard:${BUILD_NUMBER} to DockerHub Status: Failure", link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME, webhookURL: "$WEBHOOK_URL_DISCORD"
                }
            }
        }
    }
    
    post {
        aborted {
            echo "Post Aborted"
            discordSend description: "Fundhub Dashboard Deployment", footer: "Fundhub Dashboard Staging Deployment Status: Aborted", link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME, webhookURL: "$WEBHOOK_URL_DISCORD"
        }

        success {
            echo "Post Success"
            discordSend description: "Fundhub Dashboard Deployment", footer: "Fundhub Dashboard Staging Deployment Status: Success", link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME, webhookURL: "$WEBHOOK_URL_DISCORD"
        }
        
        failure {
            echo "Post Failure"
            discordSend description: "Fundhub Dashboard Deployment", footer: "Fundhub Dashboard Staging Deployment Status: Failure", link: env.BUILD_URL, result: currentBuild.currentResult, title: JOB_NAME, webhookURL: "$WEBHOOK_URL_DISCORD"
        }
    }
}