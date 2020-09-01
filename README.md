# 介绍

`k8s-mirror-webhook`是一个webhook，部署到kubernetes集群后，它会监听 `pod` 和`deployment`的部署，把要拉取的image改为镜像站点中的image。

# 配置
 
在 deploy/deployment.yml 文件中, 设置image原站点和对应的镜像站点，如下所示:

    ---
    apiVersion: v1
    kind: ConfigMap
    metadata:
    name: mirror-webhook-config
    data:
    conf.hcl: |
        mirror {
            gcr.io = "registry.aliyuncs.com/google_containers"
            k8s.gcr.io = "registry.aliyuncs.com/google_containers"
        }

# 安装

    kubectl apply -f deploy/deployment.yml

