pkill cloudprober
cloudprober -config_file cloudprober.cfg &
sleep 2
/app/service-prober &> sp.log