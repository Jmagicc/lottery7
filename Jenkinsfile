pipeline {
    agent any

    environment {
        // Harbor 相关配置
        HARBOR_HOST = '119.91.228.85:6656'
        HARBOR_REPO = 'library'
        HARBOR_USER = 'fucewei98'
        HARBOR_PASSWD = 'Asdmjh5452831'
        TASK_NAME = 'lottery7'
    }

    parameters {
        string(name: 'TAG', defaultValue: 'latest', description: 'Git tag or branch to build')
    }

    stages {
        stage('拉取Git代码') {
            steps {
                script {
                    // 使用参数化构建的 TAG 值
                    def tag = params.TAG

                    // 从 Gitee 拉取代码
                    checkout([
                        $class: 'GitSCM',
                        branches: [[name: "refs/heads/main"]], // 直接拉取 main 分支
                        userRemoteConfigs: [[
                            url: 'https://github.com/Jmagicc/lottery7'
                        ]]
                    ])
                }
            }
        }

        stage('构建Docker镜像') {
            steps {
                script {
                    // 使用 Jenkins Job 名称和 Tag 构建 Docker 镜像
                    sh "docker build -t ${JOB_NAME}:${params.TAG} ."
                }
            }
        }

        stage('推送镜像到Harbor') {
            steps {
                script {
                    // 登录 Harbor
                    sh "docker login -u ${HARBOR_USER} -p ${HARBOR_PASSWD} ${HARBOR_HOST}"

                    // 打 Tag
                    sh "docker tag ${JOB_NAME}:${params.TAG} ${HARBOR_HOST}/${HARBOR_REPO}/${TASK_NAME}:${params.TAG}"

                    // 推送镜像
                    sh "docker push ${HARBOR_HOST}/${HARBOR_REPO}/${TASK_NAME}:${params.TAG}"
                }
            }
        }
        stage('基于Harbor部署工程') {
            steps {
                script {
                    // 1. 停止并删除旧容器
                    sh """
                        docker stop ${TASK_NAME} || true
                        docker rm ${TASK_NAME} || true
                    """

                    // 2. 拉取镜像到宿主机
                    sh "docker pull ${HARBOR_HOST}/${HARBOR_REPO}/${TASK_NAME}:${params.TAG}"

                    // 3. 运行新容器
                    sh """
                        docker run -d \
                            --name ${TASK_NAME} \
                            --restart always \
                            -p 10025:10025 \
                            ${HARBOR_HOST}/${HARBOR_REPO}/${TASK_NAME}:${params.TAG}
                    """
                }
            }
        }
    }

    post {
        success {
            echo 'Pipeline 执行成功！'
        }
        failure {
            echo 'Pipeline 执行失败！'
        }
    }
}
