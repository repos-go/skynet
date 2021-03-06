package main

import (
  "fmt"
  "flag"
  "github.com/bketelsen/skynet/skylib"
  "github.com/4ad/doozer"
  "os"
  "strconv"
)

var VersionFlag *string = flag.String("version", "", "service version")
var ServiceNameFlag *string = flag.String("service", "", "service name")
var HostFlag *string = flag.String("host", "", "host")
var RegionFlag *string = flag.String("region", "", "region")

var DC *doozer.Conn

func main() {
  flag.Parse()
  Connect()

  query := &skylib.Query{
    DoozerConn: DC,
    Service: *ServiceNameFlag,
    Version: *VersionFlag,
    Host: *HostFlag,
    Region: *RegionFlag,
  }

  switch flag.Arg(0) {
    case "help", "h":
      Help()
    case "services":
      ListServices(query)
    case "hosts":
      ListHosts(query)
    case "regions":
      ListRegions(query)
    case "instances":
      ListInstances(query)
    case "versions":
      ListServiceVersions(query)
    case "topology":
      PrintTopology(query)

    default:
      Help()
  }
}

func Connect(){
  defer func() {
      if r := recover(); r != nil {
          fmt.Println("Failed to connect to Doozer")
          os.Exit(1)
      }
  }()

  DC = skylib.DoozerConnect();
}

func ListInstances(q *skylib.Query){
  results := q.FindInstances()

  for _, instance := range *results {
    fmt.Println(instance.IPAddress + ":" + strconv.Itoa(instance.Port) + " - " + instance.Name + " (" + instance.Version + ")")
  }
}

func ListHosts(q *skylib.Query){
  results := q.FindHosts()

  for _, host := range *results {
    fmt.Println(host)
  }
}

func ListRegions(q *skylib.Query){
  results := q.FindRegions()

  for _, region := range *results {
    fmt.Println(region)
  }
}

func ListServices(q *skylib.Query){
  results := q.FindServices()

  for _, service := range *results {
    fmt.Println(service)
  }
}

func ListServiceVersions(q *skylib.Query){
  if *ServiceNameFlag == "" { 
    fmt.Println("Service name is required")
    os.Exit(1)
  }

  results := q.FindServiceVersions()

  for _, version := range *results {
    fmt.Println(version)
  }
}

func PrintTopology(q *skylib.Query){
  topology := make(map[string]map[string]map[string]map[string][]*skylib.Service)

  results := q.FindInstances()

  // Build topology hash first
  for _, instance := range *results {
    if topology[instance.Region] == nil {
      topology[instance.Region] = make(map[string]map[string]map[string][]*skylib.Service)
    }

    if topology[instance.Region][instance.IPAddress] == nil {
      topology[instance.Region][instance.IPAddress] = make(map[string]map[string][]*skylib.Service)
    }

    if topology[instance.Region][instance.IPAddress][instance.Name] == nil {
      topology[instance.Region][instance.IPAddress][instance.Name] = make(map[string][]*skylib.Service)
    }

    if topology[instance.Region][instance.IPAddress][instance.Name][instance.Version] == nil {
      topology[instance.Region][instance.IPAddress][instance.Name][instance.Version] = make([]*skylib.Service, 0)
    }

    topology[instance.Region][instance.IPAddress][instance.Name][instance.Version] = append(topology[instance.Region][instance.IPAddress][instance.Name][instance.Version], instance)
  }


  // Now we can print the correct heirarchy
  for regionName, region := range topology {
    fmt.Println("Region: " + regionName)

    for hostName, host := range region {
      fmt.Println("\tHost: " + hostName)

      for serviceName, service := range host {
        fmt.Println("\t\tService: " + serviceName)

        for versionName, version := range service {
          fmt.Println("\t\t\tVersion: " + versionName)

          for _, instance := range version {
            fmt.Println("\t\t\t\t" + instance.IPAddress + ":" + strconv.Itoa(instance.Port) + " - " + instance.Name + " (" + instance.Version + ")")
          }
        }
      }
    }
  }
}

func Help(){
  // TODO: check to see if a specific command is in the arguments List
  // if so just desplay info for that

  fmt.Println("Usage:\n\t sky -option1=value -option2=value command <arguments>")

  fmt.Print(
    "\nCommands:\n" +
    "\n\thosts: List all hosts available that meet the specified criteria" +
    "\n\t\t-service - limit results to hosts running the specified service" +
    "\n\t\t-version - limit results to hosts running the specified version of the service (-service required)" +
    "\n\t\t-region - limit results to hosts in the specified region" +

    "\n\tinstances: List all instances available that meet the specified criteria" +
    "\n\t\t-service - limit results to instances of the specified service" +
    "\n\t\t-version - limit results to instances of the specified version of service" +
    "\n\t\t-region - limit results to instances in the specified region" +
    "\n\t\t-host - limit results to instances on the specified host" +


    "\n\tregions: List all regions available that meet the specified criteria" +

    "\n\tservices: List all services available that meet the specified criteria" +
    "\n\t\t-host - limit results to the specified host" +
    "\n\t\t-region - limit results to hosts in the specified region" +

    "\n\n\tservice-versions: List all services available that meet the specified criteria" +
    "\n\t\t-service - service name (required)" +
    "\n\t\t-host - limit results to the specified host" +
    "\n\t\t-region - limit results to hosts in the specified region" +

    "\n\n\ttopology: Print detailed heirarchy of regions/hosts/services/versions/instances" +
    "\n\t\t-service - limit results to instances of the specified service" +
    "\n\t\t-version - limit results to instances of the specified version of service" +
    "\n\t\t-region - limit results to instances in the specified region" +
    "\n\t\t-host - limit results to instances on the specified host" +

  "\n\n\n")
}
