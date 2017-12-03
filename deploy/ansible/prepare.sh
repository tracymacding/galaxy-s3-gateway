version='0.1.0'

git clone https://git.coding.net/Mr-x/galaxy-release.git

mkdir -p roles/galaxy-s3-gateway/files
cp galaxy-release/galaxy-s3-gateway/$version/galaxy-s3-gateway.tar.gz  roles/galaxy-s3-gateway/files/
cp ~/Downloads/go1.7.3.linux-amd64.tar.gz  roles/galaxy-s3-gateway/files/
