#!/bin/bash

source hack/common.sh

$OUTDIR_BIN/functests --ocs-registry-image="${FULL_IMAGE_NAME}" --local-storage-registry-image="${LOCAL_STORAGE_IMAGE_NAME}"  $@
if [ $? -ne 0 ]; then
	echo "dumping debug information"
	echo "--- PODS ----"
	oc get pods -n openshift-storage
	echo "---- PVs ----"
	oc get pv
	echo "--- StorageClasses ----"
	oc get storageclass --all-namespaces
	echo "--- StorageCluster ---"
	oc get storagecluster --all-namespaces -o yaml
	echo "--- CephCluster ---"
	oc get cephcluster --all-namespaces -o yaml
	echo "ERROR: Functest failed."
	exit 1
fi

