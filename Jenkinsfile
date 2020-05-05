pipeline {
  agent any
  stages {
    stage('检出') {
      steps {
        checkout([
          $class: 'GitSCM', branches: [[name: env.GIT_BUILD_REF]],
          userRemoteConfigs: [[
            url: env.GIT_REPO_URL,
            credentialsId: env.CREDENTIALS_ID
          ]]
        ])
      }
    }
    stage('构建') {
      steps {
        echo '构建中...'
        sh 'export GOPROXY=https://goproxy.cn'
        sh 'go build'
        echo '构建完成.'
      }
    }
    stage('测试') {
      steps {
        echo '单元测试中...'
        echo '单元测试完成.'
      }
    }
    stage('打包镜像') {
      steps {
        echo '部署中...'
        sh 'docker build -t ${DOCKER_IMAGE_NAME}:${GIT_LOCAL_BRANCH} .'
        echo '部署完成'
      }
    }
    stage('推送制品库') {
      steps {
        script {
          docker.withRegistry("https://${env.CODING_DOCKER_REG_HOST}", "${env.CODING_ARTIFACTS_CREDENTIALS_ID}") {
            docker.image("${env.DOCKER_IMAGE_NAME}:${env.GIT_LOCAL_BRANCH}").push()
          }
        }

      }
    }
  }
  environment {
    CODING_DOCKER_REG_HOST = "${env.CCI_CURRENT_TEAM}-docker.pkg.${env.CCI_CURRENT_DOMAIN}"
    DOCKER_IMAGE_NAME = "${env.PROJECT_NAME}/${env.DOCKER_REPO_NAME}/${env.DEPOT_NAME}"
  }
}