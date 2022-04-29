#!/usr/bin/env bash

reconcile () {
	namespace=$1
	name=$2

	echo "[$namespace/$name] Reconciling"
	loadtest=$(kubectl get loadtest -ojson -n $namespace $name)

	#ltlabel=$(echo $loadtest | jq -r '.metadata.labels.loadtest')
	#[[ "$ltlabel" != "yes" ]] && echo "Not loadtest, ignoring" && exit 0

	requests=$(echo $loadtest | jq -r '.status.requests')
	[[ "$requests" != "null" ]] && echo "Load test already ran, ignoring" && exit 0

	method=$(echo $loadtest | jq -r '.spec.method')
	address=$(echo $loadtest | jq -r '.spec.address')
	duration=$(echo $loadtest | jq -r '.spec.duration')

	echo "[$namespace/$name] Run load test (method=$method, address=$address, duration=$duration)"
	results=$(echo "$method $address" | vegeta attack --duration=$duration | vegeta report --type json | jq)

	requests=$(echo $results | jq -r '.requests')
	echo "[$namespace/$name] Requests = $requests"
	kubectl patch loadtest -n $namespace $name --type='json' --patch='[{"op": "replace", "path": "/status", "value": {"requests":'$requests'}}]'
}

# The function needs to be exported so that is can be called with xargs.
export -f reconcile

kubectl get loadtest --all-namespaces --watch --output go-template='{{printf "%s %s\n" .metadata.namespace .metadata.name}}' | xargs -L1 bash -c 'reconcile "$@"' _

