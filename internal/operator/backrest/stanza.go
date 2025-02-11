package backrest

/*
 Copyright 2019 - 2021 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

import (
	"context"

	"github.com/qingcloud/postgres-operator/internal/config"
	"github.com/qingcloud/postgres-operator/internal/kubeapi"
	"github.com/qingcloud/postgres-operator/internal/operator"
	"github.com/qingcloud/postgres-operator/internal/util"
	crv1 "github.com/qingcloud/postgres-operator/pkg/apis/qingcloud.com/v1"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CleanStanzaCreateResources deletes any existing stanza-create pgtask and job.  Useful during a
// restore when an existing stanza-create pgtask or Job might still be present from initial
// creation of the cluster.
func CleanStanzaCreateResources(namespace, clusterName string, clientset kubeapi.Interface) error {
	ctx := context.TODO()

	resourceName := clusterName + "-" + crv1.PgtaskBackrestStanzaCreate

	if err := clientset.QingcloudV1().Pgtasks(namespace).
		Delete(ctx, resourceName, metav1.DeleteOptions{}); err != nil && !kerrors.IsNotFound(err) {
		return err
	}

	// job name is the same as the task name
	deletePropagation := metav1.DeletePropagationBackground
	if err := clientset.BatchV1().Jobs(namespace).Delete(ctx, resourceName,
		metav1.DeleteOptions{
			PropagationPolicy: &deletePropagation,
		}); err != nil && !kerrors.IsNotFound(err) {
		return err
	}

	return nil
}

func StanzaCreate(namespace, clusterName string, clientset kubeapi.Interface) {
	ctx := context.TODO()
	taskName := clusterName + "-" + crv1.PgtaskBackrestStanzaCreate

	// look up the backrest-repo pod name
	selector := config.LABEL_PG_CLUSTER + "=" + clusterName + "," + config.LABEL_PGO_BACKREST_REPO + "=true"
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if len(pods.Items) != 1 {
		log.Errorf("pods len != 1 for cluster %s", clusterName)
		return
	}
	if err != nil {
		log.Error(err)
		return
	}

	podName := pods.Items[0].Name

	// get the cluster to determine the proper storage type
	cluster, err := clientset.QingcloudV1().Pgclusters(namespace).Get(ctx, clusterName, metav1.GetOptions{})
	if err != nil {
		return
	}

	// create the stanza-create task
	spec := crv1.PgtaskSpec{}
	spec.Name = taskName

	jobName := clusterName + "-" + crv1.PgtaskBackrestStanzaCreate

	spec.TaskType = crv1.PgtaskBackrest
	spec.Parameters = make(map[string]string)
	spec.Parameters[config.LABEL_JOB_NAME] = jobName
	spec.Parameters[config.LABEL_PG_CLUSTER] = clusterName
	spec.Parameters[config.LABEL_POD_NAME] = podName
	spec.Parameters[config.LABEL_CONTAINER_NAME] = "qingcloud-pgbackrest-repo"
	// pass along the appropriate image prefix for the backup task
	// this will be used by the associated backrest job
	spec.Parameters[config.LABEL_IMAGE_PREFIX] = util.GetValueOrDefault(cluster.Spec.CCPImagePrefix, operator.Pgo.Cluster.CCPImagePrefix)
	spec.Parameters[config.LABEL_BACKREST_COMMAND] = crv1.PgtaskBackrestStanzaCreate
	// Get 'true' or 'false' for setting the pgBackRest S3 verify TLS value
	spec.Parameters[config.LABEL_BACKREST_S3_VERIFY_TLS] = operator.GetS3VerifyTLSSetting(cluster)

	// Handle stanza creation for a standby cluster, which requires some
	// additional consideration.
	// Since the primary will not be directly accessible to the standby cluster,
	// ensure the stanza created in offline mode
	if cluster.Spec.Standby {
		spec.Parameters[config.LABEL_BACKREST_OPTS] = "--no-online"
	}

	newInstance := &crv1.Pgtask{
		ObjectMeta: metav1.ObjectMeta{
			Name: taskName,
		},
		Spec: spec,
	}

	newInstance.ObjectMeta.Labels = make(map[string]string)
	newInstance.ObjectMeta.Labels[config.LABEL_PG_CLUSTER] = clusterName

	_, err = clientset.QingcloudV1().Pgtasks(namespace).Create(ctx, newInstance, metav1.CreateOptions{})
	if err != nil {
		log.Error(err)
	}
}
