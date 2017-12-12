# Introduction section
%define name galaxy-s3-gateway
%define version RELEASE_VERSION
%define release RELEASE_OS
%define arch RELEASE_ARCH

Name: %{name}
Summary: Galaxy AWS-S3 Gateway
Version: %{version}
Release: %{release}
BuildArch: %{arch}
# Group: System Environment/Libraries

# company info
vendor: HangZhou Galaxy Ltd.
Packager: tracymacding <tracymacding@gmail.com>
# URL: homepage  
# Copyright:
License: GPLv2

# spcefic source files, use %{version}
# service config
SOURCE: %{name}-%{version}.tar.bz2 

# pathes section:no patches 
# spec Requirements 
# see subpackage

# spec Provides 
# Provides:
# see subpackage

# spec Conflicts
# Conflicts:

# spec Obsoletes
# Obsoletes:

%description
Galaxy S3 Gateway is s3 front-end from galaxy LTD. 

# prep section
%prep
%setup -q

# build section
# now need no build section, because we package the binaries
# %build

# install section
# this is important, here not use make install, use install
%install

# mkdir 
mkdir -p ${RPM_BUILD_ROOT}/lib/systemd/system
mkdir -p ${RPM_BUILD_ROOT}/usr/bin
mkdir -p ${RPM_BUILD_ROOT}/usr/script/galaxy-s3-gateway

# copy all file to dest
install -m 0755 -t ${RPM_BUILD_ROOT}/lib/systemd/system lib/systemd/system/*
install -m 0755 -t ${RPM_BUILD_ROOT}/usr/bin usr/bin/*
install -m 0755 -t ${RPM_BUILD_ROOT}/usr/script/galaxy-s3-gateway usr/scripts/galaxy-s3-gateway/*

# clean section
%clean
#rm -rf $RPM_BUILD_ROOT

# files section
# list the files to go into the binary RPM, along with defined file attributes
# see subpackage

#%post -p /sbin/ldconfig
#%postun -p /sbin/ldconfig

# create Subpackages, name derive to galaxyfs-%{packagename}
# all subpackages section should provide name

#%package -n galaxy-s3-gateway
#Summary: Galaxy AWS-S3 gateway

%description -n galaxy-s3-gateway
Galaxy AWS-S3 gateway

# Requires: galaxyfs-common >= %{version}
# Requires: curl
# Requires: rrdtool
# Requires: net-snmp-libs

%pre -n galaxy-s3-gateway
if [ ! -d "/opt/galaxy" ]; then
     mkdir -m 755 /opt/galaxy
fi

if [ ! -d "/opt/galaxy/galaxy-s3-gw" ]; then
     mkdir -m 755 -p /opt/galaxy/galaxy-s3-gw
fi
if [ ! -d "/opt/galaxy/galaxy-s3-gw/bin" ]; then
     mkdir -m 755 -p /opt/galaxy/galaxy-s3-gw/bin
fi

%files -n galaxy-s3-gateway
%defattr(755,root,root)
/lib/systemd/system/s3-gateway.service
/usr/bin/galaxy-s3-gateway
/usr/script/galaxy-s3-gateway/run.sh
/usr/script/galaxy-s3-gateway/galaxy-s3-gateway.cfg

%post -n galaxy-s3-gateway
cp -d /usr/script/galaxy-s3-gateway/run.sh /opt/galaxy/galaxy-s3-gw/bin/run.sh
cp -d /usr/script/galaxy-s3-gateway/galaxy-s3-gateway.cfg /opt/galaxy/galaxy-s3-gw/bin/galaxy-s3-gateway.cfg
cp -d /usr/bin/galaxy-s3-gateway /opt/galaxy/galaxy-s3-gw/bin/galaxy-s3-gateway
systemctl enable s3-gateway
#/sbin/chkconfig metanode on

%preun -n galaxy-s3-gateway
if [ "$1" = 0 ]; then
  systemctl stop s3-gateway > /dev/null 2>&1
  systemctl disable s3-gateway
fi

%changelog
