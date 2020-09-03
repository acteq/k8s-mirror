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

# 测试

拉取 k8s.gcr.io/busybox 镜像

    kubectl apply -f test/pod.yml
    kubectl apply -f test/deployment.yml

# 使用 kubeadm 安装 k8s 集群时，访问国内镜像

使用 kubeadm 安装kubernetes集群时，默认也会到k8s.gcr.io 拉取镜像，幸好可以另行指定kubeadmin的镜像仓库地址， 比如阿里云镜像仓库地址 registry.aliyuncs.com/google_containers。

    $ kubeadm init \
    --apiserver-advertise-address=10.0.52.13 \
    --image-repository registry.aliyuncs.com/google_containers \
    --kubernetes-version v1.13.3 \
    --service-cidr=10.1.0.0/16 \
    --pod-network-cidr=10.244.0.0/16

