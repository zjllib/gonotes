package main

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/net/context"
	"gonotes/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// 业务实现方法的容器
type server struct{}

// 为server定义 DoMD5 方法 内部处理请求并返回结果
// 参数 (context.Context[固定], *proto.Req[相应接口定义的请求参数])
// 返回 (*proto.Res[相应接口定义的返回参数，必须用指针], error)
func (s *server) DoMD5(ctx context.Context, in *proto.Req) (*proto.Res, error) {
	fmt.Println("MD5方法请求JSON:" + in.JsonStr)
	return &proto.Res{BackJson: "MD5 :" + fmt.Sprintf("%x", md5.Sum([]byte(in.JsonStr)))}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8028")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	rpcService := grpc.NewServer() //创建gRPC服务

	/**注册接口服务
	 * 以定义proto时的service为单位注册，服务中可以有多个方法
	 * (proto编译时会为每个service生成Register***Server方法)
	 * 包.注册服务方法(gRpc服务实例，包含接口方法的结构体[指针])
	 */
	proto.RegisterWaiterServer(rpcService, &server{})
	/**如果有可以注册多个接口服务,结构体要实现对应的接口方法
	 * user.RegisterLoginServer(s, &server{})
	 * minMovie.RegisterFbiServer(s, &server{})
	 */
	// 在gRPC服务器上注册反射服务
	reflection.Register(rpcService)

	// 将监听交给gRPC服务处理
	err = rpcService.Serve(listen)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
