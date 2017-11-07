package main

import (
  "fmt"
  "os"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/0xAX/notificator"
)

var notify *notificator.Notificator

func main() {

  validCmd := []string{"start", "stop"}

  notify = notificator.New(notificator.Options{
    AppName: "ec2_control",
  })

  if len(os.Args) < 2 {
    fmt.Println("Error!")
    os.Exit(1)
  }

  ec2Cmd := os.Args[1]
  ec2Id := os.Args[2]

  if stringInSlice( ec2Cmd, validCmd ) {
    // do nuffing
  } else {
    fmt.Printf( "%s is not a valid command!\n", ec2Cmd)
    os.Exit(1)
  }

  sess := session.Must(session.NewSessionWithOptions( session.Options{
    SharedConfigState: session.SharedConfigEnable,
  }))

  // create an EC2 Client
  ec2Svc := ec2.New( sess )

  if ec2Cmd == "start" {
    input := &ec2.StartInstancesInput{
      InstanceIds: []*string{
        aws.String(ec2Id),
      },
      DryRun: aws.Bool(false),
    }

    result, err := ec2Svc.StartInstances(input)
    if err != nil {
      fmt.Println(err)
    } else {
      notify.Push(
        "Starting AWS Instance",
        aws.StringValue(result.StartingInstances[0].CurrentState.Name),
        "",
        notificator.UR_NORMAL,
      )
    }
  }

  if ec2Cmd == "stop" {
    input := &ec2.StopInstancesInput{
      InstanceIds: []*string{
        aws.String(ec2Id),
      },
      DryRun: aws.Bool(false),
    }

    result, err := ec2Svc.StopInstances(input)
    if err != nil {
      fmt.Println(err)
    } else {
      notify.Push(
        "Stopping AWS Instance",
        aws.StringValue(result.StoppingInstances[0].CurrentState.Name),
        "",
        notificator.UR_NORMAL,
      )
    }
  }
}

func stringInSlice(a string, list []string ) bool {
  for _, b := range list {
    if b == a {
      return true
    }
  }
  return false
}
