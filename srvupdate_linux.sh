go install stonesrv
cp /home/flynnyal/application/go/bin/stonesrv /home/flynnyal/application/go/src/stonesrv/package
cp /home/flynnyal/application/go/src/stonesrv/language /home/flynnyal/application/go/src/stonesrv/package
cp /home/flynnyal/application/go/src/stonesrv/updates /home/flynnyal/application/go/src/stonesrv/package
rm /home/flynnyal/application/go/src/stonesrv/package/language/*.go
scp -r /home/flynnyal/application/go/src/stonesrv/package root@202.182.127.17:/root/application/go/src/stonesrv
