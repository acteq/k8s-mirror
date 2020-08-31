# k8s-mirror-webhook

使用 k8s的时候，常常遇到镜像下载不了的问题，可以改从镜像站点下载。`k8s-mirror-webhook`是一个webhook，它监听pod和deployment的部署，把pod或deployment中的image改为镜像站点中的image。

# 配置

在 deploy/deployment.yml 文件中, 增加你需要的改写的原站点和对应的镜像站点，如下所示:

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

