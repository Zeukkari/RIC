#!/bin/bash

IMAGELIBS=libwebp-dev libtiff5-dev libpng12-dev libjpeg-dev liblqr-1-0-dev
OPENCL=ocl-icd-opencl-dev opencl-headers ocl-icd-libopencl1
DEVEL=git-all gcc g++ binutils automake autoconf pkg-config

# Install the absolute requirements
apt-get update
apt-get install -y ${IMAGELIBS}
apt-get install -y ${OPENCL}
apt-get install -y ${DEVEL}

