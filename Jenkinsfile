pipeline {
    agent any
   
    stages {
        stage ('build and push') {
            steps {
                checkout scmGit(branches: [[name: '*/main']], extensions: [], userRemoteConfigs: [[credentialsId: 'github', url: 'https://github.com/ququiz/scoring-service']])
                sh 'chmod 777 ./push.sh'
                sh './push.sh'
                sh 'docker stop ququiz-scoring-service && docker rm ququiz-scoring-service'
                sh 'docker rmi lintangbirdas/scoring-service:v1'
            }
        }
        stage ('docker compose up') {
            steps {
                build job: "ququiz-compose", wait: true
            }
        }
    }

}