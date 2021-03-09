#!/usr/bin/env bash

echo "协议生成开始"

set -e

Platform=$(go env GOHOSTOS)

echo ${Platform}

# cellmesh服务绑定 cmprotogen
#go build -v -o ${CellMeshProtoGen} github.com/bobwong89757/cellmesh/tool/protogen

# 协议生成 protoplus

# pb插件
#go build -v -o ${GoGoFaster} github.com/gogo/protobuf/protoc-gen-gogofaster

# 路由工具 routegen
#go build -v -o ${RouteGen} github.com/bobwong89757/tool/routegen

OutputPath=../proto

echo "生成服务器协议的go消息..."
protoplus -package=proto -go_out=${OutputPath}/msgsvc_gen.go `source ./protolist.sh svc`

echo "生成服务器协议的消息绑定..."
cmprotogen -package=proto -cmgo_out=${OutputPath}/msgbind_gen.go `source ./protolist.sh all`

echo "生成客户端协议的protobuf proto文件..."
protoplus --package=proto -pb_out=${OutputPath}/clientmsg_gen.proto `source ./protolist.sh client`

echo "生成客户端协议的protobuf的go消息...."
protoc --plugin=protoc-gen-gogofaster --gogofaster_out=${OutputPath} --proto_path=${OutputPath}/ clientmsg_gen.proto

echo "完成协议生成"