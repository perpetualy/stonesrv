go install stonesrv
cp /home/fy/application/go/bin/stonesrv /home/fy/application/go/src/stonesrv/package
cp -r /home/fy/application/go/src/stonesrv/language /home/fy/application/go/src/stonesrv/package
cp -r /home/fy/application/go/src/stonesrv/updates /home/fy/application/go/src/stonesrv/package
cp -r /home/fy/application/go/src/stonesrv/conf /home/fy/application/go/src/stonesrv/package
rm /home/fy/application/go/src/stonesrv/package/language/*.go
rm /home/fy/application/go/src/stonesrv/package/conf/*.go
scp -r /home/fy/application/go/src/stonesrv/package root@202.182.127.17:/root/application/go/src/stonesrv

