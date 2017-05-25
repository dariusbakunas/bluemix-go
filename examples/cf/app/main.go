package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM-Bluemix/bluemix-go/api/cf/cfv2"
	"github.com/IBM-Bluemix/bluemix-go/session"
	"github.com/IBM-Bluemix/bluemix-go/trace"
)

func main() {

	var path string
	flag.StringVar(&path, "path", "", "Bluemix path for application")

	var org string
	flag.StringVar(&org, "org", "", "Bluemix Organization")

	var name string
	flag.StringVar(&name, "name", "", "Bluemix app name")

	var space string
	flag.StringVar(&space, "space", "", "Bluemix Space")

	var diego bool
	flag.BoolVar(&diego, "diego", false, "Bluemix Diego")

	var dockerImage string
	flag.StringVar(&dockerImage, "docker", "", "Docker image")

	var sharedDomain string
	flag.StringVar(&sharedDomain, "shared_domain", "mybluemix.net", "Bluemix shared domain")

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 120*time.Second, "Maximum time to wait for application to start")

	var routeName string
	flag.StringVar(&routeName, "route", "", "Bluemix app route")

	var buildpack string
	flag.StringVar(&buildpack, "buildpack", "", "Bluemix buildpack")

	var serviceOffering string
	flag.StringVar(&serviceOffering, "so", "cleardb", "Bluemix Service Offering")

	var servicePlan string
	flag.StringVar(&servicePlan, "plan", "spark", "Bluemix Service Plan")

	var instance int
	flag.IntVar(&instance, "instance", 2, "Bluemix App Instance")

	var serviceInstanceName string
	flag.StringVar(&serviceInstanceName, "svcname", "myservice", "Bluemix service instance name for the cloudantnosqldb offering")

	var memory int
	flag.IntVar(&memory, "memory", 128, "Bluemix app memory")

	var diskQuota int
	flag.IntVar(&diskQuota, "diskQuota", 512, "Bluemix app diskquota")

	var clean bool
	flag.BoolVar(&clean, "clean", false, "If set to true it will delete the application")

	flag.Parse()

	if name == "" || space == "" || org == "" || path == "" || routeName == "" {
		flag.Usage()
		os.Exit(1)
	}

	trace.Logger = trace.NewLogger("true")

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	client, err := cfv2.New(sess)

	if err != nil {
		log.Fatal(err)
	}
	orgAPI := client.Organizations()
	myorg, err := orgAPI.FindByName(org)

	if err != nil {
		log.Fatal(err)
	}

	spaceAPI := client.Spaces()
	myspace, err := spaceAPI.FindByNameInOrg(myorg.GUID, space)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(myorg.GUID, myspace.GUID)

	serviceOfferingAPI := client.ServiceOfferings()
	myserviceOff, err := serviceOfferingAPI.FindByLabel(serviceOffering)
	if err != nil {
		log.Fatal(err)
	}
	servicePlanAPI := client.ServicePlans()
	plan, err := servicePlanAPI.GetServicePlan(myserviceOff.GUID, servicePlan)
	if err != nil {
		log.Fatal(err)
	}

	serviceInstanceAPI := client.ServiceInstances()
	myService, err := serviceInstanceAPI.Create(serviceInstanceName, plan.GUID, myspace.GUID, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	appAPI := client.Apps()
	_, err = appAPI.FindByName(myspace.GUID, name)

	if err == nil {
		log.Fatal(err)
	}

	var appPayload = &cfv2.AppCreateRequest{
		Name:        name,
		SpaceGUID:   myspace.GUID,
		BuildPack:   buildpack,
		Instances:   instance,
		Memory:      memory,
		DiskQuota:   diskQuota,
		Diego:       diego,
		DockerImage: dockerImage,
	}

	newapp, err := appAPI.Create(appPayload)
	if err != nil {
		log.Fatal(err)
	}

	sharedDomainAPI := client.SharedDomains()
	domain, err := sharedDomainAPI.FindByName(sharedDomain)
	fmt.Println(domain)
	if err != nil {
		log.Fatal(err)
	}

	routeAPI := client.Routes()
	route, err := routeAPI.Find(routeName, domain.GUID)
	fmt.Println(route)
	if err != nil {
		log.Fatal(err)
	}

	if len(route) == 0 {
		req := cfv2.RouteRequest{
			Host:       routeName,
			DomainGUID: domain.GUID,
			SpaceGUID:  myspace.GUID,
		}
		newroute, err := routeAPI.Create(req)
		fmt.Println(newroute)
		if err != nil {
			log.Fatal(err)
		}
		bindRoute, err := appAPI.BindRoute(newapp.Metadata.GUID, newroute.Metadata.GUID)
		fmt.Println(bindRoute)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		bindRoute, err := appAPI.BindRoute(newapp.Metadata.GUID, route[0].GUID)
		fmt.Println(bindRoute)
		if err != nil {
			log.Fatal(err)
		}

	}
	if dockerImage == "" {
		uploadResponse, err := appAPI.Upload(newapp.Metadata.GUID, path)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(uploadResponse)
	}

	appState, err := appAPI.Start(newapp.Metadata.GUID, timeout)
	if err != nil {
		log.Fatal(err)
	}
	if appState.PackageState != cfv2.AppStagedState {
		log.Fatalf("Application couldn't be staged, current status is  %s", appState.PackageState)
	}
	if appState.InstanceState != cfv2.AppRunningState {
		log.Fatalf("Application is not yet running, current status is  %s", appState.InstanceState)
	}

	sbAPI := client.ServiceBindings()

	sb, err := sbAPI.Create(cfv2.ServiceBindingRequest{
		ServiceInstanceGUID: myService.Metadata.GUID,
		AppGUID:             newapp.Metadata.GUID,
	})

	if err != nil {
		log.Fatal(err)
	}
	sbFields, err := sbAPI.Get(sb.Metadata.GUID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*sbFields)

	apps, err := appAPI.Get(newapp.Metadata.GUID)
	fmt.Println(apps)
	if err != nil {
		log.Fatal(err)
	}
	var appUpdatePayload = &cfv2.AppUpdateRequest{
		Name:      "testappupdate",
		SpaceGUID: myspace.GUID,
	}

	updateapp, err := appAPI.Update(newapp.Metadata.GUID, appUpdatePayload)
	fmt.Println(updateapp)
	if err != nil {
		log.Fatal(err)
	}

	appInstances, err := appAPI.Instances(updateapp.Metadata.GUID)
	fmt.Println(appInstances)
	if err != nil {
		log.Fatal(err)
	}

	if clean {
		err = appAPI.Delete(updateapp.Metadata.GUID)
		if err != nil {
			log.Fatal(err)
		}
	}

}