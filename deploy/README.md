### galaxy-s3-gateway自动化部署

#### 说明

galaxy-s3-gateway自动化部署用于自动在Vagrant虚拟机上安装galaxy-s3-gateway及其相关依赖，主要包括：

* galaxy-s3-gateway: 默认安装路径/home/$account/galaxy_s3_gateway
* mongodb: 通过service方式安装和启动

#### 步骤

1. 环境准备: cd deploy/ansilbe & bash prepare.sh, 这将准备好自动部署依赖的所有代码/第三方依赖等
2. 启动vagrant虚拟机: cd deploy/ & vagrant up --no-provision, 根据配置自动启动galaxy-s3-gateway-01
3. provision vagrant虚拟机: cd deploy/ansilbe & vagrant provision galaxy-s3-gateway-01

#### 注意事项

* 在MacOS下启动vagrant虚拟机时设置虚拟机ip的时候发现ip最后一个部分超过300好像就会启动失败
