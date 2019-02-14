go install xmvideosrv
cp /home/flynnyal/application/go/bin/xmvideosrv /home/flynnyal/application/go/src/xmvideo/package
cp -r /home/flynnyal/application/go/src/xmvideo/language /home/flynnyal/application/go/src/xmvideo/package
cp -r /home/flynnyal/application/go/src/xmvideo/updates /home/flynnyal/application/go/src/xmvideo/package
cp -r /home/flynnyal/application/go/src/xmvideo/conf /home/flynnyal/application/go/src/xmvideo/package
rm /home/flynnyal/application/go/src/xmvideo/package/language/*.go
rm /home/flynnyal/application/go/src/xmvideo/package/conf/*.go
scp -r /home/flynnyal/application/go/src/xmvideo/package root@47.107.60.107:/root/xmvideo

