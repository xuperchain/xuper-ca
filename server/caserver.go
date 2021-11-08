/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */
package server

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/xuperchain/xuper-ca/pb"
	"github.com/xuperchain/xuper-ca/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/xuperchain/xuper-ca/config"
	"github.com/xuperchain/xuper-ca/crypto"
	"github.com/xuperchain/xuper-ca/util"
)

type caServer struct{}

// 接口层签名校验, 检验的data根据接口不同而不同
func verifyRequest(sign *pb.Sign, data []byte) bool {
	return true
	if sign == nil {
		log.Warning("request sign is nil")
		return false
	}
	cryptoClient := crypto.GetCryptoClient()
	pubKey, err := cryptoClient.GetEcdsaPublicKeyFromJSON([]byte(sign.PublicKey))
	if err != nil {
		log.Debug("crypto GetEcdsaPublicKeyFromJSON error")
		return false
	}
	ok, err := cryptoClient.VerifyECDSA(pubKey, []byte(sign.Sign), []byte(string(data)+sign.Nonce))
	if err != nil {
		log.Debug("crypto VerifyECDSA error")
	}
	return ok
}

// 添加一个网络及其管理员
// 签名字段为address+net
// conf中配置的 ca 根地址可以操作该接口
func (ca *caServer) NetAdminEnroll(ctx context.Context, in *pb.EnrollNetRequest) (*pb.EnrollResponse, error) {
	if in.GetLogid() == "" {
		in.Logid = util.GetLogId()
	}
	log.AddHook(util.NewlogIdHook(in.GetLogid()))

	if ok := verifyRequest(in.Sign, []byte(in.Address+in.Net)); !ok {
		log.Error("NetAdminEnroll sign is not right")
		return &pb.EnrollResponse{
			Logid: in.Logid,
		}, ErrSign
	}

	if ok := service.CheckCaAdmin(in.Sign.Address); !ok {
		log.Error(ErrAuth, ", address:", in.Sign.Address)
		return &pb.EnrollResponse{
			Logid: in.Logid,
		}, ErrAuth
	}

	err := service.AddNetAdmin(in.Net, in.Address, in.Isgm)
	if err != nil {
		log.Error("AddNetAdmin add admin failed:", err)
		return &pb.EnrollResponse{
			Logid: in.Logid,
		}, ErrAddNetAndAdmin
	}
	log.Infof("add net admin success, net: %+v, address: %+v", in.Net, in.Address)
	return &pb.EnrollResponse{
		Logid: in.Logid,
	}, nil
}

// 在一个网络中授权节点
// 签名字段为address+net
// conf 配置中的ca根地址  和网络的管理员地址可以操作该接口
func (ca *caServer) NodeEnroll(ctx context.Context, in *pb.EnrollNodeRequest) (*pb.EnrollResponse, error) {
	if in.GetLogid() == "" {
		in.Logid = util.GetLogId()
	}
	log.AddHook(util.NewlogIdHook(in.GetLogid()))

	if ok := verifyRequest(in.Sign, []byte(in.Address+in.Net)); !ok {
		log.Error("NodeEnroll sign is not right")
		return &pb.EnrollResponse{
			Logid: in.Logid,
		}, ErrSign
	}

	if ok := service.CheckNetAdmin(in.Sign.Address, in.Net); !ok {
		log.Warning("address isn't net admin")
		return &pb.EnrollResponse{
			Logid: in.Logid,
		}, ErrAuth
	}

	err := service.AddNode(in.Net, in.AdminAddress, in.Address)
	if err != nil {
		log.Error("NodeEnroll add node failed:", err)
		return &pb.EnrollResponse{
			Logid: in.Logid,
		}, ErrAddNode
	}
	log.Infof("add node success, net: %+v, address: %+v", in.Net, in.Address)
	return &pb.EnrollResponse{
		Logid: in.Logid,
	}, nil
}

// 获取一个节点的证书
// 签名字段为address+net
// conf中配置的ca根地址、网络的管理员地址 及节点的地址 可以操作该接口
func (ca *caServer) GetCurrentCert(ctx context.Context, in *pb.CurrentCertRequest) (*pb.CurrentCertResponse, error) {
	if in.GetLogid() == "" {
		in.Logid = util.GetLogId()
	}
	log.AddHook(util.NewlogIdHook(in.GetLogid()))

	if ok := verifyRequest(in.Sign, []byte(in.Address+in.Net)); !ok {
		log.Error("GetCurrentCert sign is not right")
		return &pb.CurrentCertResponse{
			Logid: in.Logid,
		}, ErrSign
	}

	if ok := service.CheckNode(in.Sign.Address, in.Net); !ok {
		log.Warning("address isn't a node's address")
		return &pb.CurrentCertResponse{
			Logid: in.Logid,
		}, ErrAuth
	}

	cert, nodeHdPriKey, err := service.GetNode(in.Net, in.Address)
	if err != nil {
		log.Warning("GetCurrentCert failed, err:", err)
		return nil, err
	}

	// 传输使用公钥进行加密 @todo
	return &pb.CurrentCertResponse{
		Logid:        in.Logid,
		CaCert:       cert.CaCert,
		Cert:         cert.Cert,
		PrivateKey:   cert.PrivateKey,
		NodeHdPriKey: nodeHdPriKey,
	}, nil
}

// 获取一个网络的增量撤销列表, 撤销时间为SerialNum证书撤销之后的
// 签名字段为serialNum+net
// 所有节点都可以访问, 暂不限制为网络内的节点
func (ca *caServer) GetRevokeList(ctx context.Context, in *pb.RevokeListRequest) (*pb.RevokeListResponse, error) {
	if in.GetLogid() == "" {
		in.Logid = util.GetLogId()
	}
	log.AddHook(util.NewlogIdHook(in.GetLogid()))

	if ok := verifyRequest(in.Sign, []byte(in.SerialNum+in.Net)); !ok {
		log.Error("GetRevokeList sign is not right")
		return &pb.RevokeListResponse{
			Logid: in.Logid,
		}, ErrSign
	}

	ret, err := service.GetRevokeList(in.Net, in.SerialNum)
	if err != nil {
		return nil, err
	}

	resp := &pb.RevokeListResponse{}
	for _, revokeNode := range *ret {
		resp.List = append(resp.List, &pb.RevokeNode{
			Id:         int64(revokeNode.Id),
			SerialNum:  revokeNode.SerialNum,
			CreateTime: int64(revokeNode.CreateTime),
		})
	}
	return resp, nil
}

// 撤销一个节点
// 签名字段为address+net
// conf中配置的ca根地址 及 网络的管理员地址具备该接口访问权限
func (ca *caServer) RevokeCert(ctx context.Context, in *pb.RevokeNodeRequest) (*pb.RevokeNodeResponse, error) {
	if in.GetLogid() == "" {
		in.Logid = util.GetLogId()
	}
	log.AddHook(util.NewlogIdHook(in.GetLogid()))

	if ok := verifyRequest(in.Sign, []byte(in.Address+in.Net)); !ok {
		log.Error("GetCurrentCert sign is not right")
		return &pb.RevokeNodeResponse{
			Logid: in.Logid,
		}, ErrSign
	}

	if ok := service.CheckNetAdmin(in.Sign.Address, in.Net); !ok {
		log.Warning("address isn't net admin")
		return &pb.RevokeNodeResponse{
			Logid: in.Logid,
		}, ErrAuth
	}

	ret, err := service.RevokeNode(in.Net, in.Address)
	if ret == true {
		return &pb.RevokeNodeResponse{
			Logid: in.Logid,
		}, nil
	}
	return nil, err
}

// 根据根私钥解密网络交易
// 签名字段为address+net+childHdpubKey+cypherText
// conf中配置的ca根地址、网络管理员地址可以操作该接口
func (ca *caServer) DecryptByHdKey(ctx context.Context, in *pb.DecryptByHdKeyRequest) (*pb.DecryptByHdKeyResponse, error) {
	if in.GetLogid() == "" {
		in.Logid = util.GetLogId()
	}
	log.AddHook(util.NewlogIdHook(in.GetLogid()))

	if ok := verifyRequest(in.Sign, []byte(in.Address+in.Net+in.ChildHdpubKey+in.CypherText)); !ok {
		log.Error("DecryptByHdKey sign is not right")
		return &pb.DecryptByHdKeyResponse{
			Logid: in.Logid,
		}, ErrSign
	}

	if ok := service.CheckNetAdmin(in.Sign.Address, in.Net); !ok {
		log.Warning("address isn't net admin")
		return &pb.DecryptByHdKeyResponse{
			Logid: in.Logid,
		}, ErrAuth
	}

	realMsg, err := service.DecryptByHdKey(in.Net, in.Address, in.ChildHdpubKey, in.CypherText)
	if err != nil {
		log.Warning("DecryptByHdKey failed, err:", err)
		return nil, err
	}

	// 传输使用公钥进行加密 @todo
	return &pb.DecryptByHdKeyResponse{
		Logid:   in.Logid,
		RealMsg: realMsg,
	}, nil
}

// 启动服务, 然后阻塞等待系统的关停信号
func Start(quit chan int) {
	/*
		// 注册优雅关停信号, 包括ctrl + C 和 kill 信号
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(sigc)
		// 失败退出
		quit := make(chan int)
	*/
	go startGrpcServer(quit)
	go startHttpServer(quit)
	/*
		for {
			select {
			case <-sigc:
				pprof.StopCPUProfile()
				return
			case <-quit:
				pprof.StopCPUProfile()
				return
			}
		}
	*/
}

func startGrpcServer(quit chan int) {
	// start server

	log.Info("caserver port is ", config.GetServerPort(), "start...")
	lis, err := net.Listen("tcp", config.GetServerPort())
	if err != nil {
		log.Errorf("failed to listen: %v \n", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(LogInterceptor))
	pb.RegisterCaserverServer(s, &caServer{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v  \n", err)
		quit <- 1
	}
}

func startHttpServer(quit chan int) {
	// start a http server at the same time
	log.Println("caserver http port is ", config.GetHttpPort(), "start...")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterCaserverHandlerFromEndpoint(ctx, mux, config.GetServerPort(), opts)
	if err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
	err = http.ListenAndServe(config.GetHttpPort(), mux)
	if err != nil {
		log.Printf("failed to serve: %v\n", err)
		quit <- 1
	}
}

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Debugf("before handling. Method: %+v, Req: %+v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Debugf("after handling. Method: %+v, Resp: %+v ,Err: %+v\n", info.FullMethod, resp, err)
	return resp, err
}
