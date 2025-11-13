package web

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"log"

	"github.com/glekoz/online-shop_proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// временно здесь
// type UserService struct {
// 	Client    user.UserClient
// 	publicKey *rsa.PublicKey
// }

type Weber struct {
	// временно здесь
	UserClient    user.UserClient
	PublicKeyUser *rsa.PublicKey
}

func New() *Weber {
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("gRPC connection can not be established")
	}
	client := user.NewUserClient(conn)
	pub, err := client.GetRSAPublicKey(context.Background(), &user.Empty{})
	if err != nil {
		log.Fatal("gRPC call failed")
	}
	der := make([]byte, base64.StdEncoding.DecodedLen(len(pub.Key)))
	base64.StdEncoding.Decode(der, pub.Key)
	pubKey, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		log.Fatal("key parsing failed")
	}
	publicKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		log.Fatal("key is not RSA")
	}
	return &Weber{
		UserClient:    client,
		PublicKeyUser: publicKey,
	}
}
