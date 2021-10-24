package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello i'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect %v\n", err)
	}

	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)

	//doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)

	doErrorUnary(c)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot RPC...")
	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	fmt.Println("Starting to do a SquareRoot RPC...")
	// correct call
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
		}
	}

	log.Println("Result of square root of %v: %v", number, res.GetNumberRoot())

}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting do a client streaming gRPC")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 5,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 7,
		},
	}
	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}
	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receive response from LongGreet: %v", err)
	}

	fmt.Printf("ComputeAverage response: %v\n", res)

}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting do a server streaming gRPC")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecompositionRequest: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happend: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting do a Unary gRPC")

	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 7,
	}
	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Greet gRPC: %v", err)
	}

	log.Printf("Response from Calculator: %v", res.SumResult)
}
