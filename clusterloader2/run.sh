workers=10
replicas=3
extra_flag="_hn"
sched_qps=800_800
ctrl_qps=800_800

tag="${sched_qps}_s_${ctrl_qps}_c_${workers}_nodes_${replicas}_pods_typha${extra_flag}"

outdir="$HOME/public/k8s-performance-test/reports"
reportdir="report"

mkdir -p "$outdir/$tag"

for h in `seq 1 1`;

do ( 
# clusterloader --provider local  --report-dir "$outdir/$reportdir" --nodes $workers --kubeconfig ~/.kube/config --testconfig ../config/tests/cl2-density-config.yaml --etcd-key /atlas-home/1/avolio/k8s/etcd.p1/etcd-client.key --etcd-certificate /atlas-home/1/avolio/k8s/etcd.p1/etcd-client.crt  --masterip pc-tdq-k8m-master.cern.ch --mastername pc-tdq-k8m-master.cern.ch;

./clusterloader --provider local  --report-dir "$outdir/$reportdir" --nodes 10 --kubeconfig "$HOME/public/config-tbed" --testconfig "$HOME/public/k8s-performance-tests/config/tests/cl2-density-config.yaml" --masterip pc-tbed-k8m-12.cern.ch --mastername pc-tbed-k8m-12.cern.ch

mv "$outdir/$reportdir/junit.xml" "$outdir/$tag/junit_${tag}_${h}.xml";
mv timeline.csv "$outdir/$tag/timeline_${tag}_${h}.csv";

echo "Sleeping..."
sleep 180;

);

done
